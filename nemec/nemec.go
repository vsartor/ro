// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package nemec

import (
	"errors"
	"fmt"
	"github.com/vsartor/ro/linus"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Nemec handles storage of internal folders and files for Ro.

const (
	envRoDirPath = "RO_DIR"
	defaultRoDir = "$HOME/.ro"
)

// CompletePath to Ro's root folder.
var rootPath string

// Initializes the path
func init() {
	envPath := os.Getenv(envRoDirPath)
	if envPath != "" {
		rootPath = envPath
	} else {
		rootPath = defaultRoDir
	}
	rootPath = os.ExpandEnv(rootPath)

}

// Returns complete path given a Ro relative path.
func CompletePath(path string) string {
	return filepath.Join(rootPath, path)
}

// Opens a file, creating it if it doesn't exist.
// The path is relative to Ro's internal folder.
func File(path, emptyValue string, flag int) (*os.File, error) {
	// Get complete path
	path = CompletePath(path)

	if !linus.Exists(path) {
		// Create directory where the file should be located
		dirPath := filepath.Dir(path)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return nil, err
		}

		// Create file
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		_, err = file.WriteString(emptyValue)
		if err != nil {
			return nil, err
		}
		file.Close()
	} else if linus.IsDir(path) {
		// If it exists but is a directory, give out an error
		errMsg := fmt.Sprintf("cannot read directory %s as a file", path)
		return nil, errors.New(errMsg)
	}

	return os.OpenFile(path, flag, 0644)
}

// Lists elements of a directory, creating it if it doesn't exist.
// The path is relative to Ro's internal folder.
func DirListing(path string) ([]os.FileInfo, error) {
	// Get complete path
	path = CompletePath(path)

	if !linus.Exists(path) {
		// Create directory if it doesn't exist
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return nil, err
		}
	} else if linus.IsFile(path) {
		errMsg := fmt.Sprintf("cannot read file %s as a directory", path)
		return nil, errors.New(errMsg)
	}

	return ioutil.ReadDir(path)
}
