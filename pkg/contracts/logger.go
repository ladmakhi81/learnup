package contracts

import (
	"github.com/ladmakhi81/learnup/pkg/dtos"
)

type Log interface {
	Print(message dtos.LogMessage)
	Error(message dtos.LogMessage)
	Warning(message dtos.LogMessage)
}
