// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

import (
	"os"
	"strings"
)

type paramIterator struct {
	currIdx int
}

// Resets the iterator
func (pi *paramIterator) Reset() {
	pi.currIdx = 0
}

// Rewinds the iterator one step.
func (pi *paramIterator) Rewind() {
	if pi.currIdx > 0 {
		pi.currIdx--
	}
}

// Returns the current parameter, and a flag indicating if this was
// a valid request.
func (pi *paramIterator) Curr() string {
	return os.Args[pi.currIdx]
}

// Returns the chain of Args as a string up to the current command.
func (pi *paramIterator) Path() string {
	fullPath := os.Args[:pi.currIdx+1]

	var argsPath strings.Builder
	for _, arg := range fullPath {
		if !strings.HasPrefix(arg, "-") {
			argsPath.WriteString(arg)
			argsPath.Write([]byte{' '})
		}
	}
	return argsPath.String()
}

// Returns the next parameter, and a flag indicating if there was
// an argument to return.
func (pi *paramIterator) Next() (string, bool) {
	if pi.currIdx+1 == len(os.Args) {
		return "", false
	}
	pi.currIdx++
	return os.Args[pi.currIdx], true
}
