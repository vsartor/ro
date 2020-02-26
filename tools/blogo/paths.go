// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blogo

import "path/filepath"

func postsPath() string {
	return filepath.Join(settings.SrcPath, "posts")
}

func templatePath(name string) string {
	return filepath.Join(settings.SrcPath, "templates", name+".html")
}
