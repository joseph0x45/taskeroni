package tui

import "fmt"

func (app App) renderTasks() string {
	ui := fmt.Sprintf("Showing tasks in %s\n", app.board)
	ui += baseStyle.Render(app.tasksTable.View()) + "\n"
	return ui
}

func (app App) renderTaskCreationForm() string {
	return "creating tasks"
}
