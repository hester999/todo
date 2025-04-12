package entity

import "time"

type Task struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	CreateTime  time.Time `json:"createTime"`
}
