// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func uploadFile(
	srcPath string,
	bucketName string,
	bucketPath string,
	ctx context.Context,
	client *storage.Client,
) error {
	logger.Info("Uploading '%s' to '%s:%s'.", srcPath, bucketName, bucketPath)

	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	dst := client.Bucket(bucketName).Object(bucketPath).NewWriter(ctx)
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	if err := dst.Close(); err != nil {
		return err
	}

	return nil
}

func uploadDir(srcPath, bucketName, bucketPath string, ctx context.Context, client *storage.Client) error {
	logger.Info("Uploading directory '%s' to '%s:%s'.", srcPath, bucketName, bucketPath)

	err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		// Ignore directories
		if !info.IsDir() {
			relPath := strings.TrimPrefix(path, srcPath)
			storagePath := bucketPath + relPath
			return uploadFile(path, bucketName, storagePath, ctx, client)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
