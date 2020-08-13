// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a version 3 of the GNU General
// Public License, license that can be found in the LICENSE file.

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
	provider LogProvider
	// is log output enabled,1: enable, 0: disable
	has uint32
}

// NewLogger 创建一个新的日志，采用指定prefix前缀
func NewLogger(prefix string) Clog {
	return Clog{
		defaultLogger{
			log.New(os.Stdout, prefix, log.LstdFlags),
		},
		0,
	}
}

// LogMode set enable or disable log output when you has set provider
func (sf *Clog) LogMode(enable bool) {
	if enable {
		atomic.StoreUint32(&sf.has, 1)
	} else {
		atomic.StoreUint32(&sf.has, 0)
	}
}

// SetLogProvider set provider provider
func (sf *Clog) SetLogProvider(p LogProvider) {
	if p != nil {
		sf.provider = p
	}
}

// Critical Log CRITICAL level message.
func (sf Clog) Critical(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Critical(format, v...)
	}
}

// Error Log ERROR level message.
func (sf Clog) Error(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Error(format, v...)
	}
}

// Warn Log WARN level message.
func (sf Clog) Warn(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Warn(format, v...)
	}
}

// Debug Log DEBUG level message.
func (sf Clog) Debug(format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Debug(format, v...)
	}
}

// default log
type defaultLogger struct {
	*log.Logger
}

var _ LogProvider = (*defaultLogger)(nil)

// Critical Log CRITICAL level message.
func (sf defaultLogger) Critical(format string, v ...interface{}) {
	sf.Printf("[C]: "+format, v...)
}

// Error Log ERROR level message.
func (sf defaultLogger) Error(format string, v ...interface{}) {
	sf.Printf("[E]: "+format, v...)
}

// Warn Log WARN level message.
func (sf defaultLogger) Warn(format string, v ...interface{}) {
	sf.Printf("[W]: "+format, v...)
}

// Debug Log DEBUG level message.
func (sf defaultLogger) Debug(format string, v ...interface{}) {
	sf.Printf("[D]: "+format, v...)
}
