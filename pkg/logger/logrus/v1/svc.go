package logrusv1

import (
	"github.com/ladmakhi81/learnup/pkg/logger"
	"github.com/sirupsen/logrus"
)

type LogrusLoggerSvc struct {
	log *logrus.Logger
}

func NewLogrusLoggerSvc() *LogrusLoggerSvc {
	log := logrus.New()
	//TODO: enable this comment in the production
	//log.SetFormatter(&logrus.JSONFormatter{})
	return &LogrusLoggerSvc{
		log: log,
	}
}

func (svc LogrusLoggerSvc) Print(printInfo logger.LogMessage) {
	svc.log.WithFields(
		map[string]any{
			"metadata": printInfo.Metadata,
		},
	).Info(printInfo.Message)
}

func (svc LogrusLoggerSvc) Error(errorInfo logger.LogMessage) {
	svc.log.WithFields(
		map[string]any{
			"metadata": errorInfo.Metadata,
		},
	).Fatalln(errorInfo.Message)
}

func (svc LogrusLoggerSvc) Warning(warningInfo logger.LogMessage) {
	svc.log.WithFields(
		map[string]any{
			"metadata": warningInfo.Metadata,
		},
	).Warnln(warningInfo.Message)
}
