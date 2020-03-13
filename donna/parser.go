// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

import (
	"errors"
	"fmt"
	"strings"
)

var iterator paramIterator

// Parses a single cli parameter
func parseCliParam(iterator paramIterator, global bool) error {
	// Get current parameter and trim dashes
	param := iterator.Curr()
	paramName := strings.TrimPrefix(param, "--")

	// Check whether the argument is valid, and if so check whether
	// it's a boolean flag.
	isExpected, isFlag := validate(paramName, global)
	if !isExpected {
		return errors.New(fmt.Sprintf(
			"Unexpected flag '%s'.", param,
		))
	}

	if isFlag {
		// Since parsing localFlags is simpler, handle this branch earlier.
		if global {
			globalFlags = append(globalFlags, paramName)
		} else {
			localFlags = append(localFlags, paramName)
		}
		return nil
	}

	// Assert that option received an associated value
	optionValue, ok := iterator.Next()
	if !ok {
		return errors.New(fmt.Sprintf(
			"Option '%s' did not receive an associated value.", paramName,
		))
	}

	// TODO: Validate option type

	// Save option name and value
	if global {
		globalOptions[paramName] = optionValue
	} else {
		localOptions[paramName] = optionValue
	}

	return nil
}

// Parse global command line arguments.
func ParseGlobal() error {
	iterator.Reset()

	// Parse global flags/options
	for param, ok := iterator.Next(); ok; param, ok = iterator.Next() {
		if strings.HasPrefix(param, "--") {
			err := parseCliParam(iterator, true)
			if err != nil {
				return err
			}
		} else {
			// Found regular Arg. Stop global parsing.
			iterator.Rewind()
			break
		}
	}

	// Load regular args
	for param, ok := iterator.Next(); ok; param, ok = iterator.Next() {
		if strings.HasPrefix(param, "--") {
			// Found a flag/option; stop global parsing
			iterator.Rewind()
			break
		} else {
			args.arguments = append(args.arguments, param)
		}
	}

	return nil
}

// Parses local command line arguments.
func Parse() error {
	for param, ok := iterator.Next(); ok; param, ok = iterator.Next() {
		if strings.HasPrefix(param, "--") {
			err := parseCliParam(iterator, false)
			if err != nil {
				return err
			}
		} else {
			// Should not have arguments at this point.
			return errors.New(fmt.Sprintf(
				"Unexpected argument '%s' after localFlags/localOptions.", param,
			))
		}
	}

	return nil
}
