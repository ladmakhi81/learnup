package handler

import (
	forumService "github.com/ladmakhi81/learnup/internals/forum/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
)

type Handler struct {
	forumSvc       forumService.ForumService
	translationSvc contracts.Translator
}

func NewHandler(
	forumSvc forumService.ForumService,
	translationSvc contracts.Translator,
) *Handler {
	return &Handler{
		forumSvc:       forumSvc,
		translationSvc: translationSvc,
	}
}
