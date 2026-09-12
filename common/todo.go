// Package common contains code deliberately shared by every API implementation.
// Keeping it here makes the framework folders compare only HTTP-framework choices.
package common

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Todo struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Completed   bool      `json:"completed"`
	DueDate     *string   `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TodoInput struct {
	UserID      int64   `json:"user_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Completed   bool    `json:"completed"`
	DueDate     *string `json:"due_date"`
}

type Store struct{ db *sql.DB }

const fields = "id, user_id, title, description, completed, due_date, created_at, updated_at"

func Env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func OpenStore() (*Store, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", Env("DB_USER", "admin"), Env("DB_PASSWORD", "jamurkembang"), Env("DB_HOST", "127.0.0.1"), Env("DB_PORT", "3306"), Env("DB_NAME", "todos"))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error                   { return s.db.Close() }
func (s *Store) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }

func Validate(input TodoInput) error {
	if input.UserID < 1 || input.Title == "" {
		return errors.New("user_id and title are required")
	}
	return nil
}

func scan(row interface{ Scan(...any) error }) (Todo, error) {
	var todo Todo
	err := row.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Completed, &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt)
	return todo, err
}

func (s *Store) List(ctx context.Context, userID string) ([]Todo, error) {
	query, args := "SELECT "+fields+" FROM todos", []any{}
	if userID != "" {
		query += " WHERE user_id = ?"
		args = append(args, userID)
	}
	rows, err := s.db.QueryContext(ctx, query+" ORDER BY id DESC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todos := []Todo{}
	for rows.Next() {
		todo, err := scan(rows)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, rows.Err()
}

func (s *Store) Get(ctx context.Context, id string) (Todo, error) {
	return scan(s.db.QueryRowContext(ctx, "SELECT "+fields+" FROM todos WHERE id = ?", id))
}

func (s *Store) Create(ctx context.Context, input TodoInput) (Todo, error) {
	result, err := s.db.ExecContext(ctx, "INSERT INTO todos (user_id, title, description, completed, due_date) VALUES (?, ?, ?, ?, ?)", input.UserID, input.Title, input.Description, input.Completed, input.DueDate)
	if err != nil {
		return Todo{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Todo{}, err
	}
	return s.Get(ctx, fmt.Sprint(id))
}

func (s *Store) Update(ctx context.Context, id string, input TodoInput) (Todo, error) {
	result, err := s.db.ExecContext(ctx, "UPDATE todos SET user_id=?, title=?, description=?, completed=?, due_date=? WHERE id=?", input.UserID, input.Title, input.Description, input.Completed, input.DueDate, id)
	if err != nil {
		return Todo{}, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return Todo{}, err
	}
	if count == 0 {
		return Todo{}, sql.ErrNoRows
	}
	return s.Get(ctx, id)
}

func (s *Store) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM todos WHERE id=?", id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
