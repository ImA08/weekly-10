package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"minitask1.go/internal/handlers"
	"minitask1.go/internal/middleware"
	"minitask1.go/internal/repositories"
)

func AddScheduleRoutes(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client, mdw *middleware.Middleware) {
	scheduleRepo := repositories.NewScheduleRepository(db, rdb)
	scheduleHandler := handlers.NewScheduleHandler(scheduleRepo)

	scheduleGroup := router.Group("/schedules")

	scheduleGroup.GET("/:id", scheduleHandler.GetScheduleHandler)
	scheduleGroup.POST("", scheduleHandler.CreateScheduleHandler)
}
