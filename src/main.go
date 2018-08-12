package main

import (
	"fmt"
	"io/ioutil"
	"kubernetes_api"
	"path/filepath"

	"go.uber.org/zap"

	yaml "gopkg.in/yaml.v2"
)

var log *zap.SugaredLogger

type conf struct {
	OpsManager map[string]string `yaml:"ops_manager"`
	Kubernetes map[string]string `yaml:"kubernetes"`
	Logger     string            `yaml:"logger"`
}

func (c *conf) getConf() *conf {
	filename, _ := filepath.Abs("config.yaml")
	confFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("Error when opening config file: #%v ", err)
	}
	err = yaml.Unmarshal(confFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func initLogger(preset string) {
	var logger *zap.Logger
	var err error

	switch preset {
	case "PROD":
		logger, err = zap.NewProduction()
	case "DEV":
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		fmt.Println("Failed to create logger, will use the default one")
		fmt.Println(err)
	}

	zap.ReplaceGlobals(logger)
	log = zap.S()
}

func main() {
	var config conf
	config.getConf()

	initLogger(config.Logger)

	k8s := kubernetes_api.New(config.Kubernetes["namespace"])

	k8s.CreateEnvironment(config.OpsManager["project"], config.OpsManager["api_user"], config.OpsManager["api_password"], config.OpsManager["base_url"])

	//k8s.CreateStandalone("stand", "3.4.10")
	//k8s.CreateReplicaSet("rs", "4.0.0", 3)
	//k8s.CreateShardedCluster("shcluster", "3.6.2", 1, 3, 3, 2)

	//k8s.DeleteStandalone("stand")
	//k8s.DeleteReplicaSet("rs")
	//k8s.DeleteShardedCluster("shcluster")
	//k8s.DeleteEnvironment()
}
