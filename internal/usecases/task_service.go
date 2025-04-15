package usecases

import (
	"fmt"
	"time"
	"todo/internal/entity"
	"todo/internal/utils"
)

type TaskRepository interface {
	Save(task entity.Task) (entity.Task, error)
	FindAll() ([]entity.Task, error)
	FindById(id string) (entity.Task, error)
	DeleteById(id string) error
	DeleteAll() error
	Edit(id string, task entity.Task) (entity.Task, error)
}

type TaskServiceImpl struct {
	repo TaskRepository
}

func NewTaskServiceImpl(repo TaskRepository) *TaskServiceImpl {
	return &TaskServiceImpl{repo: repo}
}

func (t *TaskServiceImpl) CreateTask(title, description string) (entity.Task, error) {
	id, err := utils.GenerateUUID()
	if err != nil {
		return entity.Task{}, fmt.Errorf("error generating id: %w", err)
	}

	task := entity.Task{
		Id:          id,
		Title:       title,
		Description: description,
		Status:      "in progress",
		CreateTime:  time.Now(),
	}
	savedTask, err := t.repo.Save(task)
	if err != nil {
		return entity.Task{}, fmt.Errorf("error saving task: %w", err)
	}

	return savedTask, nil
}

func (t *TaskServiceImpl) GetAllTasks() ([]entity.Task, error) {
	tasks, err := t.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error getting all tasks: %w", err)
	}
	return tasks, nil
}

func (t *TaskServiceImpl) GetTaskById(taskId string) (entity.Task, error) {
	if err := utils.UUIDValidator(taskId); err != nil {
		return entity.Task{}, fmt.Errorf("error getting task by id: %w", err)
	}
	task, err := t.repo.FindById(taskId)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to get task by id: %w", err)
	}
	return task, nil
}

func (t *TaskServiceImpl) DeleteTaskById(taskId string) error {
	_, err := t.repo.FindById(taskId)
	if err != nil {
		return fmt.Errorf("error getting task by id: %w", err)
	}

	if err := t.repo.DeleteById(taskId); err != nil {
		return fmt.Errorf("failed to delete task %w", err)
	}
	return nil
}

func (t *TaskServiceImpl) DeleteAllTasks() error {
	if err := t.repo.DeleteAll(); err != nil {
		return fmt.Errorf("failed to delete all tasks %w", err)
	}
	return nil
}

func (t *TaskServiceImpl) UpdateTask(id string, task entity.Task) (entity.Task, error) {
	if err := utils.UUIDValidator(id); err != nil {
		return entity.Task{}, fmt.Errorf("invalid task id %w", err)
	}

	task.Id = id
	_, err := t.repo.FindById(id)
	if err != nil {
		return entity.Task{}, fmt.Errorf("error getting task by id: %w", err)
	}

	updatedTask, err := t.repo.Edit(id, task)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to update task %w", err)
	}
	return updatedTask, nil
}
