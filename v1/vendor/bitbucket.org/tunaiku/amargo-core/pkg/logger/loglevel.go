package logger

import (
	"github.com/sirupsen/logrus"
	"strings"
)

type LogLevel string

func (s LogLevel) ToLogrusLevelOrFallback(fallbackLevel logrus.Level) logrus.Level {
	return getLogrusLevel(string(s), fallbackLevel)
}

func getLogrusLevel(logLevel string, fallbackLevel logrus.Level) logrus.Level {
	switch strings.ToUpper(logLevel) {
	case "FATAL":
		return logrus.FatalLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "TRACE":
		return logrus.TraceLevel
	case "PANIC":
		return logrus.PanicLevel
	case "WARN":
		return logrus.WarnLevel
	default:
		I().Infof("Unsupported log level: %s,so returning using fallback level instead", logLevel)
		return fallbackLevel
	}
}
