package token

import "github.com/ladmakhi81/learnup/types"

type Token interface {
	GenerateToken(userID uint) (string, error)
	VerifyToken(tokenString string) (*types.TokenClaim, error)
}
