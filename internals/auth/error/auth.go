package error

import "github.com/ladmakhi81/learnup/types"

var (
	Auth_InvalidCredentials = types.NewNotFoundError("auth.errors.invalid_credentials")
)
