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

// Blogo is the tool used to compile my blog's source files into the
// appropriate structure.
//
// The expected format is a folder with five contents:
//  * posts - A directory containing blog posts.
//  * static - A directory containing posts that are static pages.
//  * templates - A directory with HTML templates, which will be managed
//    by the pages submodule.
//  * resources - A directory whose contents will be copied over to the
//    destination root.
//  * config.json - A JSON file containing certain configuration.
//
// Posts are written in Markdown, but with the first five lines being
// prefixed by an exclamation mark and having the following meanings:
//  1) Title
//  2) Date string; e.g. February 1st, 2019
//  3) Post tags, with css colors; e.g. old:lightcoral dev:burlywood
//  4) Preview; a paragraph in Markdown describing the post's contents
//  5) Flags; specific flags used by the parser for special behaviours
//
// Note that for static pages, lines 3 and 4 are empty as they're not
// applicable. The other difference between static and regular posts is
// that posts are indexed in <base_url>/posts.html and the most recent
// posts are indexed at <base_url>/index.html. Posts are all stored in
// <base_url>/posts/<post_name>.html whereas static posts are stored
// at the root, <base_url>/<post_name>.html. Note also that the post
// files are indexed with a prefix of the form "<index>_" on their file
// name.
//
// The configuration holds four different settings:
//  1) title: The blog title, to be embedded in the head section.
//  2) base_url: The blog's base URL.
//  3) image_dir: The name of the resources subdirectory where post
//     images are included.
//  4) num_home_posts: How many of the most recent posts are indexed
//     in the home page.

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
