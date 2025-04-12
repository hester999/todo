package usecases

import (
	"fmt"
	"strings"
	"time"
	"todo/internal/entity"
	"todo/internal/errors"
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
		return entity.Task{}, errors.NewInternal("failed to generate UUID", err)
	}
	if err := utils.ValidateTitle(title, 30); err != nil {
		return entity.Task{}, errors.NewBadRequest("invalid title: "+err.Error(), nil)
	}
	if err := utils.ValidateDescription(description, 500); err != nil {
		return entity.Task{}, errors.NewBadRequest("invalid description: "+err.Error(), nil)
	}
	task := entity.Task{
		Id:          id,
		Title:       title,
		Description: description,
		Status:      false,
		CreateTime:  time.Now(),
	}
	savedTask, err := t.repo.Save(task)
	if err != nil {
		return entity.Task{}, errors.NewInternal("failed to save task", err)
	}

	return savedTask, nil
}

func (t *TaskServiceImpl) GetAllTasks() ([]entity.Task, error) {
	tasks, err := t.repo.FindAll()
	if err != nil {
		return nil, errors.NewInternal("failed to get all tasks", err)
	}
	return tasks, nil
}

func (t *TaskServiceImpl) GetTaskById(taskId string) (entity.Task, error) {
	if err := utils.UUIDValidator(taskId); err != nil {
		return entity.Task{}, errors.NewBadRequest("invalid task id", err)
	}
	task, err := t.repo.FindById(taskId)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to get task by id: %w", err)
	}
	return task, nil
}

func (t *TaskServiceImpl) DeleteTaskById(taskId string) error {
	if err := utils.UUIDValidator(taskId); err != nil {
		return errors.NewBadRequest("invalid task id", err)
	}
	_, err := t.repo.FindById(taskId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.NewNotFound("task not found", err)
		}
		return errors.NewInternal("failed to check task existence", err)
	}
	if err := t.repo.DeleteById(taskId); err != nil {
		return errors.NewInternal("failed to delete task", err)
	}
	return nil
}

func (t *TaskServiceImpl) DeleteAllTasks() error {
	if err := t.repo.DeleteAll(); err != nil {
		return errors.NewInternal("failed to delete all tasks", err)
	}
	return nil
}

func (t *TaskServiceImpl) UpdateTask(id string, task entity.Task) (entity.Task, error) {
	if err := utils.UUIDValidator(id); err != nil {
		return entity.Task{}, errors.NewBadRequest("invalid task id", err)
	}
	if err := utils.ValidateDescription(task.Description, 100); err != nil {
		return entity.Task{}, errors.NewBadRequest("invalid description: "+err.Error(), nil)
	}
	if err := utils.ValidateTitle(task.Title, 30); err != nil {
		return entity.Task{}, errors.NewBadRequest("invalid title: "+err.Error(), nil)
	}
	task.Id = id
	_, err := t.repo.FindById(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return entity.Task{}, errors.NewNotFound("task not found", err)
		}
		return entity.Task{}, errors.NewInternal("failed to check task existence", err)
	}

	updatedTask, err := t.repo.Edit(id, task)
	if err != nil {
		return entity.Task{}, errors.NewInternal("failed to update task", err)
	}
	return updatedTask, nil
}
