// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package weems

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	INFO = iota
	WARNING
	ERROR
	FATAL
)

const (
	callDepth  = 3
	timeFormat = "15:04:05"
)

var msgLevelString = [...]string{
	"\033[34m\033[40mINFO \033[0m",
	"\033[33m\033[40mWARN \033[0m",
	"\033[31m\033[40mERROR\033[0m",
	"\033[37m\033[101mFATAL\033[0m",
}

type Logger struct {
	level      int
	writeMutex sync.Mutex
	output     io.Writer
}

func NewLogger(level int) Logger {
	return Logger{level: level, output: os.Stderr}
}

func (logger *Logger) SetWriter(writer io.Writer) {
	logger.output = writer
}

func (logger Logger) GetLevel() int {
	return logger.level
}

func (logger *Logger) SetLevel(level int) {
	logger.level = level
}

func (logger *Logger) Info(msg string, args ...interface{}) {
	logger.log(INFO, msg, args...)
}

func (logger *Logger) Warn(msg string, args ...interface{}) {
	logger.log(WARNING, msg, args...)
}

func (logger *Logger) Error(msg string, args ...interface{}) {
	logger.log(ERROR, msg, args...)
}

func (logger *Logger) Fatal(msg string, args ...interface{}) {
	logger.log(FATAL, msg, args...)
	os.Exit(1)
}

func (logger *Logger) log(level int, msg string, args ...interface{}) {
	// Avoid any work if unecessary
	if level < logger.level {
		return
	}

	// Make sure we get timestamp as early as possible
	now := time.Now()

	// Get filename and line number by fetching runtime information
	_, filename, line, ok := runtime.Caller(callDepth)
	if !ok {
		filename = "<unknown>.go"
		line = 0
	} else {
		filename = filepath.Base(filename)
	}

	// String

	// Generate logger format string
	format := fmt.Sprintf(
		"%s %s:%d %s %%s\n", now.Format(timeFormat), filename, line, msgLevelString[level],
	)

	// Format string if necessary
	fmtMsg := msg
	if len(args) > 0 {
		fmtMsg = fmt.Sprintf(msg, args...)
	}

	// Perform the writing operation, thread-safely
	logger.writeMutex.Lock()
	_, _ = fmt.Fprintf(logger.output, format, fmtMsg)
	logger.writeMutex.Unlock()
}
