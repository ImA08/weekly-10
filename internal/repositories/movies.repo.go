package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"minitask1.go/internal/models"
	"minitask1.go/internal/utils"
)

type MoviesRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewMovieRepository(db *pgxpool.Pool, rdb *redis.Client) *MoviesRepository {
	return &MoviesRepository{db: db, rdb: rdb}
}

// Cache duration constants
const (
	MoviesCacheDuration    = 5 * time.Minute
	UpcomingCacheDuration  = 1 * time.Minute
	CacheKeyAllMovies      = "movies:all:%d"
	CacheKeyUpcomingMovies = "movies:upcoming"
)

// func (r *MoviesRepository) GetAllMovies(ctx context.Context, filterBy models.ShowMovie, page int) ([]models.ShowMovie, error) {
// 	cacheKey := fmt.Sprintf(CacheKeyAllMovies, page)

// 	// 1. Try getting from Redis first
// 	cached, err := r.rdb.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var movies []models.ShowMovie
// 		if err := json.Unmarshal([]byte(cached), &movies); err == nil {
// 			return movies, nil
// 		}
// 	}

// 	// 2. Query from database if cache miss
// 	const limit = 4
// 	offset := (page - 1) * limit

// 	query := `SELECT
// 			id,
// 			title,
// 			poster_path,
// 			genres
// 		FROM (
// 			SELECT
// 				m.id,
// 				m.title,
// 				m.synopsis,
// 				m.poster_path,
// 				ARRAY_AGG(DISTINCT g.genre) AS genres
// 			FROM movies m
// 			LEFT JOIN movie_genres mg ON mg.movie_id = m.id
// 			LEFT JOIN genres g ON mg.genre_id = g.id
// 			GROUP BY m.id
// 			ORDER BY m.id
// 			LIMIT 12 OFFSET $1
// 	) sq WHERE`

// 	values := []any{}
// 	if filterBy.Id > 0 {
// 		query += fmt.Sprintf(`id = $%d`, len(values)+1)
// 		values = append(values, filterBy.Id)
// 	}

// 	if filterBy.Title != "" {
// 		query += fmt.Sprintf(`LOWER(sq.title) LIKE '%%' || LOWER($%d) || '%%'`, len(values)+1)
// 		values = append(values, filterBy.Title)
// 	}

// 	if len(filterBy.Genres) > 0 {
// 		if len(values) > 0 {
// 			query += " AND "
// 		}
// 		query += "("
// 		for i, genre := range filterBy.Genres {
// 			if i > 0 {
// 				query += " OR "
// 			}
// 			query += fmt.Sprintf(`LOWER(ARRAY_TO_STRING(sq.genres, ',')) LIKE '%%' || LOWER($%d) || '%%'`, len(values)+1)
// 			values = append(values, genre)
// 		}
// 		query += ")"
// 	}

// 	rows, err := r.db.Query(ctx, query, values, limit, offset)
// 	if err != nil {
// 		return nil, fmt.Errorf("db query failed: %w", err)
// 	}
// 	defer rows.Close()

// 	var movies []models.ShowMovie
// 	for rows.Next() {
// 		var movie models.ShowMovie
// 		if err := rows.Scan(
// 			&movie.Id,
// 			&movie.Title,
// 			&movie.Image,
// 			&movie.Genres,

// 		); err != nil {
// 			return nil, fmt.Errorf("db scan failed: %w", err)
// 		}
// 		movies = append(movies, movie)
// 	}

// 	// 3. Cache the results
// 	if jsonData, err := json.Marshal(movies); err == nil {
// 		if err := r.rdb.Set(ctx, cacheKey, jsonData, MoviesCacheDuration).Err(); err != nil {
// 			fmt.Printf("Failed to cache movies: %v\n", err)
// 		}
// 	}wok

// 	return movies, nil
// }

