// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blogo

import (
	"bytes"
	"os/exec"
	"strings"
)

// This function serves as a wrapper to external conversion methods to
// facilitate swapping out to alternative in the future.
func markdownToHtml(markdownContent string) string {
	// Pandoc just works really well, although there are other, pure Go
	// alternatives such as gomarkdown/markdown.

	// Invoke pandoc to convert from Markdown to HTML.
	command := exec.Command("pandoc", "-f", "markdown", "-t", "html")

	// Pass the Markdown content to the command's standard input
	command.Stdin = strings.NewReader(markdownContent)

	// Save command outputs into a driver
	var commandOutput bytes.Buffer
	command.Stdout = &commandOutput
	command.Stderr = &commandOutput

	// Run commands and log a Fatal message on errors
	err := command.Run()
	if err != nil {
		logger.Fatal("Pandoc error: %s", err)
	}

	return commandOutput.String()
}
