// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gcp

import (
	"github.com/vsartor/ro/donna"
	"github.com/vsartor/ro/weems"
)

var logger weems.Logger

func initGcp() {
	logger = weems.NewLogger("gcp")
}

func Cmd() {
	initGcp()

	donna.ForgetDispatch()
	donna.ExpectDispatch("cluster", "Create a Dataproc cluster.")
	donna.ExpectDispatch("upload", "Upload a file or directory to a bucket.")

	command, ok := donna.NextArg()
	if !ok {
		donna.DisplayDispatchHelp()
	}

	switch command {
	case "cluster":
		clusterCmd()
	case "upload":
		uploadCmd()
	default:
		donna.DisplayDispatchHelp()
	}
}
