package pold

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
	Port  int    `yaml:"port:`
}

func NewConfig(path string) (conf Config, err error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, errors.Wrap(err, "err read config file")
	}

	if err := yaml.Unmarshal(data, &conf); err != nil {
		return conf, errors.Wrap(err, "fail mapping config file")
	}

	if conf.Title == "" {
		return conf, fmt.Errorf("please set blog title")
	}

	return conf, nil
}
