/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/omnitrix-sh/cli/internals/app"
	"github.com/omnitrix-sh/cli/internals/database"
	"github.com/omnitrix-sh/cli/internals/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "omnitrix.sh",
	Short: "cursor for terminal",
	Long:  `a claude cli alternative with no vendor lock-in`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flag("help").Changed {
			cmd.Help()
			return nil
		}
		debug, _ := cmd.Flags().GetBool("debug")
		viper.Set("debug", debug)
		if debug {
			viper.Set("log.level", "debug")
		}

		conn, err := database.Connect()
		if err != nil {
			return err
		}
		ctx := context.Background()

		app := app.New(ctx, conn)
		app.Logger.Info("Starting termai...")
		tui := tea.NewProgram(
			tui.New(app),
			tea.WithAltScreen(),
		)
		app.Logger.Info("Setting up subscriptions...")
		ch, unsub := setupSubscriptions(app)
		defer unsub()

		go func() {
			for msg := range ch {
				tui.Send(msg)
			}
		}()
		if _, err := tui.Run(); err != nil {
			return err
		}
		return nil
	},
}

func setupSubscriptions(app *app.App) (chan tea.Msg, func()) {
	ch := make(chan tea.Msg)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(app.Context)

	if viper.GetBool("debug") {
		sub := app.Logger.Subscribe(ctx)
		wg.Add(1)
		go func() {
			for ev := range sub {
				ch <- ev
			}
			wg.Done()
		}()
	}
	{
		sub := app.Sessions.Subscribe(ctx)
		wg.Add(1)
		go func() {
			for ev := range sub {
				ch <- ev
			}
			wg.Done()
		}()
	}
	return ch, func() {
		cancel()
		wg.Wait()
		close(ch)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func loadConfig() {
	viper.SetConfigName(".omnitrix")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$XDG_CONFIG_HOME/omnitrix")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("OMNITRIX")
	// SET DEFAULTS
	viper.SetDefault("log.level", "info")
	viper.SetDefault("data.dir", ".omnitrix")

	//
	viper.ReadInConfig()
}

func init() {
	loadConfig()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.Flags().BoolP("help", "h", false, "Help")
	rootCmd.Flags().BoolP("debug", "d", false, "Help")
}
