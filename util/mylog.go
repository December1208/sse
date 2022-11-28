package util

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SentryLog struct {
	L   *zap.Logger
	Ctx *gin.Context
}

func (log *SentryLog) CapError(msg string, fields ...zap.Field) {
	sentrygin.GetHubFromContext(log.Ctx).CaptureMessage(msg)
	log.L.Error(msg, fields...)
}

func GetLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

var MyLogger = GetLogger()
