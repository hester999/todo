package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
	"todo/internal/apperr"
	"todo/internal/entity"
)

type TaskRepositoryImpl struct {
	db *sqlx.DB
}

func NewTaskRepositoryImpl(db *sqlx.DB) *TaskRepositoryImpl {

	return &TaskRepositoryImpl{db: db}
}

func (t *TaskRepositoryImpl) Save(task entity.Task) (entity.Task, error) {
	queryString := "INSERT INTO tasks (id, title, description, status, create_time) VALUES ($1, $2, $3, $4, $5) RETURNING id,title,description,status,create_time;"

	savedTask := struct {
		Id          string            `db:"id"`
		Title       string            `db:"title"`
		Description string            `db:"description"`
		Status      entity.TaskStatus `db:"status"`
		CreateTime  time.Time         `db:"create_time"`
	}(task)

	err := t.db.Get(&savedTask, queryString, task.Id, task.Title, task.Description, task.Status, task.CreateTime)

	if err != nil {

		return entity.Task{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return entity.Task{savedTask.Id, savedTask.Title, savedTask.Description, savedTask.Status, savedTask.CreateTime}, nil

}
func (t *TaskRepositoryImpl) FindAll() ([]entity.Task, error) {
	query := "SELECT id, title, description, status, create_time FROM tasks"
	// TODO sqlx есть SelectContext

	var res []struct {
		Id          string            `db:"id"`
		Title       string            `db:"title"`
		Description string            `db:"description"`
		Status      entity.TaskStatus `db:"status"`
		CreateTime  time.Time         `db:"create_time"`
	}

	err := t.db.Select(&res, query)

	if err != nil {
		return nil, err
	}
	result := make([]entity.Task, 0, len(res))

	for _, row := range res {
		result = append(result, entity.Task{row.Id, row.Title, row.Description, row.Status, row.CreateTime})
	}

	return result, nil
}

func (t *TaskRepositoryImpl) FindById(id string) (entity.Task, error) {
	query := "SELECT id, title, description, status, create_time FROM tasks WHERE id = $1"

	var res struct {
		Id          string            `db:"id"`
		Title       string            `db:"title"`
		Description string            `db:"description"`
		Status      entity.TaskStatus `db:"status"`
		CreateTime  time.Time         `db:"create_time"`
	}

	err := t.db.Get(&res, query, id)

	if err != nil {
		return entity.Task{}, apperr.NotFoundError
	}

	return entity.Task{res.Id, res.Title, res.Description, res.Status, res.CreateTime}, nil
}

func (t *TaskRepositoryImpl) DeleteById(id string) error {
	queryString := "DELETE FROM tasks WHERE id = $1"
	_, err := t.db.Exec(queryString, id)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (t *TaskRepositoryImpl) DeleteAll() error {
	queryString := "DELETE FROM tasks"
	_, err := t.db.Exec(queryString)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (t *TaskRepositoryImpl) Edit(id string, task entity.Task) (entity.Task, error) {
	queryString := "UPDATE tasks SET title = $1, description = $2, status = $3, create_time = $4 WHERE id = $5 RETURNING id, title, description, status, create_time"

	updateTask := struct {
		Id          string            `db:"id"`
		Title       string            `db:"title"`
		Description string            `db:"description"`
		Status      entity.TaskStatus `db:"status"`
		CreateTime  time.Time         `db:"create_time"`
	}{}

	err := t.db.Get(&updateTask, queryString, task.Title, task.Description, task.Status, task.CreateTime, id)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to execute query: %w", err)
	}
	return entity.Task{updateTask.Id, updateTask.Title, updateTask.Description, updateTask.Status, updateTask.CreateTime}, nil
}
