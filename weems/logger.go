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
	"time"
)

const (
	TRACE = iota
	INFO
	WARN
	ERROR
	FATAL
)

const (
	callDepth  = 2
	timeFormat = "15:04:05"
)

var msgLevelString = [...]string{
	"\033[96mTRACE\033[0m",
	"\033[94mINFO \033[0m",
	"\033[33mWARN \033[0m",
	"\033[31mERROR\033[0m",
	"\033[37m\033[101mFATAL\033[0m",
}

var loggingLevels = [...]int{TRACE, INFO, WARN, ERROR, FATAL}

type Logger struct {
	name   string
	level  *int
	output *io.Writer
}

func NewLogger(name string) Logger {
	return Logger{
		name:   name,
		level:  &globalLevel,
		output: &globalWriter,
	}
}

func (logger *Logger) SetWriter(writer *io.Writer) {
	logger.output = writer
}

func (logger Logger) GetLevel() int {
	return *logger.level
}

func (logger *Logger) SetLevel(level int) {
	logger.level = &loggingLevels[level]
}

func (logger *Logger) Trace(msg string, args ...interface{}) {
	logger.log(TRACE, msg, args...)
}

func (logger *Logger) Info(msg string, args ...interface{}) {
	logger.log(INFO, msg, args...)
}

func (logger *Logger) Warn(msg string, args ...interface{}) {
	logger.log(WARN, msg, args...)
}

func (logger *Logger) Error(msg string, args ...interface{}) {
	logger.log(ERROR, msg, args...)
}

func (logger *Logger) Fatal(msg string, args ...interface{}) {
	logger.log(FATAL, msg, args...)
	os.Exit(1)
}

func (logger *Logger) log(level int, msg string, args ...interface{}) {
	// Avoid any work if unnecessary
	if level < *logger.level {
		return
	}

	// Make sure we get timestamp as early as possible
	now := time.Now().Format(timeFormat)

	// Get filename and line number by fetching runtime information
	_, filename, line, ok := runtime.Caller(callDepth)
	if !ok {
		filename = "<unknown>.go"
		line = 0
	} else {
		filename = filepath.Base(filename)
	}

	// Generate logger format string
	format := fmt.Sprintf(
		"%s %s %s \033[94m%s:%d\033[0m %%s\n", now, msgLevelString[level], logger.name, filename, line,
	)

	// Format string if necessary
	fmtMsg := msg
	if len(args) > 0 {
		fmtMsg = fmt.Sprintf(msg, args...)
	}

	// Perform the writing operation, thread-safely
	globalMutex.Lock()
	_, _ = fmt.Fprintf(*logger.output, format, fmtMsg)
	globalMutex.Unlock()
}
