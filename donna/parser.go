// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var iterator paramIterator

// Parses a single cli parameter.
func parseCliParam(global bool) error {
	// Get current parameter and trim dashes.
	param := iterator.Curr()
	paramName := strings.Trim(param, "-")

	// Check whether the argument is valid, and if so check whether
	// it's a boolean flag.
	paramInfo, isExpected := expectedInfo(paramName, global)
	if !isExpected {
		errorMsg := fmt.Sprintf("Unexpected flag %q.", param)
		return errors.New(errorMsg)
	}

	if paramInfo.kind == ParamFlag {
		// Since parsing a flag is simpler, handle this branch earlier.
		var paramRef *ParamInfo
		if global {
			paramRef = globalParams[paramInfo.name]
		} else {
			paramRef = localParams[paramInfo.name]
		}
		paramRef.ToggleFlag()

		return nil
	}

	// Assert that option received an associated value.
	optionVal, ok := iterator.Next()
	if !ok {
		errorMsg := fmt.Sprintf("Option %q did not receive an associated value.", paramInfo.name)
		return errors.New(errorMsg)
	}

	if paramInfo.kind == ParamInt {
		// Validate the parsedVal type
		parsedVal, err := strconv.Atoi(optionVal)
		if err != nil {
			errorMsg := fmt.Sprintf(
				"Option %q did not receive a valid integer parsedVal.", paramInfo.name,
			)
			return errors.New(errorMsg)
		}

		var paramRef *ParamInfo
		if global {
			paramRef = globalParams[paramInfo.name]
		} else {
			paramRef = localParams[paramInfo.name]
		}
		paramRef.SetIntValue(parsedVal)

		return nil
	}

	// String case, nothing to do but register the value.
	var paramRef *ParamInfo
	if global {
		paramRef = globalParams[paramInfo.name]
	} else {
		paramRef = localParams[paramInfo.name]
	}
	paramRef.SetStrValue(optionVal)

	return nil
}

// Parse global command line arguments.
func ParseGlobal() error {
	iterator.Reset()

	// Initialize Parameter information with default values based
	// on expected parameter information.
	for _, expectedParam := range expectedGlobalParams {
		globalParams[expectedParam.name] = NewParamInfo(expectedParam)
	}

	// Parse global flags/options.
	for param, ok := iterator.Next(); ok; param, ok = iterator.Next() {
		if strings.HasPrefix(param, "-") {
			err := parseCliParam(true)
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
		if strings.HasPrefix(param, "-") {
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
	// Initialize Parameter information with default values based
	// on expected parameter information.
	for _, expectedParam := range expectedLocalParams {
		localParams[expectedParam.name] = NewParamInfo(expectedParam)
	}

	// Parse command specific parameters
	for param, ok := iterator.Next(); ok; param, ok = iterator.Next() {
		if strings.HasPrefix(param, "-") {
			err := parseCliParam(false)
			if err != nil {
				return err
			}
		} else {
			// Should not have arguments at this point.
			errorMsg := fmt.Sprintf(
				"Unexpected argument %q after flags/options.", param,
			)
			return errors.New(errorMsg)
		}
	}

	return nil
}
