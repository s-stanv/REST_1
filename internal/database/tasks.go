package database

import (
	"REST_1/internal/models"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db}
}

func (ts *TaskStore) GetAll() ([]models.Task, error) {
	tasks := make([]models.Task, 0)
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
		;`
	err := ts.db.Select(&tasks, query)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ts *TaskStore) GetByID(id int) (*models.Task, error) {
	var task models.Task
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM tasks
		WHERE id = $1
		;`
	err := ts.db.Get(&task, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf(`task with id %d not found`, id)
	}
	log.Println(task)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (ts *TaskStore) Create(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task
	query := `
		INSERT INTO tasks (title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
		;`
	now := time.Now()
	err := ts.db.QueryRowx(query, input.Title, input.Description, input.Completed, now, now).StructScan(&task)

	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (ts *TaskStore) Update(id int, input models.UpdateTaskInput) (*models.Task, error) {
	task, err := ts.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}
	task.UpdatedAt = time.Now()
	query := `
		UPDATE tasks
		SET title = $1, description = $2, completed = $3, updated_at = $4
		WHERE ID = $5
		RETURNING *
	;`
	var updatedTask models.Task
	err = ts.db.QueryRowx(query, task.Title, task.Description, task.Completed, task.UpdatedAt, id).StructScan(&updatedTask)
	if err != nil {
		return nil, err
	}
	return &updatedTask, nil
}

func (ts *TaskStore) Delete(id int) error {
	query := `
		DELETE FROM tasks WHERE id = $1
	;`
	result, err := ts.db.Exec(query, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}
