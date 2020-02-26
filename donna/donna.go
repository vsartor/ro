// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package donna

import "github.com/vsartor/ro/weems"

// Donna handles command line arguments for Ro.

var logger weems.Logger

func init() {
	logger = weems.NewLogger("donna")
}
