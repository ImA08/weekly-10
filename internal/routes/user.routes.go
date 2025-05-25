package routes

import (
	"github.com/gin-gonic/gin"
	"minitask1.go/internal/handlers"
	"minitask1.go/internal/repositories"
)

func addUserRouter(router *gin.Engine, userRepo *repositories.UserRepository) {

	userHandler := handlers.NewUserHandler(userRepo)
	authRouter := router.Group("/auth")

	{
		authRouter.POST("", userHandler.Login)
		authRouter.POST("/register", userHandler.Register)
	}

}
