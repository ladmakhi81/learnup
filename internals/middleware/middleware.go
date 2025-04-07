package middleware

import (
	"github.com/ladmakhi81/learnup/pkg/cache"
	"github.com/ladmakhi81/learnup/pkg/token"
)

type Middleware struct {
	tokenSvc token.Token
	redisSvc cache.Cache
}

func NewMiddleware(tokenSvc token.Token, redisSvc cache.Cache) Middleware {
	return Middleware{
		tokenSvc: tokenSvc,
		redisSvc: redisSvc,
	}
}
