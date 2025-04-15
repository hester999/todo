package entity

import "time"

type Task struct {
	Id          string
	Title       string
	Description string
	Status      string
	CreateTime  time.Time
}
