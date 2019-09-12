package v1

import (
	"api/types/v1"

	"k8s.io/client-go/kubernetes"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type MongoDBV1Interface interface {
	MongoDBs(namespace string) MongoDBInterface
	Core(namespace string) CoreInterface
}

type KubeClient struct {
	restClient rest.Interface
	coreV1     *kubernetes.Clientset
}

func NewForConfig(c *rest.Config) (*KubeClient, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1.GroupName, Version: v1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	restClient, err := rest.RESTClientFor(&config)
	if err != nil {
		zap.S().Debugf("Failed to initialise clientset for MongoDB CRD")
		return nil, err
	}
	zap.S().Debugf("Clientset for MongoDB CRD initialised")

	coreV1Client, err := kubernetes.NewForConfig(&config)
	if err != nil {
		zap.S().Debugf("Failed to initialise clientset for Core V1 API")
		return nil, err
	}
	zap.S().Debugf("Clientset for Core V1 initialised")

	return &KubeClient{restClient: restClient, coreV1: coreV1Client}, nil
}

func (c *KubeClient) MongoDBs(namespace string) MongoDBInterface {
	return &mongoDBClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *KubeClient) Core(namespace string) CoreInterface {
	return &coreClient{
		client: c.coreV1,
		ns:     namespace,
	}
}
