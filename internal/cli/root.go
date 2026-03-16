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
		showHelp()
		return nil
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

func showHelp() {
	commands := CommandsRegistry

	fmt.Printf("A CLI proxy server that intercepts, captures, and replays webhook requests.\n\n")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("\thr <command> [arguments]\n\n")
	fmt.Printf("The commands are:\n\n")

	for name, command := range commands {
		fmt.Printf("\t%s\t%s\n", name, command.Description)
	}

	fmt.Println()
}
