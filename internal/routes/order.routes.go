package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"minitask1.go/internal/handlers"
	"minitask1.go/internal/middleware"
	"minitask1.go/internal/repositories"
)

func AddOrderRoutes(router *gin.Engine, db *pgxpool.Pool, mdw *middleware.Middleware) {
	orderRepo := repositories.NewOrderRepository(db)
	orderHandler := handlers.NewOrderHandler(orderRepo)

	orderGroup := router.Group("/order")
	{
		orderGroup.POST("", mdw.VerifyToken, mdw.AccessGate("user"), orderHandler.CreateOrder)
		orderGroup.GET("/:orderId", mdw.VerifyToken, mdw.AccessGate("user"), orderHandler.FindOrderById)
	}
}
