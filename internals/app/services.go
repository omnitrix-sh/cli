package app

import (
	"context"
	"database/sql"

	"github.com/omnitrix-sh/cli/internals/config"
	"github.com/omnitrix-sh/cli/internals/database"
	"github.com/omnitrix-sh/cli/internals/logging"
	"github.com/omnitrix-sh/cli/internals/message"
	permission "github.com/omnitrix-sh/cli/internals/permissions"
	session "github.com/omnitrix-sh/cli/internals/sessions"
)

type App struct {
	Context context.Context

	Sessions    session.Service
	Messages    message.Service
	Permissions permission.Service

	Logger logging.Interface
}

func New(ctx context.Context, conn *sql.DB) *App {
	q := database.New(conn)
	log := logging.NewLogger(logging.Options{
		Level: config.Get().Log.Level,
	})
	sessions := session.NewService(ctx, q)
	messages := message.NewService(ctx, q)

	return &App{
		Context:     ctx,
		Sessions:    sessions,
		Messages:    messages,
		Permissions: permission.Default,
		Logger:      log,
	}
}
