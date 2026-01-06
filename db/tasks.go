package db

import (
	"fmt"

	"github.com/joseph0x45/taskeroni/internal/models"
)

func (c *Conn) GetUserTasks(userID string) ([]models.Task, error) {
	//TODO: Read more about usage of make vs raw slice here
	tasks := []models.Task{}
	const query = `
    select * from tasks where owner_id=?
  `
	err := c.db.Select(&tasks, query, userID)
	if err != nil {
		return nil, fmt.Errorf("Error while getting user tasks: %w", err)
	}
	return tasks, nil
}
