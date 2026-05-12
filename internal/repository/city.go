package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"weather-api/internal/model"
)

type CityRepo struct {
	db *pgxpool.Pool
}

func NewCityRepo(db *pgxpool.Pool) *CityRepo {
	return &CityRepo{db: db}
}

func (r *CityRepo) Add(ctx context.Context, userID int, city string) (*model.City, error) {
	var c model.City
	err := r.db.QueryRow(ctx,
		`INSERT INTO user_cities (user_id, city)
		 VALUES ($1, $2)
		 ON CONFLICT (user_id, city) DO UPDATE SET city = EXCLUDED.city
		 RETURNING id, user_id, city`,
		userID, city,
	).Scan(&c.ID, &c.UserID, &c.City)
	if err != nil {
		return nil, fmt.Errorf("add city: %w", err)
	}
	return &c, nil
}

func (r *CityRepo) GetByUser(ctx context.Context, userID int) ([]*model.City, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, city
		 FROM user_cities
		 WHERE user_id = $1
		 ORDER BY id`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get cities: %w", err)
	}
	defer rows.Close()

	var cities []*model.City
	for rows.Next() {
		var c model.City
		if err := rows.Scan(&c.ID, &c.UserID, &c.City); err != nil {
			return nil, fmt.Errorf("scan city: %w", err)
		}
		cities = append(cities, &c)
	}
	return cities, nil
}

func (r *CityRepo) Delete(ctx context.Context, userID, cityID int) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM user_cities WHERE id = $1 AND user_id = $2`,
		cityID, userID,
	)
	if err != nil {
		return fmt.Errorf("delete city: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("city not found")
	}
	return nil
}
