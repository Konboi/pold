package pold

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

func Init() error {

	err := createDir()
	if err != nil {
		return errors.Wrap(err, "error create post dir")
	}

	err = createFile()
	if err != nil {
		return errors.Wrap(err, "error create file")
	}

	return nil
}

func createDir() error {
	postRoot := fmt.Sprintf("%s/post", root)

	err := os.Mkdir(postRoot, 0755)

	if err != nil {
		return err
	}

	templateRoot := fmt.Sprintf("%s/templates", root)
	err = os.Mkdir(templateRoot, 0755)

	if err != nil {
		return err
	}

	return nil
}

func createFile() error {
	files := make(map[string]string)

	files["archive"] = archive()
	files["footer"] = footer()
	files["header"] = header()
	files["index"] = index()
	files["post"] = post()

	for k, v := range files {
		file := []byte(v)
		err := ioutil.WriteFile(fmt.Sprintf("%s/templates/%s.html", root, k), file, 0644)

		if err != nil {
			return err
		}
	}

	config := []byte(conf())
	err := ioutil.WriteFile(fmt.Sprintf("%s/pold.yml", root), config, 0644)
	if err != nil {
		return err
	}

	return nil
}

func archive() string {
	return `{{ define "archive" }}
<h1> archive page </h1>

{{ range .Posts }}
<p><a href="/post/{{ .Path }}">{{ .Header.Title }}</a></p>
{{ end }}

{{ end }}
`
}

func footer() string {
	return `{{ define "footer" }}
</body>
</html>
{{ end }}
`
}

func header() string {
	return `{{ define "header" }}
<!doctype html>
<html lang="js">
  <head>
    <meta charset="UTF-8"/>
    <title>
      {{ .Blog.Title }}
    </title>
  </head>
  <body>
{{ end }}
`
}

func index() string {
	return `{{ define "index" }}
{{ template "header" . }}

<h1>{{ .Blog.Title }}</h1>

{{ range .Posts }}
<h1> {{ .Header.Title }} </h1>

<section>
  {{ .Content }}
</section>

{{ end }}

{{ template "footer" . }}

{{ end }}
`
}

func post() string {
	return `{{ define "post" }}
{{ template "header" . }}
<h1>{{ .Post.Header.Title }}</h1>
{{ .Post.Content }}

{{ template "footer" . }}

{{ end }}
`
}

func conf() string {
	return `title: SampleTitle
author: Sample
url: 'http://pold.example.com'
port: 8765
`
}
