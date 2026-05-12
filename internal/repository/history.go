package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"weather-api/internal/model"
)

type HistoryRepo struct {
	db *pgxpool.Pool
}

func NewHistoryRepo(db *pgxpool.Pool) *HistoryRepo {
	return &HistoryRepo{db: db}
}

func (r *HistoryRepo) Save(ctx context.Context, userID int, city string, temperature float64, description string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO weather_history (user_id, city, temperature, description)
		 VALUES ($1, $2, $3, $4)`,
		userID, city, temperature, description,
	)
	if err != nil {
		return fmt.Errorf("save history: %w", err)
	}
	return nil
}

func (r *HistoryRepo) Get(ctx context.Context, userID int, f model.HistoryFilter) ([]*model.WeatherHistory, error) {
	args := []any{userID}

	query := `SELECT id, user_id, city, temperature, description, requested_at
	          FROM weather_history
	          WHERE user_id = $1`

	if f.City != "" {
		args = append(args, f.City)
		query += fmt.Sprintf(" AND city = $%d", len(args))
	}

	query += " ORDER BY requested_at DESC"

	if f.Limit > 0 {
		args = append(args, f.Limit)
		query += fmt.Sprintf(" LIMIT $%d", len(args))
	}

	if f.Offset > 0 {
		args = append(args, f.Offset)
		query += fmt.Sprintf(" OFFSET $%d", len(args))
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}
	defer rows.Close()

	results := make([]*model.WeatherHistory, 0)
	for rows.Next() {
		var h model.WeatherHistory
		if err := rows.Scan(&h.ID, &h.UserID, &h.City, &h.Temperature, &h.Description, &h.RequestedAt); err != nil {
			return nil, fmt.Errorf("scan history: %w", err)
		}
		results = append(results, &h)
	}
	return results, nil
}
