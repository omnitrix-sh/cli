package logging

import (
	"context"

	"github.com/omnitrix-sh/cli/internals/pubsub"
)

type Interface interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Subscribe(ctx context.Context) <-chan pubsub.Event[LogMessage]

	List() []LogMessage
}
