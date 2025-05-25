package repositories

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"minitask1.go/internal/models"
)

type UserRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewUserRepository(db *pgxpool.Pool, rdb *redis.Client) *UserRepository {
	return &UserRepository{db: db, rdb: rdb}
}

func (r *UserRepository) LogInUserRepo(ctx context.Context, email string) (models.UserStruct, error) {
	redisKey := "user:" + email

	// 1. Try getting from Redis first
	cachedUser, err := r.rdb.Get(ctx, redisKey).Result()
	if err == nil {
		var user models.UserStruct
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			log.Printf("[CACHE HIT] User %s found in Redis", email)
			return user, nil
		}
		log.Println("[CACHE] Failed to unmarshal cached user:", err)
	} else if err != redis.Nil {
		log.Println("[REDIS ERROR]", err)
	}

	// 2. Query from PostgreSQL if not in cache
	query := `
		SELECT id, email, password, COALESCE(role, 'user') as role 
		FROM users 
		WHERE email = $1
	`
	var user models.UserStruct
	err = r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		return models.UserStruct{}, err
	}

	// 3. Cache the user data in Redis
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println("[MARSHAL ERROR]", err)
	} else {
		if err := r.rdb.Set(ctx, redisKey, userJSON, 2*time.Minute).Err(); err != nil {
			log.Println("[REDIS SET ERROR]", err)
		}
	}

	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, email, password string) (pgconn.CommandTag, error) {
	query := `INSERT INTO users (email, password, role) VALUES ($1, $2, 'user') RETURNING id`
	values := []any{email, password}
	var id int
	if err := r.db.QueryRow(ctx, query, values...).Scan(&id); err != nil {
		return pgconn.CommandTag{}, err
	}
	queryProfile := `INSERT INTO profiles (user_id) VALUES ($1)`
	return r.db.Exec(ctx, queryProfile, id)

}
