package router

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

var r = gin.Default()

func main() {

	dbEnv := []any{}
	dbEnv = append(dbEnv, os.Getenv("DBUSER"))
	dbEnv = append(dbEnv, os.Getenv("DBPASS"))
	dbEnv = append(dbEnv, os.Getenv("DBHOST"))
	dbEnv = append(dbEnv, os.Getenv("DBPORT"))
	dbEnv = append(dbEnv, os.Getenv("DBNAME"))

	dbString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbEnv...)
	dbClient, err := pgxpool.New(context.Background(), dbString)
	if err != nil {
		log.Printf("Unable to create connection pool : %v/n", err)
		os.Exit(1)
	}

	defer func() {
		log.Println("Closing DB...")
		dbClient.Close()
	}()

	v1 := r.Group(("/api/v1"))

	// AUTHENTICATION ROUTER
	auth := v1.Group("/auth")
	{
		auth.POST("/register", registrationHandler)
		auth.POST("/login", loginHandler)
	}

	// MOVIES ROUTER
	movies := v1.Group("/movies")
	{
		// movies.GET("", filterMoviesHandler)
		// movies.GET("/upcoming", getUpcomingMovies)
		movies.GET("/popular", getPopularMovies)
		// movies.GET("/:id", getMovieDetailHandler)
	}

	// Schedule routes
	// schedules := v1.Group("/schedules")
	{
		// schedules.GET("", getSchedulesHandler)
		// schedules.GET("/:id/seats", getAvailableSeatsHandler)
	}

	// User routes (require authentication middleware)
	// user := v1.Group("/user")
	// user.Use(authMiddleware())
	{
		// Order routes
		// user.POST("/orders", createOrderHandler)
		// user.GET("/orders", getOrderHistoryHandler)

		// Profile routes
		// user.GET("/profile", getProfileHandler)
		// user.PUT("/profile", editProfileHandler)
	}

	// Admin routes (require admin middleware)
	// admin := v1.Group("/admin")
	// admin.Use(adminMiddleware())
	{
		// 	admin.GET("/movies", getAllMoviesHandler)
		// 	admin.DELETE("/movies/:id", deleteMovieHandler)
		// 	admin.PUT("/movies/:id", updateMovieHandler)
		// 	admin.POST("/movies", createMovieHandler)
	}
}
