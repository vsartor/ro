// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

var (
	globalOptions map[string]string
	localOptions  map[string]string
)

func init() {
	globalOptions = make(map[string]string)
	localOptions = make(map[string]string)
}

func GetGlobalOption(name string) (string, bool) {
	value, exists := globalOptions[name]
	return value, exists
}

func GetOption(name string) (string, bool) {
	value, exists := localOptions[name]
	return value, exists
}
