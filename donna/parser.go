// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

import (
	"os"
	"strings"
)

var localParametersIdx int

func parseGlobalArg(arg string, idx *int) {
	argName := strings.TrimPrefix(arg, "--")

	isExpected, isFlag := validate(argName)

	if !isExpected {
		logger.Fatal("Unexpected argument '%s'.", arg)
	}

	if isFlag {
		globalFlags = append(globalFlags, argName)
		*idx++
		return
	}

	// Assert that option received an associated value
	if *idx == len(os.Args)-1 {
		logger.Fatal("Global option '%s' has no associated value.", argName)
	}

	optionValue := os.Args[*idx+1]
	globalOptions[argName] = optionValue
	*idx += 2
}

func parseLocalArg(arg string, idx *int) {
	argName := strings.TrimPrefix(arg, "--")

	isExpected, isFlag := validate(argName)

	if !isExpected {
		logger.Fatal("Unexpected argument '%s'.", arg)
	}

	if isFlag {
		flags = append(flags, argName)
		*idx++
		return
	}

	// Assert that option received an associated value
	if *idx == len(os.Args)-1 {
		logger.Fatal("Option '%s' has no associated value.", argName)
	}

	optionValue := os.Args[*idx+1]
	options[argName] = optionValue
	*idx += 2
}

// Parse global command line arguments.
func ParseGlobal() {
	idx := 1
	for idx < len(os.Args) {
		arg := os.Args[idx]

		if strings.HasPrefix(arg, "--") {
			parseGlobalArg(arg, &idx)
		} else {
			// Found a regular Argument. This means the end of
			// global parameters.
			break
		}
	}

	// Parse regular Arguments
	for idx < len(os.Args) {
		arg := os.Args[idx]

		if strings.HasPrefix(arg, "-") {
			// Found a flag/option; stop global parsing.
			localParametersIdx = idx
			break
		} else {
			args.arguments = append(args.arguments, arg)
			idx++
		}
	}
}

// Parses local command line arguments.
func Parse() {
	idx := localParametersIdx
	for idx < len(os.Args) {
		arg := os.Args[idx]

		if strings.HasPrefix(arg, "--") {
			parseLocalArg(arg, &idx)
		} else {
			// Should not have arguments at this point.
			logger.Fatal("Unexpected argument '%s' after flags/options.", arg)
		}
	}
}
