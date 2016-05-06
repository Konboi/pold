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
	Header  *PostHeader
	Content template.HTML
}

type Posts []Post

var (
	headerSeparatorStr = `---\n+`
	headerSeparator    = regexp.MustCompile(headerSeparatorStr)
)

func parseHeader(contentStr string) (*PostHeader, error) {
	header := headerSeparator.Split(contentStr, 3)

	if len(header) != 3 || header[0] != "" {
		return nil, fmt.Errorf("header format is invalid")
	}

	postHeader := &PostHeader{}
	if err := yaml.Unmarshal([]byte(header[1]), postHeader); err != nil {
		return nil, err
	}

	return postHeader, nil
}

func parseContent(contentStr string) (string, error) {
	content := headerSeparator.Split(contentStr, 3)

	if len(content) != 3 || content[0] != "" {
		return "", fmt.Errorf("content format is invalid")
	}

	return content[2], nil
}

func NewPost(path string) (*Post, error) {
	postFile, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.Wrap(err, "error read path")
	}

	contentStr := string(postFile)

	header, err := parseHeader(contentStr)
	if err != nil {
		return nil, errors.Wrap(err, "fail parse header")
	}

	parseContentStr, err := parseContent(contentStr)

	if err != nil {
		return nil, errors.Wrap(err, "fail parse content")
	}

	md := blackfriday.MarkdownCommon([]byte(parseContentStr))
	content := bluemonday.UGCPolicy().Sanitize(string(md))

	post := &Post{
		Header:  header,
		Content: template.HTML(content),
	}

	return post, nil
}
