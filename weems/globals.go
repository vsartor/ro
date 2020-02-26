// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package weems

import (
	"io"
	"os"
	"sync"
)

var (
	globalMutex  sync.Mutex
	globalLevel            = WARN
	globalWriter io.Writer = os.Stderr
)

func SetGlobalLevel(level int) {
	globalMutex.Lock()
	globalLevel = level
	globalMutex.Unlock()
}

func SetGlobalWriter(writer io.Writer) {
	globalMutex.Lock()
	globalWriter = writer
	globalMutex.Unlock()
}
