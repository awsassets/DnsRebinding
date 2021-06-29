package conf

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Domain struct {
		Main     string   `yaml:"main"`
		NS       []string `yaml:"ns"`
		IP       string   `yaml:"ip"`
		RebindIP string   `yaml:"rebind"`
	} `yaml:"domain"`
}

var (
	C config
)

// SetFromFile 从文件设置 config
func SetFromFile(c string) error {
	var (
		f   *os.File
		buf []byte
		err error
	)
	if f, err = os.Open(c); err != nil {
		return err
	}
	if buf, err = ioutil.ReadAll(f); err != nil {
		return err
	}

	if err = yaml.Unmarshal(buf, &C); err != nil {
		return err
	}
	return nil
}
