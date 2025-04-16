package contracts

import (
	"context"
	"github.com/ladmakhi81/learnup/pkg/dtos"
)

type Ai interface {
	GenerateImage(ctx context.Context, dto dtos.AiGenerateImageDto) (string, error)
}
