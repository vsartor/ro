// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

var (
	expectedGlobalFlags   []string
	expectedGlobalOptions []string
	expectedLocalFlags    []string
	expectedLocalOptions  []string
)

func init() {
	expectedGlobalFlags = make([]string, 0)
	expectedGlobalOptions = make([]string, 0)
	expectedLocalFlags = make([]string, 0)
	expectedLocalOptions = make([]string, 0)
}

// Validates that the parameter was expected, also returning
// a boolean indicating if the name refers to a flag.
func validate(name string, global bool) (bool, bool) {
	var expectedFlags []string
	var expectedOptions []string

	if global {
		expectedFlags = expectedGlobalFlags
		expectedOptions = expectedGlobalOptions
	} else {
		expectedFlags = expectedLocalFlags
		expectedOptions = expectedLocalOptions
	}

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

// Registers the name of a global flag.
func ExpectGlobalFlag(name string) {
	expectedGlobalFlags = append(expectedGlobalFlags, name)
}

// Registers the name of a command flag.
func ExpectFlag(name string) {
	expectedLocalFlags = append(expectedLocalFlags, name)
}

// Registers the name of a global option.
func ExpectGlobalOption(name string) {
	expectedGlobalOptions = append(expectedGlobalOptions, name)
}

// Registers the name of a command option.
func ExpectOption(name string) {
	expectedLocalOptions = append(expectedLocalOptions, name)
}
