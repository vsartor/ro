// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blogo

import (
	"fmt"
	"strings"
)

type tagsInfo struct {
	tags []tagInfo
}

type tagInfo struct {
	tagName  string
	tagColor string
}

const tagTemplate = `<div class="post-tag" style="background-color: %s;"><div class="tag-content">%s</div></div>`

func newTag(tagSpec string) tagInfo {
	fields := strings.Split(tagSpec, ":")
	return tagInfo{
		tagName:  fields[0],
		tagColor: fields[1],
	}
}

func newTags(tagSpecs []string) tagsInfo {
	tags := make([]tagInfo, len(tagSpecs))
	for idx, tagSpec := range tagSpecs {
		tags[idx] = newTag(tagSpec)
	}
	return tagsInfo{tags}
}

func (tag tagInfo) toHTML() string {
	return fmt.Sprintf(tagTemplate, tag.tagColor, tag.tagName)
}

func (tags tagsInfo) toHTML() string {
	var html strings.Builder

	for _, tag := range tags.tags {
		html.WriteString(tag.toHTML())
		html.WriteString("\n")
	}

	return html.String()
}
