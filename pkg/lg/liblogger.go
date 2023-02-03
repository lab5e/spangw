package lg

import "log"

type libLogger struct {
}

func (l *libLogger) Error(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
func (l *libLogger) Warning(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
func (l *libLogger) Info(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
func (l *libLogger) Debug(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
