package app

import (
	"context"
	"database/sql"

	"github.com/omnitrix-sh/cli/internals/database"
	"github.com/omnitrix-sh/cli/internals/llm"
	"github.com/omnitrix-sh/cli/internals/logging"
	"github.com/omnitrix-sh/cli/internals/message"
	session "github.com/omnitrix-sh/cli/internals/sessions"

	"github.com/spf13/viper"
)

type App struct {
	Context context.Context

	Sessions session.Service
	Messages message.Service
	LLM      llm.Service

	Logger logging.Interface
}

func New(ctx context.Context, conn *sql.DB) *App {
	q := database.New(conn)
	log := logging.NewLogger(logging.Options{
		Level: viper.GetString("log.level"),
	})
	sessions := session.NewService(ctx, q)
	messages := message.NewService(ctx, q)
	llm := llm.NewService(ctx, log, sessions, messages)

	return &App{
		Context:  ctx,
		Sessions: sessions,
		Messages: messages,
		LLM:      llm,
		Logger:   log,
	}
}
