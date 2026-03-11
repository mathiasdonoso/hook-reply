package cli

import (
	"flag"
	"fmt"
)

type HR struct{}

func NewHR(args []string) *HR {
	return &HR{}
}

func (h *HR) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	command := args[0]
	node, ok := CommandsRegistry[command]
	if !ok {
		return fmt.Errorf("unknown command: %s", command)
	}

	fs := flag.NewFlagSet(command, flag.ExitOnError)
	handler := node.Setup(fs)
	fs.Parse(args[1:])

	return handler()
}
