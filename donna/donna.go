// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

import "github.com/vsartor/ro/weems"

// Donna handles command line arguments for Ro.
//
// It handles three different types of command line parameters: flags,
// options and arguments.
//
// Flags start with a "--" and do not take arguments, their presence/absence
// indicates their value. Options start with a single dash "-" and take in a
// following value, associated to that option. Arguments are regular names.
//
// Options and flags can either be global, relating to Ro's own operations,
// preceding all arguments, or command specific, following all arguments.
// To illustrate, this is Ro's usage map:
//     ro [global flags/options] <argument chain> [command flags/options]
//
// In case this invocation rule is violated Donna will log a critical message.
//
// Flags and options are validated, with global ones being validated at
// entry point, and with command ones being validated at the command routine's
// entry point. Their values are consumed at request.
//
// Arguments are consumed one by one and validation happens online as each
// value is consumed before dispatching to another routine.

var logger weems.Logger

func init() {
	logger = weems.NewLogger("donna")
}
