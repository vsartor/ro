// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

var (
	expectedGlobalFlags   []string
	expectedGlobalOptions []string
	expectedFlags         []string
	expectedOptions       []string
)

func init() {
	expectedGlobalFlags = make([]string, 0)
	expectedGlobalOptions = make([]string, 0)
	expectedFlags = make([]string, 0)
	expectedOptions = make([]string, 0)
}

// Validates that the parameter was expected, also returning
// a boolean indicating if the name refers to a flag.
func validate(name string) (bool, bool) {
	for _, flag := range expectedFlags {
		if flag == name {
			return true, true
		}
	}

	for _, option := range expectedOptions {
		if option == name {
			return true, false
		}
	}

	return false, false
}

// Validates that the parameter was expected, also returning
// a boolean indicating if the name refers to a flag. Performs
// the operation for global parameters.
func validateGlobal(name string) (bool, bool) {
	for _, flag := range expectedGlobalFlags {
		if flag == name {
			return true, true
		}
	}

	for _, option := range expectedGlobalOptions {
		if option == name {
			return true, false
		}
	}

	return false, false
}

// Registers the name of a global flag.
func ExpectGlobalFlag(name string) {
	expectedGlobalFlags = append(expectedGlobalFlags, name)
}

// Registers the name of a command flag.
func ExpectFlag(name string) {
	expectedFlags = append(expectedFlags, name)
}

// Registers the name of a global option.
func ExpectGlobalOption(name string) {
	expectedGlobalOptions = append(expectedGlobalOptions, name)
}

// Registers the name of a command option.
func ExpectOption(name string) {
	expectedOptions = append(expectedOptions, name)
}
