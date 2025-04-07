package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/auth/constant"
	"net/http"
	"strings"
)

func (m Middleware) CheckAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("authorization")
		if authorization == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}
		authorizationSegments := strings.Split(authorization, " ")
		if len(authorizationSegments) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}
		bearer := strings.Trim(strings.ToLower(authorizationSegments[0]), " ")
		token := authorizationSegments[1]
		if bearer != "bearer" || token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}
		cachedToken, cachedTokenErr := m.redisSvc.GetVal(constant.LoginCacheKey)
		if cachedTokenErr != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if cachedToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}
		tokenClaims, tokenErr := m.tokenSvc.VerifyToken(token)
		if tokenErr != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}
		ctx.Set("AUTH", tokenClaims)
		ctx.Next()
	}
}
