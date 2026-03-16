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
		Name: "serve",
		Description: "Creates a proxy server that intercepts and forwards webhook requests.",
		Setup: func(fs *flag.FlagSet) func() error {
			port := fs.Uint("port", DefaultPort, "Port where the webhook server listens for incoming requests.")
			forward := fs.String("forward", "", "URL to forward captured webhook requests to.")

			return func() error {
				return handlers.ServeHandler(*port, *forward)
			}
		},
	},
	"log": {
		Name: "log",
		Description: "Get the last 20 requests received.",
		Setup: func(fs *flag.FlagSet) func() error {
			return func() error {
				return handlers.LogHandler()
			}
		},
	},
	"replay": {
		Name: "replay",
		Description: "Replay the selected request to the configured forward target.",
		Setup: func(fs *flag.FlagSet) func() error {
			last := fs.Bool("last", false, "Replay the most recent event.")
			times := fs.Uint("times", 1, "Number of times to replay the event (used with --delay).")
			delay := fs.Uint("delay", 0, "Delay between each replay")
			target := fs.String("target", "", "Override the forward target for the replayed request.")

			return func() error {
				return handlers.ReplayHandler(*last, *times, *delay, *target)
			}
		},
	},
}
