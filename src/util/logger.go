package util

import (
	"os"

	"github.com/Sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// NewLogger creates a new logger
func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.WarnLevel
	log.Formatter = new(prefixed.TextFormatter)
	log.Level = logrus.DebugLevel

	return log
}
