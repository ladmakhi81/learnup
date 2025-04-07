package avalaiv1

import (
	"context"
	"github.com/ladmakhi81/learnup/pkg/ai"
	"github.com/ladmakhi81/learnup/pkg/env"
	"github.com/sashabaranov/go-openai"
)

type AvalAiSvc struct {
	config *env.EnvConfig
}

func NewAvalAiSvc(config *env.EnvConfig) *AvalAiSvc {
	return &AvalAiSvc{
		config: config,
	}
}

func (svc AvalAiSvc) GenerateImage(
	ctx context.Context,
	dto ai.AiGenerateImageDto,
) (string, error) {
	apiKey := svc.config.App.OPENAI_KEY
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.avalai.ir/v1"
	client := openai.NewClientWithConfig(config)
	resp, respErr := client.CreateImage(ctx,
		openai.ImageRequest{
			Model:   "dall-e-3",
			N:       1,
			Prompt:  dto.Prompt,
			Size:    dto.Size,
			Quality: "standard",
		},
	)
	if respErr != nil {
		return "", respErr
	}
	return resp.Data[0].URL, nil
}
