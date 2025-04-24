package router

import (
	"github.com/gin-gonic/gin"
)

var r = gin.Default()

func main() {

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
