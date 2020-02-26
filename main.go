// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"ro/donna"
	"ro/weems"
)

var logger weems.Logger

func init() {
	logger = weems.NewLogger("ro")
}

func main() {
	// Add expected values
	donna.ExpectGlobalFlag("verbose")
	donna.ExpectGlobalFlag("quiet")

	// Parse and validate global values
	donna.Parse()
	name, ok := donna.ValidateGlobal()
	if !ok {
		logger.Fatal("Unexpected global flag/option '%s'.", name)
	}

	// Handle logging level flags; verbose overrides quiet by design
	if donna.HasGlobalFlag("quiet") {
		weems.SetGlobalLevel(weems.ERROR)
		logger.SetLevel(weems.ERROR)
	}
	if donna.HasGlobalFlag("verbose") {
		weems.SetGlobalLevel(weems.INFO)
		logger.SetLevel(weems.INFO)
	}

	method, ok := donna.NextArg()
	if !ok {
		logger.Fatal("Expected an argument.")
	}

	switch method {
	case "version":
		cmdVersion()
	default:
		logger.Fatal("Unexpected argument '%s'.", method)
	}
}

func cmdVersion() {
	fmt.Printf("ro version %s\n", version)
}
