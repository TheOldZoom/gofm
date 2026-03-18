package setup

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/theOldZoom/gofm/internal/config"
	"github.com/theOldZoom/gofm/internal/verbose"
)

func Run() (*config.Config, error) {
	m := NewModel()
	verbose.Printf("starting interactive setup")

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		verbose.Printf("interactive setup failed: %v", err)
		return nil, err
	}

	result := finalModel.(Model)
	cfg, err := result.Result()
	if err != nil {
		verbose.Printf("interactive setup result error: %v", err)
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("setup cancelled")
	}

	verbose.Printf("interactive setup finished for username=%s", cfg.Username)
	return cfg, nil
}
