package cmd

import (
	"fmt"
	"os"

	"github.com/omnitrix-sh/cli/internal/config"
	"github.com/spf13/cobra"
)

var automodeCmd = &cobra.Command{
	Use:   "automode [on|off]",
	Short: "Enable or disable automatic permission approval",
	Long: `Enable or disable automatic permission approval mode.

When automode is enabled, all permission requests are automatically approved without showing dialogs.
This reduces interruptions but automatically grants all tool permissions.

Use with caution as this will automatically approve potentially destructive operations.`,
	Args: cobra.ExactArgs(1),
	ValidArgs: []string{"on", "off", "enable", "disable", "true", "false", "status"},
	Example: `
  # Enable automode
  omnitrix automode on

  # Disable automode  
  omnitrix automode off
  
  # Check current status
  omnitrix automode status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config first to ensure it's initialized
		cwd, _ := cmd.Flags().GetString("cwd")
		debug, _ := cmd.Flags().GetBool("debug")
		
		if cwd == "" {
			var err error
			cwd, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %v", err)
			}
		}
		
		_, err := config.Load(cwd, debug)
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		switch args[0] {
		case "on", "enable", "true":
			err := config.SetAutoMode(true)
			if err != nil {
				return fmt.Errorf("failed to enable automode: %v", err)
			}
			fmt.Println("Automode enabled - all permissions will be automatically approved")
		case "off", "disable", "false":
			err := config.SetAutoMode(false)
			if err != nil {
				return fmt.Errorf("failed to disable automode: %v", err)
			}
			fmt.Println("Automode disabled - permission dialogs will be shown")
		case "status":
			if config.AutoModeEnabled() {
				fmt.Println("Automode is currently enabled")
			} else {
				fmt.Println("Automode is currently disabled")
			}
		default:
			return fmt.Errorf("invalid argument: %s. Use 'on', 'off', or 'status'", args[0])
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(automodeCmd)
}