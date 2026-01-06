package tui

import "fmt"

func (app App) refreshTasksList() error {
	return nil
}

func (app App) renderTasksScreen() string {
	ui := fmt.Sprintf("Showing tasks in %s\n", app.displayedTasks)
	ui += baseStyle.Render(app.tasksTable.View()) + "\n"
	return ui
}
