package main

import (
	"fmt"
	"os"

	"github.com/vibolsovichea/etch/internal/config"
	"github.com/vibolsovichea/etch/internal/ui"
	"github.com/vibolsovichea/etch/internal/version"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("etch %s (commit: %s, built: %s)\n", version.Version, version.Commit, version.Date)
		os.Exit(0)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if cfg == nil {
		p := tea.NewProgram(ui.NewSetupModel())
		m, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		setup := m.(ui.SetupModel)
		if setup.Cancelled {
			fmt.Println("Setup cancelled.")
			os.Exit(0)
		}

		cfg, err = config.Init(setup.VaultPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing vault: %v\n", err)
			os.Exit(1)
		}
	}

	p := tea.NewProgram(ui.NewAppModel(cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
