package v1

import (
	"api/types/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type MongoDBInterface interface {
	List(opts metav1.ListOptions) (*v1.MongoDBList, error)
	Get(name string, options metav1.GetOptions) (*v1.MongoDB, error)
	Create(*v1.MongoDB) (*v1.MongoDB, error)
	Delete(name string) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type mongoDBClient struct {
	restClient rest.Interface
	ns         string
}

func (c *mongoDBClient) List(opts metav1.ListOptions) (*v1.MongoDBList, error) {
	result := v1.MongoDBList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("mongodb").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *mongoDBClient) Get(name string, opts metav1.GetOptions) (*v1.MongoDB, error) {
	result := v1.MongoDB{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("mongodb").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *mongoDBClient) Create(mongodb *v1.MongoDB) (*v1.MongoDB, error) {
	result := v1.MongoDB{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("mongodb").
		Body(mongodb).
		Do().
		Into(&result)
	return &result, err
}

func (c *mongoDBClient) Delete(name string) error {
	return c.restClient.Delete().
		Namespace(c.ns).
		Resource("mongodb").
		Name(name).
		Do().
		Error()
}

func (c *mongoDBClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("mongodb").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
