package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailTaken = errors.New("email already registered")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash string) (*User, error) {
	var u User
	query := `
		INSERT INTO users (email, password, created_at, modified_at)
		VALUES ($1, $2, now(), now())
		RETURNING id, email, password, created_at, modified_at
	`
	err := r.pool.QueryRow(ctx, query, email, passwordHash).Scan(
		&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.ModifiedAt,
	)

	if err != nil {
		var pgErr interface{ ConstraintName() string }

		if errors.As(err, &pgErr) && pgErr.ConstraintName() == "users_email_key" {
			return nil, ErrEmailTaken
		}

		return nil, err
	}
	return &u, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	query := `
		SELECT
			id
			, email
			, password
			, created_at
			, modified_at
		FROM users
		WHERE email = $1
	`
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.ModifiedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}
