package types

import (
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type TokenClaim struct {
	UserID uint
	jwt.RegisteredClaims
}

func NewTokenClaim(userID uint, exp time.Time) *TokenClaim {
	return &TokenClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   strconv.Itoa(int(userID)),
		},
	}
}
