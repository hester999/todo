package repo

import (
	"github.com/jmoiron/sqlx"
	"time"
	"todo/internal/apperr"
	"todo/internal/db"
	"todo/internal/entity"
)

type TaskRepositoryImpl struct {
	cursor *sqlx.DB
}

func NewTaskRepositoryImpl() *TaskRepositoryImpl {
	cur, err := db.Connection()
	if err != nil {
		panic(err)
	}

	return &TaskRepositoryImpl{cursor: cur}
}

func (t *TaskRepositoryImpl) Save(task entity.Task) (entity.Task, error) {
	queryString := "INSERT INTO tasks (id, title, description, status, create_time) VALUES ($1, $2, $3, $4, $5) RETURNING id,title,description,status,create_time;"

	savedTask := struct {
		Id          string    `db:"id"`
		Title       string    `db:"title"`
		Description string    `db:"description"`
		Status      string    `db:"status"`
		CreateTime  time.Time `db:"create_time"`
	}(task)

	err := t.cursor.Get(&savedTask, queryString, task.Id, task.Title, task.Description, task.Status, task.CreateTime)

	if err != nil {

		return entity.Task{}, apperr.DatabaseError
	}

	return entity.Task{savedTask.Id, savedTask.Title, savedTask.Description, savedTask.Status, savedTask.CreateTime}, nil

}
func (t *TaskRepositoryImpl) FindAll() ([]entity.Task, error) {
	query := "SELECT id, title, description, status, create_time FROM tasks"
	// TODO sqlx есть SelectContext

	var res []struct {
		Id          string    `db:"id"`
		Title       string    `db:"title"`
		Description string    `db:"description"`
		Status      string    `db:"status"`
		CreateTime  time.Time `db:"create_time"`
	}

	err := t.cursor.Select(&res, query)

	if err != nil {
		return nil, apperr.DatabaseError
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
		Id          string    `db:"id"`
		Title       string    `db:"title"`
		Description string    `db:"description"`
		Status      string    `db:"status"`
		CreateTime  time.Time `db:"create_time"`
	}

	err := t.cursor.Get(&res, query, id)

	if err != nil {
		return entity.Task{}, apperr.NotFoundError
	}

	return entity.Task{res.Id, res.Title, res.Description, res.Status, res.CreateTime}, nil
}

func (t *TaskRepositoryImpl) DeleteById(id string) error {
	queryString := "DELETE FROM tasks WHERE id = $1"
	_, err := t.cursor.Exec(queryString, id)
	if err != nil {
		return apperr.DatabaseError
	}
	return nil
}

func (t *TaskRepositoryImpl) DeleteAll() error {
	queryString := "DELETE FROM tasks"
	_, err := t.cursor.Exec(queryString)
	if err != nil {
		return apperr.DatabaseError
	}
	return nil
}

func (t *TaskRepositoryImpl) Edit(id string, task entity.Task) (entity.Task, error) {
	queryString := "UPDATE tasks SET title = $1, description = $2, status = $3, create_time = $4 WHERE id = $5 RETURNING id, title, description, status, create_time"

	updateTask := struct {
		Id          string    `db:"id"`
		Title       string    `db:"title"`
		Description string    `db:"description"`
		Status      string    `db:"status"`
		CreateTime  time.Time `db:"create_time"`
	}{}

	err := t.cursor.Get(&updateTask, queryString, task.Title, task.Description, task.Status, task.CreateTime, id)
	if err != nil {
		return entity.Task{}, apperr.DatabaseError
	}
	return entity.Task{updateTask.Id, updateTask.Title, updateTask.Description, updateTask.Status, updateTask.CreateTime}, nil
}
