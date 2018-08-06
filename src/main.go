package main

import (
	"fmt"
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
	//fmt.Println(config.OpsManager["project"])
	k8s := kubernetes_api.New("enrique-test")
	fmt.Println("TEST")
	//k8s.CreateEnvironment()
	//k8s.DeleteEnvironment()
}
