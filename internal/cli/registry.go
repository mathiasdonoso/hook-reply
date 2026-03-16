package cli

import (
	"flag"

	"github.com/mathiasdonoso/hook-replay/internal/cli/handlers"
)

type CommandNode struct {
	Name        string
	Description string
	Setup       func(fs *flag.FlagSet) func() error
}

const DefaultPort = 3000

var CommandsRegistry = map[string]*CommandNode{
	"serve": {
		Name:        "serve",
		Description: "Creates a proxy server that intercepts and forwards webhook requests.",
		Setup: func(fs *flag.FlagSet) func() error {
			port := fs.Uint("port", DefaultPort, "Port where the webhook server listens for incoming requests.")
			forward := fs.String("forward", "", "URL to forward captured webhook requests to.")

			return func() error {
				return handlers.ServeHandler(*port, *forward)
			}
		},
	},
}
