package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "github.com/charmbracelet/log"
	"github.com/joseph0x45/taskeroni/db"
	"github.com/joseph0x45/taskeroni/internal/models"
	"github.com/joseph0x45/taskeroni/shared"
)

type App struct {
	board      string
	tasks      []models.Task
	tasksTable table.Model
	renderer   *lipgloss.Renderer
	conn       *db.Conn
	inputs     map[string]textinput.Model
	state      string
}

func (app App) View() string {
	if app.state == shared.StateBrowsing {
		return app.renderTasks()
	}
	return app.renderTaskCreationForm()
}

func (app App) nextBoard() string {
	switch app.board {
	case shared.BacklogBoard:
		return shared.InProgressBoard
	case shared.InProgressBoard:
		return shared.DoneBoard
	case shared.DoneBoard:
		return shared.BacklogBoard
	}
	return app.board
}

func (app App) buildTasksTable(tasks []models.Task) table.Model {
	columns := []table.Column{
		{Title: "Label", Width: 40},
		{Title: "Due Date", Width: 40},
	}
	rows := []table.Row{}
	for _, task := range tasks {
		rows = append(rows, table.Row{task.Label, task.DueDate})
	}
	tasksTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
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
	return tasksTable
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return app, tea.Quit
		case "s":
			app.board = app.nextBoard()
			tasks, err := app.conn.GetTasks(app.board)
			if err != nil {
				//TODO: Improve this
				panic(err)
			}
			app.tasks = tasks
			app.tasksTable = app.buildTasksTable(tasks)
		case "r":
			tasks, err := app.conn.GetTasks(app.board)
			if err != nil {
				//TODO: Improve this
				panic(err)
			}
			app.tasks = tasks
			app.tasksTable = app.buildTasksTable(tasks)
		case "n":
			app.state = shared.StateCreating
		case "esc":
			return app.handleEscape()
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
}

func InitApp(renderer *lipgloss.Renderer, conn *db.Conn) *App {
	appInputsMap := make(map[string]textinput.Model)
	initAppInputs(appInputsMap)
	tasks, err := conn.GetTasks(shared.BacklogBoard)
	if err != nil {
		panic(err)
	}
	app := &App{
		board:    shared.BacklogBoard,
		tasks:    tasks,
		renderer: renderer,
		conn:     conn,
		inputs:   appInputsMap,
		state:    shared.StateBrowsing,
	}
	app.tasksTable = app.buildTasksTable(tasks)
	return app
}
