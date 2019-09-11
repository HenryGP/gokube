package v1

import (
	"api/types/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type MongoDBV1Interface interface {
	MongoDBs(namespace string) MongoDBInterface
}

type MongoDBV1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*MongoDBV1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1.GroupName, Version: v1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &MongoDBV1Client{restClient: client}, nil
}

func (c *MongoDBV1Client) MongoDBs(namespace string) MongoDBInterface {
	return &mongoDBClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
