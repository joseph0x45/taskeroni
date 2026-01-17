package models

type TaskStatus string

type Task struct {
	ID    string `json:"id" db:"id"`
	Label string `json:"label" db:"label"`
	//Status can be backlog, todo, in progress or done
	Status  string `json:"status" db:"status"`
	DueDate string `json:"due_date" db:"due_date"`
}
