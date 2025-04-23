package contracts

import (
	"github.com/ladmakhi81/learnup/shared/types"
)

type Token interface {
	GenerateToken(userID uint) (string, error)
	VerifyToken(tokenString string) (*types.TokenClaim, error)
}
