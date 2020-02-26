// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

type Args struct {
	numConsumed int
	arguments   []string
}

var args Args

func init() {
	args.arguments = make([]string, 0)
}

// Returns the next command line argument in the chain.
func NextArg() (string, bool) {
	// Make sure we've not consumed all arguments
	if args.numConsumed == len(args.arguments) {
		return "", false
	}

	arg := args.arguments[args.numConsumed]
	args.numConsumed++
	return arg, true
}
