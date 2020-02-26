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

// Returns whether flag name is expected.
func isExpected(name string, expectedName []string) bool {
	for _, flag := range expectedName {
		if flag == name {
			return true
		}
	}
	return false
}

// Performs generalized validation logic.
func validate(
	flags []string,
	expectedFlags []string,
	options map[string]string,
	expectedOptions []string,
) {
	for _, flag := range flags {
		if !isExpected(flag, expectedFlags) {
			logger.Fatal("Unexpected flag '%s'.", flag)
		}
	}

	for option, _ := range options {
		if !isExpected(option, expectedOptions) {
			logger.Fatal("Unexpected option '%s'.", option)
		}
	}
}

// Validates that only expected global flags/options were found.
// Returns the name of the unexpected flag/option and a boolean value
// which is true if there were no errors.
func ValidateGlobal() {
	validate(globalFlags, expectedGlobalFlags, globalOptions, expectedGlobalOptions)
}

// Validates that only command-specific flag/options were found.
// Returns the name of the unexpected flag/option and a boolean value
// which is true if there were no errors.
func ValidateLocal() {
	validate(flags, expectedFlags, options, expectedOptions)
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
