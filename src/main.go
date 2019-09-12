package main

import (
	"api/types/v1"
	clientV1 "clientset/v1"
	"flag"
	"io/ioutil"
	"path/filepath"

	logging "logging"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		CreateConfigMap(appConfig.Kubernetes["config_map_name"], appConfig.OpsManager["org_id"], appConfig.OpsManager["base_url"])

	clientSet.Core(appConfig.Kubernetes["namespace"]).
		CreateSecret(appConfig.Kubernetes["secret_name"], appConfig.OpsManager["api_user"], appConfig.OpsManager["api_password"])

	/*
		var shardedCluster = v1.MongoDB{
			TypeMeta: metav1.TypeMeta{
				Kind:       "MongoDB",
				APIVersion: "mongodb.com/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "shardedclustertest",
				Namespace: appConfig.Kubernetes["namespace"],
			},
			Spec: v1.MongoSpec{
				Version:              "4.0.4",
				Credentials:          appConfig.Kubernetes["secret_name"],
				Project:              appConfig.Kubernetes["config_map_name"],
				Type:                 "ShardedCluster",
				ConfigServerCount:    3,
				MongoDsPerShardCount: 3,
				MongosCount:          3,
				ShardCount:           2,
			},
		}

		var result *v1.MongoDB
		result, err = clientSet.MongoDBs(appConfig.Kubernetes["namespace"]).
			Create(&shardedCluster)

		if err != nil {
			logger.Panic(err.Error())
		}
		logger.Infof("mongodb %d", result)


			var replicaSet = v1.MongoDB{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MongoDB",
					APIVersion: "mongodb.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "replicasettest",
					Namespace: appConfig.Kubernetes["namespace"],
				},
				Spec: v1.MongoSpec{
					Version:     "4.0.4",
					Credentials: appConfig.Kubernetes["secret_name"],
					Project:     appConfig.Kubernetes["config_map_name"],
					Type:        "ReplicaSet",
					Members:     3,
				},
			}

			var result *v1.MongoDB
			result, err = clientSet.MongoDBs(appConfig.Kubernetes["namespace"]).
				Create(&replicaSet)

			if err != nil {
				logger.Panic(err.Error())
			}
			logger.Infof("mongodb %d", result)


				var standalone = v1.MongoDB{
					TypeMeta: metav1.TypeMeta{
						Kind:       "MongoDB",
						APIVersion: "mongodb.com/v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "standalonetest",
						Namespace: appConfig.Kubernetes["namespace"],
					},
					Spec: v1.MongoSpec{
						Version:     "4.0.4",
						Credentials: appConfig.Kubernetes["secret_name"],
						Project:     appConfig.Kubernetes["config_map_name"],
						Type:        "Standalone",
					},
				}

				var result *v1.MongoDB
				result, err = clientSet.MongoDBs(appConfig.Kubernetes["namespace"]).Create(&standalone)

				if err != nil {
					logger.Panic(err.Error())
				}
				logger.Infof("mongodb %d", result)

		mongodbs, err := clientSet.MongoDBs(appConfig.Kubernetes["namespace"]).List(metav1.ListOptions{})
		if err != nil {
			logger.Panic(err.Error())
		}

		logger.Infof("mongodbs found: %+v\n", mongodbs)*/
}
