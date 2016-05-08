package pold

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

	var posts Posts
	err := filepath.Walk(postRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			path = strings.Replace(path, fmt.Sprintf("%s/post", root), "", -1)
			path = strings.Replace(path, ".md", "", -1)

			post, err := NewPost(path)

			if err != nil {
				return err
			}

			posts = append(posts, post)

			if count <= len(posts) {
				return nil
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "fail get post files")
	}

	return posts, nil
}
