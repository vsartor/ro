// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blogo

import (
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
	donna.ExpectFlag("local")
	donna.ValidateLocal()

	// Attempt to fetch srcDir and dstDir through Donna arguments,
	// and if not possible attempt to fetch from OS environment
	// variables.

	srcDir, ok := donna.NextArg()
	if !ok {
		srcDir = os.Getenv(srcDirEnvKey)
		if srcDir != "" {
			logger.Info("Fetched srcDir from environment as '%s'.", srcDir)
		}
	}
	if srcDir == "" {
		logger.Fatal("Pass srcDir as argument or through $%s.", srcDirEnvKey)
	}

	dstDir, ok := donna.NextArg()
	if !ok {
		dstDir = os.Getenv(dstDirEnvKey)
		if dstDir != "" {
			logger.Info("Fetched dstDir from environment as '%s'.", dstDir)
		}
	}
	if dstDir == "" {
		logger.Fatal("Pass dstDir as argument or through $%s.", dstDirEnvKey)
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
