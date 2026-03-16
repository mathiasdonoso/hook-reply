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

// TODO: define flag descriptions in a single place to avoid duplication
var CommandsRegistry = map[string]*CommandNode{
	"serve": {
		Name: "serve",
		Description: `Creates a proxy server that intercepts and forwards webhook requests.
		Flags:

		--port:
		Port where the webhook server listens for incoming requests. Default: 3000.

		--forward:
		URL where captured webhook requests will be forwarded.
		`,
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
		Description: `Get the last 20 requests received.`,
		Setup: func(fs *flag.FlagSet) func() error {
			return func() error {
				return handlers.LogHandler()
			}
		},
	},
	"replay": {
		Name: "replay",
		Description: `Replay the selected request to the configured forward target.
		Flags:

		--last:
		Replay the most recent event. Default: false.

		--times:
		Number of times to replay the event (used with --delay). Default: 1.

		--delay:
		Delay between each replay (in milliseconds). Default: 0.

		--target:
		Override the forward target for the replayed request.
		`,
		Setup: func(fs *flag.FlagSet) func() error {
			last := fs.Bool("last", false, "Replay the most recent event.")
			times := fs.Uint("times", 1, "Number of times to replay the event (used with --delay).")
			delay := fs.Uint("delay", 0, "Delay between each replay (in milliseconds).")
			target := fs.String("target", "", "Override the forward target for the replayed request.")

			return func() error {
				return handlers.ReplayHandler(fs.Arg(0), *last, *times, *delay, *target)
			}
		},
	},
}
