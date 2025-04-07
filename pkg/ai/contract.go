package ai

import "context"

type Ai interface {
	GenerateImage(ctx context.Context, dto AiGenerateImageDto) (string, error)
}
