package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"minitask1.go/internal/handlers"
	"minitask1.go/internal/middleware"
	"minitask1.go/internal/repositories"
)

func AddMovieRoutes(router *gin.Engine, db *pgxpool.Pool, mdw *middleware.Middleware, rdb *redis.Client) {
	movieRepo := repositories.NewMovieRepository(db, rdb)
	movieHandler := handlers.NewMovieHandler(movieRepo)

	movieGroup := router.Group("/movies")
	{
		movieGroup.GET("", movieHandler.GetMovies)
		movieGroup.GET("/upcoming", mdw.VerifyToken, mdw.AccessGate("user"), movieHandler.GetUpcomingMovies)
	}
}
