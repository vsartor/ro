// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blogo

import (
	"fmt"
	"github.com/vsartor/ro/linus"
	"github.com/vsartor/ro/tools/blogo/pages"
	"path/filepath"
	"strconv"
	"strings"
)

type postInfo struct {
	index   int
	isHome  bool
	tags    tagsInfo
	url     string
	title   string
	date    string
	preview string
	content string
}

func (post postInfo) ToHtml(preview bool) string {
	page, err := pages.NewPage(templatePath("post_content"))
	if err != nil {
		logger.Fatal("Could not load post_content: %s", err)
	}

	page.Replace("post_url", post.url)
	page.Replace("post_title", post.title)
	page.Replace("post_date", post.date)
	page.Replace("tags", post.tags.toHTML())

	if preview {
		page.Replace("post_content", post.preview)
	} else {
		page.Replace("post_content", post.content)
	}

	return page.ToString()
}

// Applies necessary treatments to Markdown content before further
// transformations can be made on them.
func treatContent(content string) string {
	// Currently, the only treatment necessary is doubling up the number of
	// slashes  in the content, since if we do not do it, the slashes that
	// should be present in the output will be interpreted as escape characters
	// instead.
	return strings.ReplaceAll(content, "\\", "\\\\")
}

// Converts Markdown content into an HTML string.
func parseContent(content string) string {
	treatedContent := treatContent(content)
	return markdownToHtml(treatedContent)
}

// Parses information regarding a post.
func parsePost(name string, isStatic bool) (postInfo, error) {
	var directory string
	if isStatic {
		directory = "static"
	} else {
		directory = "posts"
	}

	path := filepath.Join(settings.SrcPath, directory, name)
	lines, err := linus.ReadLines(path)
	if err != nil {
		return postInfo{}, err
	}

	// Fetch title, date
	title := strings.TrimSpace(lines[0][1:])
	date := strings.TrimSpace(lines[1][1:])

	// Fetch tags present in the fourth line
	tags := newTags(strings.Fields(lines[2])[1:])

	// Fetch post preview
	preview := parseContent(strings.TrimSpace(lines[3][1:]))

	// Parse any flags from the fifth line
	isHome := strings.Contains(lines[4], "Home")

	// The rest of the file is regular content
	content := parseContent(strings.Join(lines[5:], "\n"))

	var url string
	var index int
	if isStatic {
		// Split to remove extension
		urlName := strings.Split(name, ".")[0]
		url = fmt.Sprintf("%s/%s.html", settings.BaseUrl, urlName)
	} else {
		indexString := strings.Split(name, "_")[0]
		index, err = strconv.Atoi(indexString)
		if err != nil {
			return postInfo{}, fmt.Errorf("invalid index: %s", err)
		}

		urlName := strings.Split(strings.TrimPrefix(name, indexString+"_"), ".")[0]
		url = fmt.Sprintf("%s/posts/%s.html", settings.BaseUrl, urlName)
	}

	return postInfo{
		index,
		isHome,
		tags,
		url,
		title,
		date,
		preview,
		content,
	}, nil
}
