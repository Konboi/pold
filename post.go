package pold

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
)

type Post struct {
	Title   string
	Content template.HTML
}

type Posts []Post

type PostHeader struct {
	Title string `yaml:"title"`
	Tag   string `yaml:"tag"`
}

var (
	headerSeparatorStr = `---\n+`
	headerSeparator    = regexp.MustCompile(headerSeparatorStr)
)

func parseHeader(contentStr string) (*PostHeader, error) {
	header := headerSeparator.Split(contentStr, 3)

	if len(header) != 3 || header[0] != "" {
		return nil, fmt.Errorf("header format is invalid")
	}

	postHeader := &{}

	return nil, nil
}

func NewPost(path string) (*Post, error) {
	postFile, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.Wrap(err, "error read path")
	}

	contentStr := string(postFile)
	md := blackfriday.MarkdownCommon(postFile)

	_, err = parseHeader(contentStr)

	if err != nil {
		return nil, errors.Wrap(err, "fail parse header")
	}

	content := bluemonday.UGCPolicy().Sanitize(string(md))

	post := &Post{
		Title:   "",
		Content: template.HTML(content),
	}

	return post, nil
}
