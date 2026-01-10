package database

import (
	"REST_1/internal/models"

	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (ts *TaskStore) GetAll() ([]models.Task, error) {
	tasks := make([]models.Task, 0)
	query := `SELECT id, title, description, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC;`
	err := ts.db.Select(&tasks, query)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ts *TaskStore) GetByID(id int) (*models.Task, error) {
	var task models.Task
	query := `SELECT id, title, description, created_at, updated_at
		FROM tasks
		WHERE id = $1;`
	err := ts.db.Get(&task, query, id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
