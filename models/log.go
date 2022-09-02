package models

// Logger log libaray
type Logger interface {
	Debugf(format string, args ...interface{})
	Debug(v ...interface{})
	Infof(format string, args ...interface{})
	Info(v ...interface{})
	Warnf(format string, args ...interface{})
	Warn(v ...interface{})
	Errorf(format string, args ...interface{})
	Error(v ...interface{})
}
