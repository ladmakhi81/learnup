package koanf

import (
	koanfEnv "github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"strings"
)

type KoanfEnvSvc struct{}

func NewKoanfEnvSvc() *KoanfEnvSvc {
	return &KoanfEnvSvc{}
}

func (svc KoanfEnvSvc) LoadLearnUp() (*dtos.EnvConfig, error) {
	k := koanf.New(".")
	provider := koanfEnv.Provider("LEARNUP_", "__", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "LEARNUP_"))
	})
	if err := k.Load(provider, nil); err != nil {
		return nil, err
	}
	var envData dtos.EnvConfig
	if err := k.Unmarshal("", &envData); err != nil {
		return nil, err
	}
	return &envData, nil
}