func (r *MoviesRepository) GetAllMovies(ctx context.Context, filterBy models.ShowMovie, page int) ([]models.ShowMovie, error) {
	cacheKey := fmt.Sprintf("movies:%d:%v", page, filterBy)

	// 1. Try Redis cache first
	if cached, err := r.rdb.Get(ctx, cacheKey).Result(); err == nil {
		var movies []models.ShowMovie
		if err := json.Unmarshal([]byte(cached), &movies); err == nil {
			return movies, nil
		}
	}

	// 2. Build PostgreSQL query
	const limit = 12
	offset := (page - 1) * limit

	baseQuery := `
		SELECT 
			m.id, 
			m.title, 
			m.poster_path,
			ARRAY_AGG(DISTINCT g.genre) AS genres
		FROM movies m
		LEFT JOIN movie_genres mg ON mg.movie_id = m.id
		LEFT JOIN genres g ON mg.genre_id = g.id`

	var filters []string
	var values []interface{}

	// ID filter
	if filterBy.Id > 0 {
		filters = append(filters, fmt.Sprintf("m.id = $%d", len(values)+1))
		values = append(values, filterBy.Id)
	}

	// Title filter
	if filterBy.Title != "" {
		filters = append(filters, fmt.Sprintf("LOWER(m.title) LIKE '%%' || $%d || '%%'", len(values)+1))
		values = append(values, strings.ToLower(filterBy.Title))
	}

	// Genre filter (using pgx array support)
	if len(filterBy.Genres) > 0 {
		filters = append(filters, fmt.Sprintf(
			`EXISTS (
				SELECT 1 FROM movie_genres mg2
				JOIN genres g2 ON mg2.genre_id = g2.id
				WHERE mg2.movie_id = m.id
				AND LOWER(g2.genre) = ANY($%d)
			)`, len(values)+1))
		values = append(values, utils.LowerCaseStrings(filterBy.Genres))
	}

	// Combine query parts
	query := baseQuery
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += `
		GROUP BY m.id
		ORDER BY m.id
		LIMIT $` + fmt.Sprint(len(values)+1) + `
		OFFSET $` + fmt.Sprint(len(values)+2)
	values = append(values, limit, offset)

	// 3. Execute query (using pgx)
	rows, err := r.db.Query(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("db query failed: %w", err)
	}
	defer rows.Close()

	var movies []models.ShowMovie
	for rows.Next() {
		var movie models.ShowMovie
		if err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.Image,
			&movie.Genres, // pgx automatically scans arrays into slices
		); err != nil {
			return nil, fmt.Errorf("db scan failed: %w", err)
		}
		movies = append(movies, movie)
	}

	// 4. Cache results
	if jsonData, err := json.Marshal(movies); err == nil {
		r.rdb.Set(ctx, cacheKey, jsonData, 24*time.Hour)
	}

	return movies, nil
}

func (r *MoviesRepository) UpcomingMovies(ctx context.Context) ([]models.UpcomingMovie, error) {
	// 1. Try getting from Redis first
	cached, err := r.rdb.Get(ctx, CacheKeyUpcomingMovies).Result()
	if err == nil {
		var movies []models.UpcomingMovie
		if err := json.Unmarshal([]byte(cached), &movies); err == nil {
			return movies, nil
		}
	}

	// 2. Query from database if cache miss
	query := `
        SELECT 
            m.id, m.title, COALESCE(m.synopsis, ''),
            m.duration, c.name, c.location, c.city,
            s.show_time, TO_CHAR(s.show_time, 'YYYY-MM-DD'),
            (
                SELECT JSON_AGG(g.genre)
                FROM movie_genres mg
                JOIN genres g ON mg.genre_id = g.id
                WHERE mg.movie_id = m.id
            ) as genres
        FROM schedules s
        JOIN movies m ON s.movie_id = m.id
        JOIN cinemas c ON s.cinema_id = c.id
        WHERE s.show_time >= CURRENT_DATE
        GROUP BY m.id, c.id, s.show_time
        ORDER BY s.show_time ASC
        LIMIT 15
    `

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db query failed: %w", err)
	}
	defer rows.Close()

	var movies []models.UpcomingMovie
	for rows.Next() {
		var movie models.UpcomingMovie
		var genresJSON []byte

		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Synopsis,
			&movie.Duration,
			&movie.Cinemas,
			&movie.Location,
			&movie.City,
			&movie.Schedule,
			&movie.Date,
			&genresJSON,
		); err != nil {
			return nil, fmt.Errorf("db scan failed: %w", err)
		}

		// Handle genres
		if genresJSON != nil {
			if err := json.Unmarshal(genresJSON, &movie.Genres); err != nil {
				return nil, fmt.Errorf("failed to decode genres: %w", err)
			}
		} else {
			movie.Genres = []string{}
		}

		movies = append(movies, movie)
	}

	// 3. Cache the results
	if jsonData, err := json.Marshal(movies); err == nil {
		if err := r.rdb.Set(ctx, CacheKeyUpcomingMovies, jsonData, UpcomingCacheDuration).Err(); err != nil {
			fmt.Printf("Failed to cache upcoming movies: %v\n", err)
		}
	}

	return movies, nil
}

