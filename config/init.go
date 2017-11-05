package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var Path string

func Init() (err error) {
	var barr []byte
	if barr, err = ioutil.ReadFile(Path); err != nil {
		return
	}
	if err = yaml.Unmarshal(barr, &Config); err != nil {
		return
	}
	return
}
