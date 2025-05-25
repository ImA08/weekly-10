package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"minitask1.go/internal/handlers"
	"minitask1.go/internal/middleware"
	"minitask1.go/internal/repositories"
)

func AddProfileRoutes(router *gin.Engine, db *pgxpool.Pool, mdw *middleware.Middleware, rdb *redis.Client) {
	profileRepo := repositories.NewProfileRepository(db)
	profileHandler := handlers.NewProfileHandler(profileRepo)

	profileGroup := router.Group("/profile")
	{
		profileGroup.PATCH("/edit", mdw.VerifyToken, mdw.AccessGate("user"), profileHandler.EditProfile)
		profileGroup.GET("", mdw.VerifyToken, mdw.AccessGate("user"), profileHandler.GetProfile)
	}
}
