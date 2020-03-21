// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bran

import (
	"encoding/json"
	"github.com/vsartor/ro/nemec"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Bran remember things for Ro.

// Name of the folder where memories are kept.
const branFolder = "cache"

type cacheMap = map[string]string

var loadedMaps map[string]*cacheMap

func init() {
	loadedMaps = make(map[string]*cacheMap)
}

func hashValue(name string) string {
	var value int
	for _, ch := range name {
		value += int(ch)
		value *= int(ch)
	}
	if value < 0 {
		value *= -1
	}
	str := strconv.Itoa(value)
	if len(str) < 2 {
		return str
	}
	return str[:2]
}

func hashPath(path string) string {
	var hash strings.Builder

	parts := strings.Split(path, string(filepath.Separator))

	for _, part := range parts {
		if part == "" {
			continue
		}
		hash.WriteString(hashValue(part))
	}

	return hash.String()
}

// Returns an identifier for the file that called this function.
func callerTag() string {
	_, filePath, _, _ := runtime.Caller(2)
	return hashPath(filePath)
}

func getMap(caller string) (*cacheMap, error) {
	callerMap, ok := loadedMaps[caller]
	if !ok {
		// Allocate a new map and register it
		newMap := make(cacheMap)
		loadedMaps[caller] = &newMap
		callerMap = &newMap

		// Read file
		filePath := filepath.Join(branFolder, caller) + ".json"
		file, err := nemec.File(filePath, "{}", os.O_RDONLY)
		if err != nil {
			return nil, err
		}
		contents, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		// Unmarshal into new map
		err = json.Unmarshal(contents, callerMap)
		if err != nil {
			return nil, err
		}
	}
	return callerMap, nil
}

func setMap(caller string, callerMap *cacheMap) error {
	filePath := filepath.Join(branFolder, caller) + ".json"
	file, err := nemec.File(filePath, "{}", os.O_WRONLY)
	if err != nil {
		return err
	}
	contents, err := json.Marshal(callerMap)
	if err != nil {
		return err
	}
	_, err = file.Write(contents)
	if err != nil {
		return err
	}
	return nil
}

func Get(key string) (string, error) {
	caller := callerTag()
	callerMap, err := getMap(caller)
	if err != nil {
		return "", err
	}
	value, _ := (*callerMap)[key]
	return value, nil
}

func Set(key, value string) error {
	caller := callerTag()
	callerMap, err := getMap(caller)
	if err != nil {
		return err
	}
	(*callerMap)[key] = value
	err = setMap(caller, callerMap)
	if err != nil {
		return err
	}
	return nil
}
