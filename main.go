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
	"os"
)

var (
	logger  weems.Logger
	Version string
)

func init() {
	logger = weems.NewLogger("ro")
}

func main() {
	// Add expected values
	donna.ExpectGlobalFlag("t", "trace", "Enables extremely verbose logging.")
	donna.ExpectGlobalFlag("v", "verbose", "Enables verbose logging.")
	donna.ExpectGlobalFlag("q", "quiet", "Only log errors.")

	// Parse and validate global values
	err := donna.ParseGlobal()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

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
		fmt.Printf("No commands given.\n")
		os.Exit(1)
	}

	switch method {
	case "version":
		cmdVersion()
	case "blogo":
		logger.Trace("Dispatching to blogo.")
		blogo.Cmd()
	case "gcp":
		logger.Trace("Dispatching to gcp.")
		gcp.Cmd()
	default:
		fmt.Printf("Unexpected command %q.\n", method)
		os.Exit(1)
	}
}

func cmdVersion() {
	fmt.Printf("ro 1r%s\n", Version)
}
