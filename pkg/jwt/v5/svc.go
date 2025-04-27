package jwtv5

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ladmakhi81/learnup/internals/auth/constant"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"github.com/ladmakhi81/learnup/shared/types"
	"strings"
	"time"
)

type JwtSvc struct {
	config   *dtos.EnvConfig
	redisSvc contracts.Cache
}

func NewJwtSvc(config *dtos.EnvConfig, redisSvc contracts.Cache) *JwtSvc {
	return &JwtSvc{
		config:   config,
		redisSvc: redisSvc,
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
		return nil, nil
	}
	return claims, nil
}

func (svc JwtSvc) getSecretKey() []byte {
	return []byte(svc.config.App.TokenSecretKey)
}

func (svc JwtSvc) DecodeToken(tokenString string) (*types.TokenClaim, error) {
	if tokenString == "" {
		return nil, nil
	}
	tokenStringSegments := strings.Split(tokenString, " ")
	if len(tokenStringSegments) != 2 {
		return nil, nil
	}
	tokenBearer := strings.Trim(strings.ToLower(tokenStringSegments[0]), " ")
	token := tokenStringSegments[1]
	if tokenBearer != "bearer" || token == "" {
		return nil, nil
	}
	cachedToken, err := svc.redisSvc.GetVal(constant.LoginCacheKey)
	if err != nil {
		return nil, err
	}
	if cachedToken == "" {
		return nil, nil
	}
	tokenClaims, err := svc.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	return tokenClaims, nil
}
