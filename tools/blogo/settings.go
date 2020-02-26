// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blogo

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const (
	configFileName = "config.json"
	srcDirEnvKey   = "BLOGO_SOURCE_DIRECTORY"
	dstDirEnvKey   = "BLOGO_TARGET_DIRECTORY"
)

type settingsContext struct {
	Title        string `json:"title"`
	BaseUrl      string `json:"base_url"`
	ImgDir       string `json:"image_dir"`
	NumHomePosts int    `json:"num_home_posts"`
	SrcPath      string
	DstPath      string
}

var settings settingsContext

func loadSettings(sourceDir, targetDir string) {
	configPath := filepath.Join(sourceDir, configFileName)
	fileContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Fatal("Could not open configuration file '%s': %s", configPath, err)
	}

	err = json.Unmarshal(fileContent, &settings)
	if err != nil {
		logger.Fatal("Could not parse configuration in '%s': %s", configPath, err)
	}
	settings.SrcPath = sourceDir
	settings.DstPath = targetDir
}

func redirectUrl(newUrl string) {
	logger.Warn("Redirecting build URL to '%s'.", newUrl)
	settings.BaseUrl = newUrl
}
