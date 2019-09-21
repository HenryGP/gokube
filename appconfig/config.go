package appconfig

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type AppConf struct {
	Kubernetes map[string]string `yaml:"kubernetes"`
	Logger     string            `yaml:"logger"`
}

func GetConf() (*AppConf, error) {
	var c AppConf
	filename, _ := filepath.Abs("config.yml")
	confFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(confFile, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
