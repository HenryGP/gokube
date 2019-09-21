package main

import (
	clientv1 "github.com/10gen/dredd/clientset/v1"
	typesv1 "github.com/10gen/dredd/crdapi/types/v1"

	logging "github.com/10gen/dredd/logging"

	"github.com/10gen/dredd/appconfig"
	"github.com/10gen/dredd/webapi/v1"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func InitialiseKubernetesConfig(kubeconfig string) (*rest.Config, error) {
	var err error
	var config *rest.Config
	if kubeconfig == "" {
		zap.S().Info("Using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		zap.S().Infof("Using configuration from %s", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return config, err
}

func main() {
	var err error
	var appConfig *appconfig.AppConf
	appConfig, err = appconfig.GetConf()
	if err != nil {
		panic(err.Error)
	}

	logger, err := logging.InitLogger(appConfig.Logger)
	if err != nil {
		logger.Panic(err.Error())
	}
	logger.Infof("ZAP logger initialised in %s mode", appConfig.Logger)

	zap.S().Info("Loading Kubernetes cluster configuration")
	config, err := InitialiseKubernetesConfig(appConfig.Kubernetes["kubeconfig"])
	if err != nil {
		zap.S().Panic(err.Error())
	}

	zap.S().Info("Initialising KubeClient")
	clientSet, err := clientv1.NewForConfig(config)
	if err != nil {
		logger.Panic(err.Error())
	}

	logger.Infof("Addig schemes")
	typesv1.AddToScheme(scheme.Scheme)

	// Web Server
	zap.S().Info("Initialising API server on port 8080")
	err = webapi.ServeAPI(":8080", clientSet, appConfig.Kubernetes["namespace"])
	if err != nil {
		logger.Panic(err.Error())
	}
}
