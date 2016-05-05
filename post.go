package pold

import (
	"html/template"
	"io/ioutil"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
)

func NewPost(path string) (post *Post, err error) {
	postFile, err := ioutil.ReadFile(path)

	if err != nil {
		return post, errors.Wrap(err, "error read path")
	}

	common := blackfriday.MarkdownCommon(postFile)
	content := bluemonday.UGCPolicy().Sanitize(string(common))

	post = &Post{
		Title:   "",
		Content: template.HTML(content),
	}

	return post, nil
}
