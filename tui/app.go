package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joseph0x45/taskeroni/shared"
)

type KeyHandler func(App) (App, tea.Cmd)

func (app App) handleEscape() (tea.Model, tea.Cmd) {
	if app.state == shared.StateBrowsing {
		return app, tea.Quit
	}
	app.state = shared.StateBrowsing
	return app, nil
}
