package quasar

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	DefaultPort = 3273
)

type config struct {
	Port     int      `yaml:"port"`
	Hostname string   `yaml:"hostname"`
	Daemons  []daemon `yaml:"daemons"`
}

func (c config) Address() string {
	return c.Hostname + ":" + strconv.Itoa(c.Port)
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

	usedNames := map[string]struct{}{}
	for i, d := range c.Daemons {
		n := d.Name
		if n == "" {
			n = d.Type
		}
		if _, used := usedNames[n]; used {
			return config{}, errors.Errorf(
				"conflict daemon name %s in config.",
				n,
			)
		}
		usedNames[n] = struct{}{}
		d.Name = n
		c.Daemons[i] = d
	}

	if c.Port == 0 {
		c.Port = 3273
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
