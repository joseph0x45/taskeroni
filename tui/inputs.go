package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joseph0x45/taskeroni/shared"
	"github.com/joseph0x45/taskeroni/utils"
)

func (app App) switchFocus() (tea.Model, tea.Cmd) {
	switch app.currentScreen {
	case "auth":
		for _, inputName := range []string{"usernameInput", "passwordInput"} {
			input := app.inputs[inputName]
			if input.Focused() {
				input.Blur()
			} else {
				input.Focus()
			}
			app.inputs[inputName] = input
		}
	}
	return app, nil
}

func (app App) handleAuth() (tea.Model, tea.Cmd) {
	usernameInput := app.inputs["usernameInput"]
	if usernameInput.Value() == "" {
		app.authScreenErrorText = "Username can not be empty"
		return app, nil
	}
	passwordInput := app.inputs["passwordInput"]
	if passwordInput.Value() == "" {
		app.authScreenErrorText = "Password can not be empty"
		return app, nil
	}
	user, err := app.conn.GetUserByUsername(usernameInput.Value())
	if err != nil {
		app.authScreenErrorText = "Something went wrong! Check logs for more information"
		return app, nil
	}
	if user == nil {
		app.authScreenErrorText = "Invalid credentials"
		return app, nil
	}
	if !utils.HashMatchesPassword(user.Password, passwordInput.Value()) {
		app.authScreenErrorText = "Invalid credentials"
		return app, nil
	}
	app.currentUserID = user.ID
	tasks, err := app.conn.GetUserTasks(app.currentUserID)
	if err != nil {
		app.authScreenErrorText = "Failed to load user tasks. Check logs for more information"
		return app, nil
	}
	app.tasks = tasks
	app.currentScreen = "tasks"
	//init tables
	columns := []table.Column{
		{Title: "Label", Width: 10},
		{Title: "Status", Width: 10},
		{Title: "DueDate", Width: 10},
	}
	rows := []table.Row{}
	for _, task := range app.tasks {
		row := []string{
			task.Label,
			task.Status,
			task.DueDate,
		}
		rows = append(rows, row)
	}
	tasksTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(50),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	tasksTable.SetStyles(s)
	app.tasksTable = tasksTable
	app.displayedTasks = shared.DisplayedTasksBacklog
	app.state = shared.StateBrowsing
	return app, nil
}

// TODO: rename later
func (app App) SwitchDisplayedTasks() (tea.Model, tea.Cmd) {
	return app, nil
}
