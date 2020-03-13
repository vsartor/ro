// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package linus

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Copies a single regular file from srcPath to dstPath, preserving
// file permissions.
func CopyFile(srcPath, dstPath string) error {
	// Attempt to open an already existing file on srcPath
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Attempt to create/overwrite on dstPath
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Make sure we can stat the file before attempting to perform
	// the copy operation to avoid persisting with an errored state
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	// Copy file contents from one file to another
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	// Copy file permissions from one file to another
	err = os.Chmod(dstPath, srcInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

// Recursively copies the content of srcPath into a new directory
// at dstPath.
func CopyDir(srcPath, dstPath string) error {
	// Make sure source path can be stat'd and is a directory.
	src, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	if !src.IsDir() {
		return fmt.Errorf("not a directory: %s", srcPath)
	}

	// Make sure that the destination directory does not yet exist.
	// This means that stat'ing this path should raise a "does not
	// exist" error.
	_, err = os.Stat(dstPath)
	if err == nil {
		return fmt.Errorf("path already exists: %s", dstPath)
	}
	if !os.IsNotExist(err) {
		// Error is not nil but it's something other than "does not
		// exist".
		return err
	}

	// Create the directory preserving the directory permissions.
	// Note that we do not use MkdirAll to prevent creating an
	//entire chain with incorrect permissions.
	err = os.Mkdir(dstPath, src.Mode())
	if err != nil {
		return err
	}

	// Fetch the name of all contents from source path.
	contents, err := ioutil.ReadDir(srcPath)
	if err != nil {
		return err
	}

	// Copy each content. If a directory, recurse, if a file, defer
	// to the individual copyFile function.
	for _, content := range contents {
		contentSrcPath := filepath.Join(srcPath, content.Name())
		contentDstPath := filepath.Join(dstPath, content.Name())

		if content.IsDir() {
			err = CopyDir(contentSrcPath, contentDstPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(contentSrcPath, contentDstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Returns whether path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// Returns whether the path leads to a directory.
func IsDir(path string) bool {
	stat, err := os.Stat(path)

	if err != nil {
		return false
	}

	return stat.IsDir()
}

// Returns whether the path leads to a regular file.
func IsFile(path string) bool {
	stat, err := os.Stat(path)

	if err != nil {
		return false
	}

	return !stat.IsDir()
}
