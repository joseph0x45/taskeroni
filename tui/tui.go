package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joseph0x45/taskeroni/db"
	"github.com/joseph0x45/taskeroni/internal/models"
)

type App struct {
	currentUserID       string
	currentScreen       string
	tasks               []models.Task
	tasksTable          table.Model
	renderer            *lipgloss.Renderer
	conn                *db.Conn
	inputs              map[string]textinput.Model
	authScreenErrorText string
	displayedTasks      string //TODO: find better name
	state               string
}

func (app *App) renderAuthScreen() string {
	ui := "Log into Taskeroni('tab' to switch focus. Enter to submit form)\n"
	ui += app.inputs["usernameInput"].View()
	ui += "\n"
	ui += app.inputs["passwordInput"].View()
	ui += "\n"
	ui += redTextStyle.Render(app.authScreenErrorText)
	return ui
}

func (app App) View() string {
	if app.currentUserID == "" {
		return app.renderAuthScreen()
	}
	return app.renderTasksScreen()
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return app, tea.Quit
		case "enter":
			return app.handleAuth()
		case "tab":
			return app.switchFocus()
		case "s":
			// case "r":
			// 	//TODO: refresh tasks
			// 	if app.currentScreen == "tasks" {
			// 		return app, tea.Quit
			// 	}
			// 	return app, nil
		}
	}
	udpatedTasksTable, tableUpdateCmd := app.tasksTable.Update(msg)
	app.tasksTable = udpatedTasksTable
	return app, tea.Batch(app.updateInputs(msg), tableUpdateCmd)
}

func (app *App) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(app.inputs))
	for inputName := range app.inputs {
		teaModel, teaCmd := app.inputs[inputName].Update(msg)
		app.inputs[inputName] = teaModel
		cmds = append(cmds, teaCmd)
	}
	return tea.Batch(cmds...)
}

func (app App) Init() tea.Cmd {
	return nil
}

func initAppInputs(inputsMap map[string]textinput.Model) {
	usernameInput := textinput.New()
	usernameInput.Focus()
	usernameInput.Placeholder = "Username"
	usernameInput.Width = 32
	inputsMap["usernameInput"] = usernameInput
	passwordInput := textinput.New()
	passwordInput.Placeholder = "Password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '*'
	passwordInput.Width = 32
	inputsMap["passwordInput"] = passwordInput
}

func InitApp(renderer *lipgloss.Renderer, conn *db.Conn) *App {
	appInputsMap := make(map[string]textinput.Model)
	initAppInputs(appInputsMap)
	return &App{
		currentUserID:       "",
		currentScreen:       "auth",
		tasks:               nil,
		renderer:            renderer,
		conn:                conn,
		inputs:              appInputsMap,
		authScreenErrorText: "",
	}
}
