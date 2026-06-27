package users

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	List(ctx context.Context) ([]User, error)
}

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) SQLiteRepository {
	return SQLiteRepository{db: db}
}

func (r SQLiteRepository) List(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, name, email, role
FROM users
ORDER BY id;`)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	result := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		result = append(result, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return result, nil
}
