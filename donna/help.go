// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

import (
	"fmt"
	"os"
	"strings"
)

type NamedSlice interface {
	Names() []string
}

func computePadding(namedObjects NamedSlice) int {
	padding := 0
	for _, name := range namedObjects.Names() {
		nameLen := len(name)
		if nameLen > padding {
			padding = nameLen
		}
	}
	return padding
}

func getKindStr(kind ParamKind) string {
	if kind == ParamStr {
		return "str"
	} else if kind == ParamInt {
		return "int"
	} else if kind == ParamFlag {
		return "nil"
	}

	// If we got here, there was a developer error.
	panic("Unexpected parameter kind received.")
}

// Displays command help based on parameter and argument expected information.
func DisplayCommandHelp() {
	// Show command usage
	var usageString strings.Builder

	usageString.WriteString("Usage: ")
	usageString.WriteString(iterator.Path())
	for _, argInfo := range expectedArgs {
		usageString.WriteString(argInfo.name)
		usageString.Write([]byte{' '})
	}
	usageString.WriteString("[parameters]")

	fmt.Println(usageString.String())

	// Arguments section
	if len(expectedArgs) > 0 {
		// Compute padding
		padding := computePadding(expectedArgs)
		padFmt := fmt.Sprintf("%%-%ds", padding)

		// Display arguments
		fmt.Printf("\nArguments:\n")
		for _, argInfo := range expectedArgs {
			paddedName := fmt.Sprintf(padFmt, argInfo.name)
			fmt.Printf("\033[94m→\033[0m %s \033[94m[str]\033[0m: %s\n", paddedName, argInfo.desc)
		}
	}

	// Parameters section
	if len(expectedLocalParams) > 0 {
		// Compute padding
		padding := computePadding(expectedLocalParams)
		padFmt := fmt.Sprintf("%%-%ds", padding)

		// Display parameters
		fmt.Printf("\nParameters:\n")
		for _, paramInfo := range expectedLocalParams {
			kind := getKindStr(paramInfo.kind)
			paddedName := fmt.Sprintf(padFmt, paramInfo.name)
			fmt.Printf("\033[94m→\033[0m %s \033[94m[%s]\033[0m: %s\n", paddedName, kind, paramInfo.desc)
		}
	}

	os.Exit(0)
}
