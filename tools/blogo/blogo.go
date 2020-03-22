// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blogo

import (
	"fmt"
	"github.com/vsartor/ro/donna"
	"github.com/vsartor/ro/weems"
	"os"
)

var logger weems.Logger

// Blogo's initialization function.
// Do not make use of regular 'init' functionality to avoid all
// tools being initialized at startup since only one will be used.
func initBlogo() {
	logger = weems.NewLogger("blogo")
}

// Blogo's command entry point.
func Cmd() {
	initBlogo()

	// Register this commands flags/options and validate them
	donna.ExpectFlag("l", "local")
	err := donna.Parse()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Attempt to fetch srcDir and dstDir through Donna arguments,
	// and if not possible attempt to fetch from OS environment
	// variables.

	srcDir, ok := donna.NextArg()
	if !ok {
		srcDir = os.Getenv(srcDirEnvKey)
		if srcDir != "" {
			fmt.Printf("Fetched source path from enviroment: %q.\n", srcDir)
		}
	}
	if srcDir == "" {
		fmt.Printf("Pass srcDir as argument or through $%s.", srcDirEnvKey)
		os.Exit(1)
	}

	dstDir, ok := donna.NextArg()
	if !ok {
		dstDir = os.Getenv(dstDirEnvKey)
		if dstDir != "" {
			fmt.Printf("Fetched destination path from enviroment: %q.\n", dstDir)
		}
	}
	if dstDir == "" {
		fmt.Printf("Pass dstDir as argument or through $%s.", dstDirEnvKey)
		os.Exit(1)
	}

	// Load settings
	loadSettings(srcDir, dstDir)

	// Handle command line parameters
	if donna.HasFlag("local") {
		redirectUrl(dstDir)
	}

	// Actually perform the blog build
	buildBlog()
}
