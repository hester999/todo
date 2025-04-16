package entity

import "time"

type TaskStatus string

const (
	InProgress TaskStatus = "In Progress"
	Finished   TaskStatus = "Finished"
)

type Task struct {
	Id          string
	Title       string
	Description string
	Status      TaskStatus
	CreateTime  time.Time
}
