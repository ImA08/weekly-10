package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"minitask1.go/internal/models"
	"minitask1.go/internal/repositories"
	"minitask1.go/pkg"
)

type UserHandler struct {
	userRepo *repositories.UserRepository
}

func NewUserHandler(userRepo *repositories.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

type UserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Register
// @summary			Register User
// @Router			/auth/register [post]
// @accept 		 	json
// @param			body body models.AuthForm true "register information"
// @produce			json
// @failure			500 {object} models.ErrorResponse
// @failure			400 {object} models.ErrorResponse
// @succes			201 {object} gin.H
func (h *UserHandler) Register(ctx *gin.Context) {

	var body models.UserStruct
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "Field validation") {
			if strings.Contains(err.Error(), "min") {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Password have to greater than 8 characters",
				})
				return
			}
			if strings.Contains(err.Error(), "max") {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Password have to less than 10 characters",
				})
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Name and password required!",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server Error",
		})
		return
	}

	hash := pkg.InitHashConfig()
	hash.UseDefaultConfig()
	hashedPass, err := hash.GenPasswordHash(body.Password)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Hash Failed",
		})
		return
	}

	cmd, err := h.userRepo.CreateUser(ctx.Request.Context(), body.Email, hashedPass)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server Error",
		})
		return
	}
	if cmd.RowsAffected() == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server Error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Selamat datang %s silahkan login", body.Email),
	})

	// var req UserRequest

	// if err := ctx.ShouldBindJSON(&req); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	// 	return
	// }

	// hashedPassword, err := bcrypt.GenerateFromPassword(
	// 	[]byte(req.Password),
	// 	bcrypt.DefaultCost,
	// )
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure password"})
	// 	return
	// }

	// user, err := h.userRepo.CreateUser(ctx.Request.Context(), req.Email, string(hashedPassword))
	// if err != nil {
	// 	ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
	// 	return
	// }

	// // Return
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"ID":    user.ID,
	// 	"Email": user.Email,
	// 	"msg":   "SignUp Succes",
	// })
}

// Login
// @summary			Login User
// @Router			/auth [post]
// @param			body body models.AuthForm true "login information"
// @accept 			json
// @produce			json
// @failure			500 {object} models.ErrorResponse
// @succes			201 {object} models.Response
func (h *UserHandler) Login(ctx *gin.Context) {
	// var req UserRequest
	var body models.UserStruct

	// 1. Input validation
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 2. Get user from repository
	user, err := h.userRepo.LogInUserRepo(ctx.Request.Context(), body.Email)

	if err != nil {

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Authentication failed",
			"detail": "Email not found or invalid",
		})
		return
	}

	hash := pkg.InitHashConfig()
	log.Println("[DEBUG]", user.Password, body.Password)

	valid, err := hash.CompareHashAndPassword(user.Password, body.Password)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server Error",
		})
		return
	}

	log.Println("[DEBUG] FIRST TEST")
	if !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid username/password",
		})
		return
	}
	log.Println("[DEBUG] SECOND TEST")
	claims := pkg.NewClaims(user.ID, user.Role)
	log.Println("[DEBUG] USER CHECK", claims)
	token, err := claims.GenerateToken()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	decoded, err := jwt.ParseWithClaims(token, &pkg.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Printf("[DEBUG] Token verification failed: %v", err)
	} else if claims, ok := decoded.Claims.(*pkg.Claims); ok {
		log.Printf("[DEBUG] Generated token contains: ID=%d, Role=%s", claims.Id, claims.Role)
	}
	// 4. Successful response
	ctx.JSON(http.StatusOK, gin.H{
		"ID":    user.ID,
		"Email": user.Email,
		"token": token,
		"msg":   "LogIn Succes",
	})
}
