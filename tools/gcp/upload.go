// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/vsartor/ro/donna"
	"github.com/vsartor/ro/linus"
	"google.golang.org/api/option"
)

func uploadCmd() {
	// Handle command line parameters

	donna.ExpectOption("cred")
	donna.ExpectOption("bucket")
	donna.ExpectOption("project")
	donna.ValidateLocal()

	credential, ok := donna.GetOption("cred")
	if !ok {
		logger.Fatal("Did not receive credential.")
	}
	bucket, ok := donna.GetOption("bucket")
	if !ok {
		logger.Fatal("Did not receive bucket.")
	}

	srcPath, ok := donna.NextArg()
	if !ok {
		logger.Fatal("Expected source path as argument.")
	}
	dstPath, ok := donna.NextArg()
	if !ok {
		logger.Fatal("Expected destination path as argument.")
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
