package utils

import (
	"time"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"
)

func init() {
	// apmlogrus.Hook will send "error", "panic", and "fatal" level log messages to Elastic APM.
	logrus.AddHook(&apmlogrus.Hook{})
}

func NewLogger(env, service string) *logrus.Entry {
	logger := logrus.New()
	customFormatter := &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	logger.SetFormatter(customFormatter)

	return logger.WithFields(logrus.Fields{"services.name": service, "services.environment": env})
}
