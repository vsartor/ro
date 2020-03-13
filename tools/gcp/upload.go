// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/vsartor/ro/donna"
	"github.com/vsartor/ro/linus"
	"google.golang.org/api/option"
	"os"
)

func uploadCmd() {
	// Handle command line parameters
	donna.ExpectOption("cred")
	donna.ExpectOption("bucket")
	err := donna.Parse()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// TODO: Get this from Rich
	credential, ok := donna.GetOption("cred")
	if !ok {
		fmt.Println("Did not receive credential file path.")
		os.Exit(1)
	}
	bucket, ok := donna.GetOption("bucket")
	if !ok {
		fmt.Println("Did not receive bucket.")
		os.Exit(1)
	}

	srcPath, ok := donna.NextArg()
	if !ok {
		fmt.Println("Expected source path as argument.")
		os.Exit(1)
	}
	dstPath, ok := donna.NextArg()
	if !ok {
		fmt.Println("Expected destination path as argument.")
		os.Exit(1)
	}

	// Create client

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credential))
	if err != nil {
		logger.Fatal("Could not create client: %s", err)
	}

	// Dispatch to correct routine

	if !linus.Exists(srcPath) {
		logger.Fatal("srcPath does not exist: %s", srcPath)
	}

	if linus.IsFile(srcPath) {
		logger.Info("Uploading a single file.")
		err := uploadFile(srcPath, bucket, dstPath, ctx, client)
		if err != nil {
			logger.Fatal("Could upload file: %s", err)
		}
	} else {
		logger.Info("Uploading a folder.")
		err := uploadDir(srcPath, bucket, dstPath, ctx, client)
		if err != nil {
			logger.Fatal("Could not upload directory: %s", err)
		}
	}
}
