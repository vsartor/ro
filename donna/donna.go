// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

// Donna Paulsen handles command line arguments for Ro.

// It handles two different types of arguments: options and arguments.

// Options are optional command line arguments with default values and
// that starts with "-", either containing a parameter if integer or
// string, or simply being a flag.

// Arguments are command line arguments that do not have a preceding dash.

// There are two kinds of options, global options, handled at startup and
// consumed by ro's entry point, and command specific options which are
// handled by the executing code.
