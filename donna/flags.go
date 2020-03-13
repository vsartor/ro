// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

var (
	globalFlags []string
	localFlags  []string
)

func init() {
	globalFlags = make([]string, 0)
	localFlags = make([]string, 0)
}

func HasGlobalFlag(name string) bool {
	for _, flag := range globalFlags {
		if flag == name {
			return true
		}
	}
	return false
}

func HasFlag(name string) bool {
	for _, flag := range localFlags {
		if flag == name {
			return true
		}
	}
	return false
}
