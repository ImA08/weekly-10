package middleware

import (
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"minitask1.go/pkg"
)

func (m *Middleware) AccessGate(allowedRole ...string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		// 1. ambil payload/claims dari context gin
		claims, exist := ctx.Get("Payload")
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Please Login",
			})
			return
		}
		// type assertion claims menjadi pkg.claims
		userClaims, ok := claims.(*pkg.Claims)
		// log.Println(userClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		// cek role yang ada di claims

		if !slices.Contains(allowedRole, userClaims.Role) {
			log.Println("[DEBUG] allowed role", allowedRole, userClaims.Role)

			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Forbidden",
			})
			return
		}
		ctx.Next()
	}
}
