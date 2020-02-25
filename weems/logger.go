package weems

import (
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"os"
)

var logger loggo.Logger

// Initialize the logger at startup.
// Any problems with Weems startup should be reason for the
// program panicking.
func init() {
	_, err := loggo.ReplaceDefaultWriter(loggocolor.NewColorWriter(os.Stdout))
	if err != nil {
		panic("Failure initializing Weem's logger.")
	}

	logger = loggo.GetLogger("Weems")

	// Default logger level should be warning
	logger.SetLogLevel(loggo.WARNING)
}

func Critical(msg string, args ...interface{}) {
	logger.Criticalf(msg, args...)
	os.Exit(1)
}

func Error(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Warningf(msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

func SetVerbose() {
	logger.SetLogLevel(loggo.INFO)
}

func SetQuiet() {
	logger.SetLogLevel(loggo.ERROR)
}
