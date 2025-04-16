package contracts

import (
	"github.com/ladmakhi81/learnup/pkg/dtos"
)

type EnvProvider interface {
	LoadLearnUp() (*dtos.EnvConfig, error)
}
