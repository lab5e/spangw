package lg

var logger Logger

func init() {
	logger = &libLogger{}
}

// ReplaceLogger replaces the default Logger with a custom implementation
func ReplaceLogger(l Logger) {
	logger = l
}

// Error logs error-level messages
func Error(fmt string, args ...interface{}) {
	logger.Error(fmt, args...)
}

// Warning logs warning-level messages
func Warning(fmt string, args ...interface{}) {
	logger.Warning(fmt, args...)
}

// Info logs info-level messages
func Info(fmt string, args ...interface{}) {
	logger.Info(fmt, args...)
}

// Debug logs debug-level messages
func Debug(fmt string, args ...interface{}) {
	logger.Debug(fmt, args...)
}
