package adapter

import "todo/internal/entity"

type TaskService interface {
	CreateTask(title, description string) (entity.Task, error)
	GetAllTasks() ([]entity.Task, error)
	GetTaskById(taskId string) (entity.Task, error)
	DeleteTaskById(taskId string) error
	DeleteAllTasks() error
	UpdateTask(id string, task entity.Task) (entity.Task, error)
}

// вот тут хз, нужноо ли это но по идее оно доолжно адаптировать логику под тг бота, и не срать в фаайле бота
type BotAdapter struct {
	taskService TaskService
}

func NewBotAdapter(taskService TaskService) *BotAdapter {
	return &BotAdapter{taskService: taskService}
}

func (adapter *BotAdapter) CreateTask(title, description string) (entity.Task, error) {
	return adapter.taskService.CreateTask(title, description)
}

func (adapter *BotAdapter) GetAllTasks() ([]entity.Task, error) {
	return adapter.taskService.GetAllTasks()
}

func (adapter *BotAdapter) GetTaskById(taskId string) (entity.Task, error) {
	return adapter.taskService.GetTaskById(taskId)
}

func (adapter *BotAdapter) DeleteTaskById(taskId string) error {
	return adapter.taskService.DeleteTaskById(taskId)
}

func (adapter *BotAdapter) DeleteAllTasks() error {
	return adapter.taskService.DeleteAllTasks()
}

func (adapter *BotAdapter) UpdateTask(id string, task entity.Task) (entity.Task, error) {
	return adapter.taskService.UpdateTask(id, task)
}
