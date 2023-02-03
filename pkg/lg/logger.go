package lg

// Logger is a (very simple) leveled logger
type Logger interface {
	// Error logs error-level messages
	Error(fmt string, args ...interface{})
	// Warning logs warning-level messages
	Warning(fmt string, args ...interface{})
	// Info logs info-level messages
	Info(fmt string, args ...interface{})
	// Debug logs debug-level messages
	Debug(fmt string, args ...interface{})
}
