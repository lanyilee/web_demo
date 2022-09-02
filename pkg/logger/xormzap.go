package logger

import (
	"go.uber.org/zap"
	"xorm.io/core"
)

type XormZapLogger struct {
	logger  *zap.Logger
	level   core.LogLevel
	showSQL bool
}

// NewXormZapLogger xorm使用zap日志
func NewXormZapLogger(logger *zap.Logger) core.ILogger {
	return &XormZapLogger{
		logger:  logger,
		level:   core.LOG_ERR,
		showSQL: true,
	}
}

// Error implement core.ILogger
func (l *XormZapLogger) Error(v ...interface{}) {
	if l.level <= core.LOG_ERR {
		l.logger.Sugar().Error(v...)
	}
	return
}

// Errorf implement core.ILogger
func (l *XormZapLogger) Errorf(format string, v ...interface{}) {
	if l.level <= core.LOG_ERR {
		l.logger.Sugar().Errorf(format, v...)
	}
	return
}

// Debug implement core.ILogger
func (l *XormZapLogger) Debug(v ...interface{}) {
	if l.level <= core.LOG_DEBUG {
		l.logger.Sugar().Debug(v...)
	}
	return
}

// Debugf implement core.ILogger
func (l *XormZapLogger) Debugf(format string, v ...interface{}) {
	if l.level <= core.LOG_DEBUG {
		l.logger.Sugar().Debugf(format, v...)
	}
	return
}

// Info implement core.ILogger
func (l *XormZapLogger) Info(v ...interface{}) {
	if l.level <= core.LOG_INFO {
		l.logger.Sugar().Info(v...)
	}
	return
}

// Infof implement core.ILogger
func (l *XormZapLogger) Infof(format string, v ...interface{}) {
	if l.level <= core.LOG_INFO {
		l.logger.Sugar().Infof(format, v...)
	}
	return
}

// Warn implement core.ILogger
func (l *XormZapLogger) Warn(v ...interface{}) {
	if l.level <= core.LOG_WARNING {
		l.logger.Sugar().Warn(v...)
	}
	return
}

// Warnf implement core.ILogger
func (l *XormZapLogger) Warnf(format string, v ...interface{}) {
	if l.level <= core.LOG_WARNING {
		l.logger.Sugar().Warnf(format, v...)
	}
	return
}

// Level implement core.ILogger
func (l *XormZapLogger) Level() core.LogLevel {
	return l.level
}

// SetLevel implement core.ILogger
func (l *XormZapLogger) SetLevel(lv core.LogLevel) {
	l.level = lv
	return
}

// ShowSQL implement core.ILogger
func (l *XormZapLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (l *XormZapLogger) IsShowSQL() bool {
	return l.showSQL
}
