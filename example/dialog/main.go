package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xhd2015/less-gen/flags"
)

func main() {
	err := Handle(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func Handle(args []string) error {
	var show bool
	var showDialog bool
	args, err := flags.Bool("--show", &show).
		Bool("--show-dialog", &showDialog).
		Parse(args)
	if err != nil {
		return err
	}
	if len(args) > 0 {
		return fmt.Errorf("unrecognized extra args: %s", strings.Join(args, " "))
	}

	// Create and configure model
	model := NewModel()
	if showDialog {
		model.showDialog = true
	}

	if show {
		// Just output the rendered view with ANSI colors stripped
		output := model.ViewStripped()
		fmt.Println(output)
		return nil
	}

	// Run the interactive program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
