package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"minitask1.go/pkg"
)

type Middleware struct{}

func InitMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) VerifyToken(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Please Login!",
		})
		return
	}

	if !strings.Contains(bearerToken, "Bearer") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Silahkan login terlebih dahulu",
		})
		return
	}

	token := strings.Split(bearerToken, " ")[1]

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Please Login!",
		})
		return
	}
	claims := &pkg.Claims{}
	if err := claims.VerifyToken(token); err != nil {
		log.Println(err.Error())
		// if err.Error() == "expired token" || err.Error() == "token has invalid claims: token is expired" {
		if strings.Contains(err.Error(), "expired") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Please Login Again!",
			})
			return
		}
		if strings.Contains(err.Error(), "malformed") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Identitas login anda rusak, Silahkan login kembali",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}
	// ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
	// 	"message": "Succes!",
	// 	"data":    claims,
	// })
	// return
	log.Println("[DEBUG] CLAIMS TEST", claims)
	ctx.Set("Payload", claims)
	ctx.Next()
}
