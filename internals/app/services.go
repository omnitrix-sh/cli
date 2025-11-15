package app

import (
	"context"
	"database/sql"

	"github.com/omnitrix-sh/cli/internals/config"
	"github.com/omnitrix-sh/cli/internals/database"
	"github.com/omnitrix-sh/cli/internals/logging"
	"github.com/omnitrix-sh/cli/internals/lsp"
	"github.com/omnitrix-sh/cli/internals/lsp/watcher"
	"github.com/omnitrix-sh/cli/internals/message"
	permission "github.com/omnitrix-sh/cli/internals/permissions"
	"github.com/omnitrix-sh/cli/internals/pubsub"
	session "github.com/omnitrix-sh/cli/internals/sessions"
	util "github.com/omnitrix-sh/cli/internals/utils"
)

type App struct {
	Context context.Context

	Sessions    session.Service
	Messages    message.Service
	Permissions permission.Service

	LSPClients map[string]*lsp.Client

	Logger logging.Interface

	Status  *pubsub.Broker[util.InfoMsg]
	ceanups []func()
}

func New(ctx context.Context, conn *sql.DB) *App {
	cfg := config.Get()
	q := database.New(conn)
	log := logging.NewLogger(logging.Options{
		Level: cfg.Log.Level,
	})
	sessions := session.NewService(ctx, q)
	messages := message.NewService(ctx, q)

	app := &App{
		Context:     ctx,
		Sessions:    sessions,
		Messages:    messages,
		Permissions: permission.Default,
		Logger:      log,
		Status:      pubsub.NewBroker[util.InfoMsg](),
		LSPClients:  make(map[string]*lsp.Client),
	}

	for name, client := range cfg.LSP {
		lspClient, err := lsp.NewClient(client.Command, client.Args...)
		app.ceanups = append(app.ceanups, func() {
			lspClient.Close()
		})
		workspaceWatcher := watcher.NewWorkspaceWatcher(lspClient)
		if err != nil {
			log.Error("Failed to create LSP client for", name, err)
			continue
		}

		_, err = lspClient.InitializeLSPClient(ctx, config.WorkingDirectory())
		if err != nil {
			log.Error("Initialize failed", "error", err)
			continue
		}
		go workspaceWatcher.WatchWorkspace(ctx, config.WorkingDirectory())
		app.LSPClients[name] = lspClient
	}
	return app
}

func (a *App) Close() {
	for _, cleanup := range a.ceanups {
		cleanup()
	}
	for _, client := range a.LSPClients {
		client.Close()
	}
	a.Logger.Info("App closed")
}
