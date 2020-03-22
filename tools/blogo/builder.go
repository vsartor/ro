// ro - Copyright (c) Victhor Sartório, 2020-. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blogo

import (
	"github.com/vsartor/ro/linus"
	"github.com/vsartor/ro/tools/blogo/pages"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type commonResources struct {
	head   pages.Page
	header pages.Page
	footer pages.Page
}

func getCommonPage(name string) pages.Page {
	page, err := pages.NewPage(templatePath(name))
	if err != nil {
		logger.Fatal("Could not fetch %s: %s", page, err)
	}

	page.Replace("page_title", settings.Title)
	page.Replace("base_url", settings.BaseUrl)

	return page
}

func getCommonResources() commonResources {
	return commonResources{
		head:   getCommonPage("head"),
		header: getCommonPage("content_header"),
		footer: getCommonPage("footer"),
	}
}

func applyCommonResources(page *pages.Page, resources commonResources) {
	page.Replace("head", resources.head.ToString())
	page.Replace("content_header", resources.header.ToString())
	page.Replace("footer", resources.footer.ToString())
}

func getPostInfo(name string, isStatic bool, c chan postInfo) {
	logger.Info("Building post page %s.", name)

	postInfo, err := parsePost(name, isStatic)
	if err != nil {
		logger.Error("Skipping %s due to: %s", name, err)
	}

	c <- postInfo
}

func getPosts() []postInfo {
	posts := make([]postInfo, 0)

	postSrcs, err := ioutil.ReadDir(postsPath())
	if err != nil {
		logger.Fatal("Could not list posts: %s", err)
	}

	ch := make(chan postInfo)
	for _, postSrc := range postSrcs {
		go getPostInfo(postSrc.Name(), false, ch)
	}

	for range postSrcs {
		posts = append(posts, <-ch)
	}

	// Sort entries, most recent first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].index > posts[j].index
	})

	return posts
}

func getStatics() []postInfo {
	statics := make([]postInfo, 0)

	staticsBasePath := filepath.Join(settings.SrcPath, "static")
	staticSrcs, err := ioutil.ReadDir(staticsBasePath)
	if err != nil {
		logger.Fatal("Could not list statics: %s", err)
	}

	ch := make(chan postInfo)
	for _, staticSrc := range staticSrcs {
		go getPostInfo(staticSrc.Name(), true, ch)
	}

	for range staticSrcs {
		statics = append(statics, <-ch)
	}

	// Sort entries, most recent first
	sort.Slice(statics, func(i, j int) bool {
		return statics[i].index > statics[j].index
	})

	return statics
}

func buildResources() {
	logger.Info("Building resources.")

	// Copy all resources to the destination folder
	srcResourcesPath := filepath.Join(settings.SrcPath, "resources")
	err := linus.CopyDir(srcResourcesPath, settings.DstPath)
	if err != nil {
		logger.Fatal("Could not build resources: %s", err)
	}
}

func buildHome(posts []postInfo, resources commonResources) {
	logger.Info("Building home page.")

	// Create page and add common resources
	index, err := pages.NewPage(templatePath("index"))
	if err != nil {
		logger.Fatal("Error reading template: %s", err)
	}
	applyCommonResources(&index, resources)

	// Add post previews by most recent
	postContent := pages.Page{}
	previewCount := 0
	for _, post := range posts {
		if !post.isHome {
			continue
		}

		logger.Info("Adding preview of %q to home.", post.title)
		previewCount++
		postContent.Append(post.ToHtml(true))

		if previewCount == settings.NumHomePosts {
			break
		}
	}
	index.Replace("post_content", postContent.ToString())

	// Write file
	dstPath := filepath.Join(settings.DstPath, "index.html")
	err = index.Write(dstPath)
	if err != nil {
		logger.Fatal("Could not write index: %s", err)
	}
}

func postFileName(post postInfo) string {
	urlTerms := strings.Split(post.url, "/")
	return urlTerms[len(urlTerms)-1]
}

func buildPost(post postInfo, isStatic bool, resources commonResources) error {
	logger.Info("Generating page for %q.", post.title)

	postPage, err := pages.NewPage(templatePath("post"))
	if err != nil {
		return err
	}
	applyCommonResources(&postPage, resources)

	postPage.Replace("post_content", post.ToHtml(false))
	postPage.ReplaceEncoded("image_dir", settings.BaseUrl+"/"+settings.ImgDir)

	fileName := postFileName(post)

	var filePath string
	if isStatic {
		filePath = filepath.Join(settings.DstPath, fileName)
	} else {
		filePath = filepath.Join(settings.DstPath, "posts", fileName)
	}

	err = postPage.Write(filePath)
	if err != nil {
		return err
	}

	return nil
}

func buildPosts(posts []postInfo, statics bool, resources commonResources) {
	if statics {
		logger.Info("Building static posts.")
	} else {
		logger.Info("Building posts.")
	}

	// For actual posts, create the subfolder they'll be contained
	if !statics {
		err := os.Mkdir(filepath.Join(settings.DstPath, "posts"), os.ModePerm)
		if err != nil {
			logger.Fatal("Could not create posts folder: %s", err)
		}
	}

	for _, post := range posts {
		err := buildPost(post, statics, resources)
		if err != nil {
			logger.Error("Skipping %q due to: %s", post.title, err)
		}
	}
}

func buildPostsPage(posts []postInfo, resources commonResources) {
	logger.Info("Building posts index.")

	// Create page and add common resources
	postsIndex, err := pages.NewPage(templatePath("posts"))
	if err != nil {
		logger.Fatal("Error reading template: %s", err)
	}
	applyCommonResources(&postsIndex, resources)

	// Add post previews by most recent
	postContent := pages.Page{}
	for _, post := range posts {
		logger.Info("Adding preview of %q to posts index.", post.title)
		postContent.Append(post.ToHtml(true))
	}
	postsIndex.Replace("post_index", postContent.ToString())

	// Write file
	dstPath := filepath.Join(settings.DstPath, "posts.html")
	err = postsIndex.Write(dstPath)
	if err != nil {
		logger.Fatal("Could not write posts index: %s", err)
	}
}

func buildBlog() {
	// Load basic resources and parse posts before any file operations are made
	resources := getCommonResources()
	posts := getPosts()
	statics := getStatics()

	// Make sure target folder is cleared in case it exists
	if linus.Exists(settings.DstPath) {
		err := os.RemoveAll(settings.DstPath)
		if err != nil {
			logger.Fatal("could not remove target folder: %s", settings.DstPath)
		}
	}

	// Run all the build steps
	buildResources()
	buildHome(posts, resources)
	buildPosts(posts, false, resources)
	buildPosts(statics, true, resources)
	buildPostsPage(posts, resources)

	logger.Info("Successfully finished building the blog.")
}
