package pold

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

type PostHeader struct {
	Title     string   `yaml:"title"`
	Tag       []string `yaml:"tag"`
	PublishAt string   `yaml:"publish_at"`
}

type Post struct {
	Header     *PostHeader
	ContentStr string
	Content    template.HTML
	Path       string
}

type Posts []*Post

var (
	postSeparatorStr = `---\n+`
	postSeparator    = regexp.MustCompile(postSeparatorStr)
)

func parsePost(postStr string) (*Post, error) {
	splitContent := postSeparator.Split(postStr, 3)

	if len(splitContent) != 3 || splitContent[0] != "" {
		return nil, fmt.Errorf("content format is invalid")
	}

	postHeader := &PostHeader{}
	if err := yaml.Unmarshal([]byte(splitContent[1]), postHeader); err != nil {
		return nil, err
	}

	post := &Post{
		Header:     postHeader,
		ContentStr: splitContent[2],
	}

	return post, nil
}

func NewPost(path string) (*Post, error) {
	filePath := fmt.Sprintf("%s/post%s.md", root, path)

	postFile, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, errors.Wrap(err, "error read path")
	}

	postStr := string(postFile)
	post, err := parsePost(postStr)

	if err != nil {
		return nil, errors.Wrap(err, "fail parse content")
	}

	md := blackfriday.MarkdownCommon([]byte(post.ContentStr))
	content := bluemonday.UGCPolicy().Sanitize(string(md))

	post.Content = template.HTML(content)
	post.Path = path

	return post, nil
}

func PublishedPosts(count int) (Posts, error) {
	postRoot := fmt.Sprintf("%s/post/", root)

	out, _ := exec.Command("git", "log", "--name-only", "--pretty=format:", postRoot).Output()
	files := strings.Split(string(out), "\n")

	var posts Posts
	for _, file := range files {
		if file == "" {
			continue
		}
		path := strings.Replace(file, "post", "", -1)
		path = strings.Replace(path, ".md", "", -1)
		post, err := NewPost(path)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)

		if count <= len(posts) {
			break
		}
	}

	return posts, nil
}

func CreatePost(name string) error {
	_, err := os.Stat("post")

	if err != nil {
		return errors.Wrap(err, "not exits post firectory")
	}

	dir := fmt.Sprintf("post/%s", time.Now().Format("2006/01/02"))
	_, err = os.Stat(dir)

	if err != nil {
		os.MkdirAll(dir, 0755)
	}

	publishAt := time.Now().Format("2006-01-02")
	tmpl := tmpl(publishAt)

	filePath := fmt.Sprintf("%s/%s.md", dir, name)
	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "fail file create")
	}
	defer file.Close()

	_, err = file.Write([]byte(tmpl))
	if err != nil {
		return errors.Wrap(err, "fail write default format")
	}

	return nil
}

func tmpl(publishAt string) string {
	return fmt.Sprintf(`---
title:
tag:
publish_at: %s
---
`, publishAt)
}
