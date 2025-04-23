package middleware

import (
	"github.com/ladmakhi81/learnup/pkg/contracts"
)

type Middleware struct {
	tokenSvc contracts.Token
	redisSvc contracts.Cache
}

func NewMiddleware(tokenSvc contracts.Token, redisSvc contracts.Cache) *Middleware {
	return &Middleware{
		tokenSvc: tokenSvc,
		redisSvc: redisSvc,
	}
}
