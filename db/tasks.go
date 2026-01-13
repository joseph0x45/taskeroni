package db

import (
	"fmt"

	"github.com/joseph0x45/taskeroni/internal/models"
)

func (c *Conn) GetTasks() ([]models.Task, error) {
	tasks := []models.Task{}
	const query = "select * from tasks"
	err := c.db.Select(&tasks, query)
	if err != nil {
		return nil, fmt.Errorf("Error while getting tasks: %w", err)
	}
	return tasks, nil
}
