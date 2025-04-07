package jwtv5

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ladmakhi81/learnup/pkg/env"
	"github.com/ladmakhi81/learnup/types"
	"time"
)

type JwtSvc struct {
	config *env.EnvConfig
}

func NewJwtSvc(config *env.EnvConfig) *JwtSvc {
	return &JwtSvc{
		config: config,
	}
}

func (svc JwtSvc) GenerateToken(userID uint) (string, error) {
	claim := types.NewTokenClaim(
		userID,
		time.Now().Add(time.Minute*60),
	)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, signedErr := token.SignedString(svc.getSecretKey())
	if signedErr != nil {
		return "", errors.New("Error happen in signed token")
	}
	return signedToken, nil
}

func (svc JwtSvc) VerifyToken(tokenString string) (*types.TokenClaim, error) {
	claims := &types.TokenClaim{}
	token, tokenErr := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return svc.getSecretKey(), nil
	})
	if tokenErr != nil {
		return nil, errors.New("Error happen in verify token")
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}

func (svc JwtSvc) getSecretKey() []byte {
	return []byte(svc.config.App.TokenSecretKey)
}
