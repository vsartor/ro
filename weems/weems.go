// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package weems

// Weems keeps a log of everything that happens.
//
// Weems provides a Logger type from which basic, thread-safe, logging
// operations can be performed.
//
// There are four different logging levels, INFO, WARN, ERROR and FATAL.
// These levels affect three different things:
//   1) Whether output is shown based on the Logger's level setting, i.e.,
//      if a Logger's level is WARN, INFO logs will be ignored.
//   2) The output itself, which will present the log level.
//   3) Control flow, as Fatal logs will result in a os.Exit(1) call after
//      the message is logged.
//
// Loggers also have an associated name and an io.Writer output.
//
// The output format of Weem's Loggers are fixed, use ANSI color escape
// codes for ease of reading, and keep to the following format:
//     <hour>:<minute>:<seconds> <level> <name> <filename>:<line-number> <message>
//
// The filename and line number are fetched at runtime using Go's runtime
// module from the standard library.
//
// Logging operations are thread-safe, as all of Weem's Loggers share a
// global mutex. This is done so that different loggers writing to the same
// io.Writer will not conflict with each other.
