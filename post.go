package pold

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

type PostHeader struct {
	Title string   `yaml:"title"`
	Tag   []string `yaml:"tag"`
}

type Post struct {
	Header     *PostHeader
	ContentStr string
	Content    template.HTML
}

type Posts []Post

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
	postFile, err := ioutil.ReadFile(path)

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

	return post, nil
}
