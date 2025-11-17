package main

import (
	"github.com/omnitrix-sh/cli/cmd"
	"github.com/omnitrix-sh/cli/internal/logging"
)

func main() {
	defer logging.RecoverPanic("main", func() {
		logging.ErrorPersist("Application terminated due to unhandled panic")
	})

	cmd.Execute()
}
