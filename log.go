// Copyright 2020 wubbalubbaaa. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package easyRpc

import "log"

// DefaultLogger instance
var DefaultLogger Logger = &StdLogger{Enable: true}

// Logger defines log interface
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// SetLogger set default logger for easyRpc
func SetLogger(l Logger) {
	DefaultLogger = l
}

// StdLogger defines default logger
type StdLogger struct {
	Enable bool
}

// Debug simply printf
func (l *StdLogger) Debug(format string, v ...interface{}) {
	if l.Enable {
		log.Printf(format, v...)
	}
}

// Info simply printf
func (l *StdLogger) Info(format string, v ...interface{}) {
	if l.Enable {
		log.Printf(format, v...)
	}
}

// Warn simply printf
func (l *StdLogger) Warn(format string, v ...interface{}) {
	if l.Enable {
		log.Printf(format, v...)
	}
}

// Error simply printf
func (l *StdLogger) Error(format string, v ...interface{}) {
	if l.Enable {
		log.Printf(format, v...)
	}
}
