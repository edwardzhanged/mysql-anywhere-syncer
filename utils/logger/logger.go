package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func NewLogger() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetOutput(os.Stdout)
	Logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
		FullTimestamp:   true,
	}
}
