// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gcp

import (
	"fmt"
	"github.com/vsartor/ro/donna"
	"github.com/vsartor/ro/weems"
	"os"
)

var logger weems.Logger

func initGcp() {
	logger = weems.NewLogger("gcp")
}

func Cmd() {
	initGcp()
	command, ok := donna.NextArg()
	if !ok {
		fmt.Printf("Expected a GCP command.\n")
		os.Exit(1)
	}

	switch command {
	case "cluster":
		clusterCmd()
	case "upload":
		uploadCmd()
	default:
		fmt.Printf("Unexpected GCP command %q.\n", command)
		os.Exit(1)
	}
}
