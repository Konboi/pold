package pold

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// TODO
// Use template at generate file

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

	files["base"] = base()
	files["index"] = index()
	files["archive"] = archive()
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

func base() string {
	return `<!doctype html>
<html lang="ja">
  <head>
    <meta charset="UTF-8"/>
    {{ block header }}
    {{ end }}
  </head>
  <body>
    <h1>Hello Pold!!!</h1>
    {{ block content }}
    {{ end }}
    {{ block footer }}
    {{ end }}
  </body>
</html>
`
}

func archive() string {
	return `{{ define "archive" }}
{{ template "header" . }}
<h1> archive page </h1>
{{ range .Posts }}
<p><a href="/post/{{ .Path }}">{{ .Header.Title }}</a></p>
{{ end }}
{{ template "footer" . }}
{{ end }}

`
}

func index() string {
	return `{{ extends "base.html" }}

{{ block header }}
<title>{{ .Blog.Title }}</title>
{{ end }}

{{ block content }}
{{ range .Posts }}
<h1> {{ .Header.Title }} </h1>

<section>
  {{ .Content }}
</section>
{{ end }}
{{ end }}
`
}

func post() string {
	return ``
}

func conf() string {
	return `title: SampleTitle
author: Sample
url: 'http://pold.example.com'
port: 8765
`
}
