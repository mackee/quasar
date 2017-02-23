package quasar

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type config struct {
	Daemons []daemon `yaml:"daemons"`
}

func ParseConfig(filename string) (config, error) {
	src, err := os.Open(filename)
	if err != nil {
		return config{}, errors.Wrap(err, "cannot read cofig")
	}
	c, err := parseConfig(src)
	if err != nil {
		return config{}, errors.Wrap(err, "fail parse config")
	}

	return c, nil
}

func parseConfig(src io.Reader) (config, error) {
	var c config
	bs, err := ioutil.ReadAll(src)
	if err != nil {
		return c, errors.Wrap(err, "cannot slurp from file")
	}
	err = yaml.Unmarshal(bs, &c)
	if err != nil {
		return c, errors.Wrap(err, "cannot parse")
	}

	return c, nil

}
