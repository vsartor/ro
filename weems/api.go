package weems

var logger Logger

// Initialize the logger at startup.
// Any problems with Weems startup should be reason for the
// program panicking.
func init() {
	logger = NewLogger(WARNING)
}

func Fatal(msg string, args ...interface{}) {
	logger.Fatal(msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger.Error(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Warn(msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger.Info(msg, args...)
}

func SetVerbose() {
	logger.SetLevel(INFO)
}

func SetQuiet() {
	logger.SetLevel(ERROR)
}
