package app

import (
	"context"
	"database/sql"

	"github.com/omnitrix-sh/cli/internals/database"
	"github.com/omnitrix-sh/cli/internals/logging"
	session "github.com/omnitrix-sh/cli/internals/sessions"

	"github.com/spf13/viper"
)

type App struct {
	Context  context.Context
	Sessions session.Service
	Logger   logging.Interface
}

func New(ctx context.Context, conn *sql.DB) *App {
	q := database.New(conn)
	log := logging.NewLogger(logging.Options{
		Level: viper.GetString("log.level"),
	})
	return &App{
		Context:  ctx,
		Sessions: session.NewService(ctx, q),
		Logger:   log,
	}
}
