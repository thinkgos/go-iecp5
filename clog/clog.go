package clog

import (
	"log"
	"os"
	"sync/atomic"
)

// LogProvider RFC5424 log message levels only Debug Warn and Error
type LogProvider interface {
	Critical(format string, v ...interface{})
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

// Clog 日志内部调试实现
type Clog struct {
	logger LogProvider
	// is log output enabled,1: enable, 0: disable
	hasLog uint32
}

// New 创建一个新的日志无前缀
func New() *Clog {
	return NewWithPrefix("")
}

// NewWithPrefix 创建一个新的日志，采用指定prefix前缀
func NewWithPrefix(prefix string) *Clog {
	return &Clog{
		logger: newDefaultLogger(prefix),
	}
}

// LogMode set enable or disable log output when you has set logger
func (this *Clog) LogMode(enable bool) {
	if enable {
		atomic.StoreUint32(&this.hasLog, 1)
	} else {
		atomic.StoreUint32(&this.hasLog, 0)
	}
}

// SetLogProvider set logger provider
func (this *Clog) SetLogProvider(p LogProvider) {
	if p != nil {
		this.logger = p
	}
}

// Critical Log CRITICAL level message.
func (this *Clog) Critical(format string, v ...interface{}) {
	if atomic.LoadUint32(&this.hasLog) == 1 {
		this.logger.Critical(format, v...)
	}
}

// Error Log ERROR level message.
func (this *Clog) Error(format string, v ...interface{}) {
	if atomic.LoadUint32(&this.hasLog) == 1 {
		this.logger.Error(format, v...)
	}
}

// Warn Log WARN level message.
func (this *Clog) Warn(format string, v ...interface{}) {
	if atomic.LoadUint32(&this.hasLog) == 1 {
		this.logger.Warn(format, v...)
	}
}

// Debug Log DEBUG level message.
func (this *Clog) Debug(format string, v ...interface{}) {
	if atomic.LoadUint32(&this.hasLog) == 1 {
		this.logger.Debug(format, v...)
	}
}

// default log
type logger struct {
	*log.Logger
}

var _ LogProvider = (*logger)(nil)

// newDefaultLogger new default logger with prefix output os.Stderr
func newDefaultLogger(prefix string) *logger {
	return &logger{
		log.New(os.Stderr, prefix, log.LstdFlags),
	}
}

// Critical Log CRITICAL level message.
func (this *logger) Critical(format string, v ...interface{}) {
	this.Printf("[C]: "+format, v...)
}

// Error Log ERROR level message.
func (this *logger) Error(format string, v ...interface{}) {
	this.Printf("[E]: "+format, v...)
}

// Warn Log WARN level message.
func (this *logger) Warn(format string, v ...interface{}) {
	this.Printf("[W]: "+format, v...)
}

// Debug Log DEBUG level message.
func (this *logger) Debug(format string, v ...interface{}) {
	this.Printf("[D]: "+format, v...)
}
