package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	swaggerFile "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "minitask1.go/docs"
	"minitask1.go/internal/middleware"
	"minitask1.go/internal/repositories"
)

func InitRouter(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	// gin engine initialization
	router := gin.Default()
	middleware := middleware.InitMiddleware()

	router.Use(middleware.CORSMiddleware)

	userRepo := repositories.NewUserRepository(db, rdb)
	addUserRouter(router, userRepo)
	AddMovieRoutes(router, db, middleware, rdb)
	AddProfileRoutes(router, db, middleware, rdb)
	AddOrderRoutes(router, db, middleware)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))

	return router

}
