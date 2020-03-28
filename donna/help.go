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

func getKindStr(kind ParamKind) string {
	if kind == ParamStr {
		return "str"
	} else if kind == ParamInt {
		return "int"
	} else if kind == ParamFlag {
		return "flag"
	}

	// If we got here, there was a developer error.
	panic("Unexpected parameter kind received.")
}

// Displays command help based on parameter and argument expected information.
func DisplayCommandHelp() {
	// Show command usage
	var usageString strings.Builder

	usageString.WriteString("Usage: ...")
	for _, argInfo := range expectedArgs {
		usageString.Write([]byte{' '})
		usageString.WriteString(argInfo.name)
	}
	usageString.WriteString(" [parameters]")

	fmt.Println(usageString.String())

	// Show argument descriptions
	if len(expectedArgs) > 0 {
		fmt.Printf("\nArguments:\n")
		for _, argInfo := range expectedArgs {
			fmt.Printf("→ %s \033[94m[str]\033[0m: %s\n", argInfo.name, argInfo.desc)
		}
	}

	// Show flag descriptions
	if len(expectedLocalParams) > 0 {
		fmt.Printf("\nParameters:\n")
		for _, paramInfo := range expectedLocalParams {
			kind := getKindStr(paramInfo.kind)
			fmt.Printf("→ %s \033[94m[%s]\033[0m: %s\n", paramInfo.name, kind, paramInfo.desc)
		}
	}

	os.Exit(0)
}
