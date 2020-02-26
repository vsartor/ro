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

func main() {
	donna.Parse()

	// Handle logging level flags; verbose overrides quiet by design
	if donna.HasGlobalFlag("quiet") {
		weems.SetQuiet()
	}
	if donna.HasGlobalFlag("verbose") {
		weems.SetVerbose()
	}

	weems.Info("Ro's method dispatching is beginning.")

	method, ok := donna.NextArg()
	if !ok {
		weems.Fatal("Expected an argument.")
	}

	switch method {
	case "version":
		version()
	case "help":
		help()
	default:
		weems.Fatal("Unexpected argument '%s'.", method)
	}
}

func version() {
	fmt.Printf("I don't know my version yet.\n")
}

func help() {
	fmt.Printf("I can't help you yet.\n")
}
