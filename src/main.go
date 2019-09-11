package main

import (
	"api/types/v1"
	clientV1 "clientset/v1"
	"flag"
	"io/ioutil"
	"path/filepath"

	logging "logging"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

type AppConf struct {
	OpsManager map[string]string `yaml:"ops_manager"`
	Kubernetes map[string]string `yaml:"kubernetes"`
	Logger     string            `yaml:"logger"`
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func (c *AppConf) getConf() (*AppConf, error) {
	filename, _ := filepath.Abs("config.yaml")
	confFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(confFile, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func main() {

	var err error
	var appConfig AppConf
	appConfig.getConf()

	logger, err := logging.InitLogger(appConfig.Logger)
	if err != nil {
		logger.Panic(err.Error())
	}
	logger.Infof("ZAP logger initialised in %s mode", appConfig.Logger)

	var config *rest.Config

	logger.Info("Loading Kubernetes cluster configuration")
	if kubeconfig == "" {
		logger.Info("Using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		logger.Infof("Using configuration from %s", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		logger.Panic(err.Error())
	}

	v1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1.NewForConfig(config)
	if err != nil {
		logger.Panic(err.Error())
	}

	clientSet.Core(appConfig.Kubernetes["namespace"]).
		CreateConfigMap(appConfig.Kubernetes["project"], appConfig.Kubernetes["base_url"])

	/*
		mongodbs, err := clientSet.MongoDBs(appConfig.Kubernetes["namespace"]).List(metav1.ListOptions{})
		if err != nil {
			logger.Panic(err.Error())
		}

		logger.Infof("mongodbs found: %+v\n", mongodbs)*/
}
