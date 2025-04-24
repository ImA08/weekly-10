package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type userStruct struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

var users = []userStruct{
	{Email: "caelus@trailblazer.hsr", Password: "lordtrashcane"},
	{Email: "march7th@trailblazer.hsr", Password: "caelussimper"},
}

// REGISTRATION HANDLER

func registrationHandler(ctx *gin.Context) {
	newUser := &userStruct{}
	if err := ctx.ShouldBind(newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input",
		})
		return
	}

	users = append(users, *newUser)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Succes",
		"data": users,
	})
}

// LOGIN HANDLER

func loginHandler(ctx *gin.Context) {
	// Bind request body ke struct
	var loginData userStruct
	if err := ctx.ShouldBind(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input format",
		})
		return
	}

	// Validasi email dan password tidak kosong
	if loginData.Email == "" || loginData.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Email and password are required",
		})
		return
	}

	// Cari user dengan email yang cocok
	var foundUser *userStruct
	for _, user := range users {
		if user.Email == loginData.Email {
			foundUser = &user
			break
		}
	}

	// Jika user tidak ditemukan
	if foundUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	// Verifikasi password
	if foundUser.Password != loginData.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	// Jika login berhasil
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"email": foundUser.Email,
		},
	})
}
