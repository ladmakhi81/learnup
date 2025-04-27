package middleware

import (
	"github.com/ladmakhi81/learnup/pkg/contracts"
)

type Middleware struct {
	tokenSvc contracts.Token
}

func NewMiddleware(tokenSvc contracts.Token) *Middleware {
	return &Middleware{tokenSvc: tokenSvc}
}
