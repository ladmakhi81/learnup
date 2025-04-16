package contracts

import "github.com/ladmakhi81/learnup/types"

type Validation interface {
	Validate(dto any) *types.ClientError
}