func (r *MoviesRepository) AddMovies(ctx context.Context, dataMovie models.MovieForm, posterPath, backdropPath string) error {
	// Begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Safe rollback if not committed

	// Insert movie
	var movieId int
	err = tx.QueryRow(
		ctx,
		`INSERT INTO movies (title, synopsis, duration, release_date, poster_path, backdrop_path)
         VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		dataMovie.Title, dataMovie.Synopsis, dataMovie.Duration,
		dataMovie.ReleaseDate, posterPath, backdropPath,
	).Scan(&movieId) // Fixed: using pointer
	if err != nil {
		return fmt.Errorf("failed to insert movie: %w", err)
	}

	// Process casts
	for _, castName := range dataMovie.Casts {
		var castID int
		err := tx.QueryRow(
			ctx,
			`SELECT id FROM casts WHERE name = $1 LIMIT 1`,
			castName,
		).Scan(&castID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = tx.QueryRow(
					ctx,
					`INSERT INTO casts (name) VALUES ($1) RETURNING id`,
					castName,
				).Scan(&castID)
				if err != nil {
					return fmt.Errorf("failed to insert cast %s: %w", castName, err)
				}
			} else {
				return fmt.Errorf("failed to query cast %s: %w", castName, err)
			}
		}

		_, err = tx.Exec(
			ctx,
			`INSERT INTO movie_casts (movie_id, cast_id) VALUES ($1, $2)`,
			movieId, castID,
		)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				continue // Skip duplicate cast-movie links
			}
			return fmt.Errorf("failed to link cast to movie: %w", err)
		}
	}

	// Process directors
	for _, directorName := range dataMovie.Directors {
		var directorID int
		err := tx.QueryRow(
			ctx,
			`SELECT id FROM directors WHERE name = $1 LIMIT 1`,
			directorName,
		).Scan(&directorID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = tx.QueryRow(
					ctx,
					`INSERT INTO directors (name) VALUES ($1) RETURNING id`,
					directorName, // Fixed typo from 'directorName'
				).Scan(&directorID)
				if err != nil {
					return fmt.Errorf("failed to insert director %s: %w", directorName, err)
				}
			} else {
				return fmt.Errorf("failed to query director %s: %w", directorName, err)
			}
		}

		_, err = tx.Exec(
			ctx,
			`INSERT INTO movie_directors (movie_id, director_id) VALUES ($1, $2)`,
			movieId, directorID,
		)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				log.Printf("Director %s already linked to movie %d", directorName, movieId)
				continue
			}
			return fmt.Errorf("failed to link director to movie: %w", err)
		}
	}

	// Insert schedule
	_, err = tx.Exec(
		ctx,
		`INSERT INTO schedules (movie_id, cinema_id, show_time, price, date) 
         VALUES ($1, $2, $3, $4, $5)`,
		movieId, dataMovie.Cinema, dataMovie.ScreenTime,
		dataMovie.Price, dataMovie.Date,
	)
	if err != nil {
		return fmt.Errorf("failed to insert schedule: %w", err)
	}

	// Commit transaction if everything succeeded
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// func (r *MoviesRepository) AddMovies(ctx context.Context, dataMovie models.DetailMovieStruct, posterPath, backdropPath string) error {

// 	queryMovie := `INSERT INTO movies (title, synopsis, duration, release_date, poster_path, backdrop_path)
// 			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
// 	var movieId int
// 	movieValues := []any{dataMovie.Title, dataMovie.Synopsis, dataMovie.Duration, dataMovie.ReleaseDate, posterPath, backdropPath}
// 	if err := r.db.QueryRow(ctx, queryMovie, movieValues...).Scan(movieId); err != nil {
// 		return err
// 	}

// 	// check cast if exist, and insert if itsn't
// 	for _, castName := range dataMovie.Casts {
// 		var castID int

// 		// Check if cast exists
// 		err := r.db.QueryRow(ctx, `SELECT id FROM casts WHERE name = $1 LIMIT 1`, castName).Scan(&castID)

// 		if err != nil {
// 			if errors.Is(err, pgx.ErrNoRows) {
// 				// Insert new cast member
// 				err = r.db.QueryRow(
// 					ctx,
// 					`INSERT INTO casts (name) VALUES ($1) RETURNING id`,
// 					castName,
// 				).Scan(&castID)

// 				if err != nil {
// 					// return fmt.Errorf("failed to insert cast %s: %w", castName, err)
// 					return err
// 				}
// 			} else {
// 				// return fmt.Errorf("failed to query cast %s: %w", castName, err)
// 				return err
// 			}
// 		}

// 		// Link cast to movie
// 		_, err = r.db.Exec(
// 			ctx,
// 			`INSERT INTO movie_casts (movie_id, cast_id) VALUES ($1, $2)`,
// 			movieId,
// 			castID,
// 		)
// 		if err != nil {
// 			var pgErr *pgconn.PgError
// 			if errors.As(err, &pgErr) {
// 				if pgErr.Code == "23505" {
// 					continue
// 				}
// 			}
// 			return err
// 		}
// 	}

// 	// check director if exist, and insert if itsn't
// 	for _, directorName := range dataMovie.Directors {
// 		var directorID int

// 		// Check if director exists
// 		err := r.db.QueryRow(ctx, `SELECT id FROM directors WHERE name = $1 LIMIT 1`, directorName).Scan(&directorID)

// 		if err != nil {
// 			if errors.Is(err, pgx.ErrNoRows) {
// 				// Insert new director
// 				err = r.db.QueryRow(
// 					ctx,
// 					`INSERT INTO directors (name) VALUES ($1) RETURNING id`,
// 					directorName,
// 				).Scan(&directorID)

// 				if err != nil {
// 					return fmt.Errorf("failed to insert director %s: %w", directorName, err)
// 				}
// 			} else {
// 				return fmt.Errorf("failed to query director %s: %w", directorName, err)
// 			}
// 		}

// 		// Link director to movie
// 		_, err = r.db.Exec(
// 			ctx,
// 			`INSERT INTO movie_directors (movie_id, director_id) VALUES ($1, $2)`,
// 			movieId,
// 			directorID,
// 		)
// 		if err != nil {
// 			var pgErr *pgconn.PgError
// 			if errors.As(err, &pgErr) {
// 				if pgErr.Code == "23505" {
// 					log.Printf("Director %s already linked to movie %d", directorName, movieId)
// 					continue
// 				}
// 			}
// 			return fmt.Errorf("failed to link director to movie: %w", err)
// 		}
// 	}

// 	querySchedule := `INSERT INTO schedules (movie_id, cinema_id, show_time, price, date) VALUES ($1, $2, $3, $4, $5)`
// 	scheduleValues := []any{movieId, dataMovie.Cinema, dataMovie.ScreenTime, dataMovie.Price, dataMovie.Date}
// 	_, err := r.db.Exec(ctx, querySchedule, scheduleValues...)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// InvalidateCache clears relevant movie caches
func (r *MoviesRepository) InvalidateCache(ctx context.Context) error {
	// Invalidate all movies cache
	keys, err := r.rdb.Keys(ctx, "movies:*").Result()
	if err != nil {
		return fmt.Errorf("failed to get cache keys: %w", err)
	}

	if len(keys) > 0 {
		if err := r.rdb.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete cache keys: %w", err)
		}
	}
	return nil
}
