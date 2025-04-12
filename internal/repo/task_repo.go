package repo

import (
	"fmt"
	"todo/internal/db"
	"todo/internal/entity"
)

type TaskRepository interface {
	Save(task entity.Task) (entity.Task, error)
	FindAll() ([]entity.Task, error)
	FindById(id string) (entity.Task, error)
	DeleteById(id string) error
	DeleteAll() error
	Edit(id string, task entity.Task) (entity.Task, error)
}

type TaskRepositoryImpl struct {
	db db.Database
}

func NewTaskRepositoryImpl(db db.Database) *TaskRepositoryImpl {
	if db == nil {
		panic("db cannot be nil")
	}
	return &TaskRepositoryImpl{db: db}
}

func (t *TaskRepositoryImpl) Save(task entity.Task) (entity.Task, error) {
	queryString := "INSERT INTO tasks (id, title, description, status, create_time) VALUES ($1, $2, $3, $4, $5) RETURNING id,title,description,status,create_time;"

	var savedTask entity.Task
	row, err := t.db.Query(queryString, task.Id, task.Title, task.Description, task.Status, task.CreateTime)

	if err != nil {
		return entity.Task{}, err
	}
	defer row.Close()

	if row.Next() {
		err := row.Scan(&savedTask.Id, &savedTask.Title, &savedTask.Description, &savedTask.Status, &savedTask.CreateTime)
		if err != nil {
			return entity.Task{}, err
		}
	} else {
		return entity.Task{}, fmt.Errorf("failed to retrieve saved task")
	}
	return savedTask, nil

}

func (t *TaskRepositoryImpl) FindAll() ([]entity.Task, error) {
	query := "SELECT id, title, description, status, create_time FROM tasks"
	rows, err := t.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreateTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (t *TaskRepositoryImpl) FindById(id string) (entity.Task, error) {
	query := "SELECT id, title, description, status, create_time FROM tasks WHERE id = $1"
	rows, err := t.db.Query(query, id)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to find task: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var task entity.Task
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreateTime)
		if err != nil {
			return entity.Task{}, fmt.Errorf("failed to scan task: %w", err)
		}
		return task, nil
	}
	return entity.Task{}, fmt.Errorf("task with id %s not found", id)
}

func (t *TaskRepositoryImpl) DeleteById(id string) error {
	queryString := "DELETE FROM tasks WHERE id = $1"
	err := t.db.Exec(queryString, id)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskRepositoryImpl) DeleteAll() error {
	queryString := "DELETE FROM tasks"
	err := t.db.Exec(queryString)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskRepositoryImpl) Edit(id string, task entity.Task) (entity.Task, error) {
	queryString := "UPDATE tasks SET title = $1, description = $2, status = $3, create_time = $4 WHERE id = $5 RETURNING id, title, description, status, create_time"
	var updatedTask entity.Task
	rows, err := t.db.Query(queryString, task.Title, task.Description, task.Status, task.CreateTime, id)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to update task: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&updatedTask.Id, &updatedTask.Title, &updatedTask.Description, &updatedTask.Status, &updatedTask.CreateTime)
		if err != nil {
			return entity.Task{}, fmt.Errorf("failed to scan updated task: %w", err)
		}
		return updatedTask, nil
	}
	return entity.Task{}, fmt.Errorf("task with id %s not found", id)
}
