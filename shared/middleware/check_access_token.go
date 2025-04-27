package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m Middleware) CheckAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("authorization")
		claim, err := m.tokenSvc.DecodeToken(authorization)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if claim == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}
		ctx.Set("AUTH", claim)
		ctx.Next()
	}
}
