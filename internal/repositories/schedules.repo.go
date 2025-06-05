package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"minitask1.go/internal/models"
)

func (r *MoviesRepository) GetScheduleMovieRepo(ctx context.Context, movieId int) ([]models.Schedules, error) {
	cacheKey := fmt.Sprintf("movie:%d:schedules", movieId)

	if cached, err := r.rdb.Get(ctx, cacheKey).Result(); err == nil {
		var schedules []models.Schedules
		if err := json.Unmarshal([]byte(cached), &schedules); err == nil {
			return schedules, nil
		}
		log.Printf("cache unmarshal error: %v", err)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `SELECT 
    c.name, c.location, c.city,
    s.show_time, s.price, s.date
    FROM schedules s LEFT JOIN cinemas c
    ON s.cinema_id = c.id
    WHERE s.movie_id = $1
    ORDER BY s.date, s.show_time`

	rows, err := tx.Query(ctx, query, movieId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var schedules []models.Schedules
	for rows.Next() {
		var schedule models.Schedules
		if err := rows.Scan(
			&schedule.Cinema,
			&schedule.Location,
			&schedule.City,
			&schedule.ShowTime,
			&schedule.Price,
			&schedule.Date,
		); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	if len(schedules) > 0 {
		if jsonData, err := json.Marshal(schedules); err == nil {
			if err := r.rdb.SetEx(ctx, cacheKey, jsonData, 24*time.Hour).Err(); err != nil {
				log.Printf("failed to set cache: %v", err)
			}
		}
	}

	return schedules, nil
}
