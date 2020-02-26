// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pages

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Holds the content of a page.
type Page struct {
	content string
}

// Creates a new page by reading a specific file.
func NewPage(path string) (Page, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return Page{}, err
	}

	return Page{string(contents)}, nil
}

// Replaces a certain tag from a Page by actual content.
func (page *Page) Replace(tagName, newContent string) {
	realTag := fmt.Sprintf("{{%s}}", tagName)
	page.content = strings.ReplaceAll(page.content, realTag, newContent)
}

// Replaces a certain tag from a Page, with the delimiters HTML encoded,
// by actual content.
func (page *Page) ReplaceEncoded(tagName, newContent string) {
	realTag := fmt.Sprintf("%%7B%%7B%s%%7D%%7D", tagName)
	page.content = strings.ReplaceAll(page.content, realTag, newContent)
}

// Appends content to an existing Page.
func (page *Page) Append(newContent string) {
	page.content = fmt.Sprintf("%s\n%s", page.content, newContent)
}

// Returns a string with the Page's content.
func (page *Page) ToString() string {
	return page.content
}

// Writes the content of a Page into a file.
func (page *Page) Write(path string) error {
	err := ioutil.WriteFile(path, []byte(page.content), 0644)
	if err != nil {
		return err
	}
	return nil
}

