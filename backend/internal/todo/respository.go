package todo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTodoNotFound = errors.New("todo not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, userID int64, title string) (*Todo, error) {
	var t Todo
	query := `
		INSERT INTO todos (user_id, title, is_done, created_at, modified_at)
		VALUES ($1, $2, false, now(), now())
		RETURNING id, user_id, title, is_done, created_at, modified_at
	`
	err := r.pool.QueryRow(ctx, query, userID, title).Scan(
		&t.ID, &t.UserID, &t.Title, &t.IsDone, &t.CreatedAt, &t.ModifiedAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *Repository) ListByUser(ctx context.Context, userID int64) ([]Todo, error) {
	query := `
		SELECT id, user_id, title, is_done, created_at, modified_at
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo

		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.IsDone, &t.CreatedAt, &t.ModifiedAt); err != nil {
			return nil, err
		}

		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *Repository) Update(ctx context.Context, id, userID int64, title *string, isDone *bool) (*Todo, error) {
	var t Todo
	query := `
		UPDATE todos
		SET title = COALESCE($1, title),
		    is_done = COALESCE($2, is_done),
		    modified_at = now()
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, title, is_done, created_at, modified_at
	`
	err := r.pool.QueryRow(ctx, query, title, isDone, id, userID).Scan(
		&t.ID, &t.UserID, &t.Title, &t.IsDone, &t.CreatedAt, &t.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) Delete(ctx context.Context, id, userID int64) error {
	query := `
		DELETE FROM todos 
		WHERE id = $1 
			AND user_id = $2
	`
	tag, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrTodoNotFound
	}

	return nil
}
