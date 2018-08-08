package main

import (
	"io/ioutil"
	"kubernetes_api"
	"log"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	OpsManager map[string]string `yaml:"ops_manager"`
}

func (c *conf) getConf() *conf {
	filename, _ := filepath.Abs("config.yaml")
	confFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error when opening config file: #%v ", err)
	}
	err = yaml.Unmarshal(confFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func main() {
	var config conf
	config.getConf()
	k8s := kubernetes_api.New("enrique-test")
	k8s.CreateEnvironment(config.OpsManager["project"], config.OpsManager["api_user"],
		config.OpsManager["api_password"], config.OpsManager["base_url"])

	k8s.CreateStandalone("stand", "3.4.10")

	//k8s.DeleteEnvironment()
}
