package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xhd2015/less-gen/flags"
)

const help = `
Interactive Todo Demo with DOM-based event handling

Usage: interactive_demo [OPTIONS]

Options:
  --debug-file <file>              enable debug logging to specified file
  -h,--help                        show help message

Examples:
  interactive_demo                  run the interactive demo
  interactive_demo --debug-file debug.log    run with debug logging enabled
`

func main() {
	err := handle(os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func handle(args []string) error {
	var debugFile string
	_, err := flags.String("--debug-file", &debugFile).
		Help("-h,--help", help).
		Parse(args)
	if err != nil {
		return err
	}

	// Create and run the Bubble Tea program
	model := NewModel(debugFile)
	p := tea.NewProgram(model, tea.WithAltScreen())
	model.program = p
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}
	return nil
}

type fileLogger struct {
	file *os.File
}

func (l *fileLogger) Log(format string, args ...interface{}) {
	newLineFmt := format
	if !strings.HasSuffix(format, "\n") {
		newLineFmt += "\n"
	}
	fmt.Fprintf(l.file, newLineFmt, args...)
}
