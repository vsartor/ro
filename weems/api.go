// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package weems

var logger Logger

func init() {
	logger = NewLogger(WARNING)
}

func Fatal(msg string, args ...interface{}) {
	logger.Fatal(msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger.Error(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Warn(msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger.Info(msg, args...)
}

func SetVerbose() {
	logger.SetLevel(INFO)
}

func SetQuiet() {
	logger.SetLevel(ERROR)
}
