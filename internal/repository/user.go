package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateWithPassword(ctx context.Context, name, email, passwordHash string) (*User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (name, email, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, email, role, created_at`,
		name, email, passwordHash,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*User, string, error) {
	var u User
	var hash string
	err := r.db.QueryRow(ctx,
		`SELECT id, name, email, role, created_at, password_hash
		 FROM users
		 WHERE email = $1 AND deleted_at IS NULL`,
		email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, "", fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, "", fmt.Errorf("get user by email: %w", err)
	}
	return &u, hash, nil
}

func (r *UserRepo) Create(ctx context.Context, name, email string) (*User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (name, email)
		 VALUES ($1, $2)
		 RETURNING id, name, email, role, created_at`,
		name, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetAll(ctx context.Context) ([]*User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, email, role, created_at
		 FROM users
		 WHERE deleted_at IS NULL
		 ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		`SELECT id, name, email, role, created_at
		 FROM users
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) Update(ctx context.Context, id int, name, email string) (*User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		`UPDATE users
		 SET name = $1, email = $2
		 WHERE id = $3 AND deleted_at IS NULL
		 RETURNING id, name, email, role, created_at`,
		name, email, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) SoftDelete(ctx context.Context, id int) error {
	result, err := r.db.Exec(ctx,
		`UPDATE users SET deleted_at = NOW()
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("soft delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
