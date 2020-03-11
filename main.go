// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"github.com/vsartor/ro/donna"
	"github.com/vsartor/ro/tools/blogo"
	"github.com/vsartor/ro/tools/gcp"
	"github.com/vsartor/ro/weems"
)

var logger weems.Logger

func init() {
	logger = weems.NewLogger("ro")
}

func main() {
	// Add expected values
	donna.ExpectGlobalFlag("trace")
	donna.ExpectGlobalFlag("verbose")
	donna.ExpectGlobalFlag("quiet")

	// Parse and validate global values
	donna.ParseGlobal()

	// Handle logging level flags; most verbose flags have precedence on purpose
	if donna.HasGlobalFlag("quiet") {
		weems.SetGlobalLevel(weems.ERROR)
	}
	if donna.HasGlobalFlag("verbose") {
		weems.SetGlobalLevel(weems.INFO)
	}
	if donna.HasGlobalFlag("trace") {
		weems.SetGlobalLevel(weems.TRACE)
	}

	method, ok := donna.NextArg()
	if !ok {
		logger.Fatal("Expected an argument.")
	}

	switch method {
	case "version":
		cmdVersion()
	case "blogo":
		logger.Info("Dispatching to blogo.")
		blogo.Cmd()
	case "gcp":
		logger.Info("Dispatching to gcp.")
		gcp.Cmd()
	default:
		logger.Fatal("Unexpected argument '%s'.", method)
	}
}

func cmdVersion() {
	fmt.Printf("ro version %s\n", version)
}
