package mongoclient

import (
	crd "kubernetes_api/mongodeployments"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

// This file implement all the (CRUD) client methods we need to access our CRD object
func CrdClient(cl *rest.RESTClient, scheme *runtime.Scheme, namespace string) *Crdclient {
	return &Crdclient{cl: cl, ns: namespace, codec: runtime.NewParameterCodec(scheme)}
}

type Crdclient struct {
	cl    *rest.RESTClient
	ns    string
	codec runtime.ParameterCodec
}

func (f *Crdclient) CreateMongoDbStandalone(obj *crd.MongoDbStandalone) (*crd.MongoDbStandalone, error) {
	var result crd.MongoDbStandalone
	err := f.cl.Post().
		Namespace(f.ns).Resource("mongodbstandalones").
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) UpdateMongoDbStandalone(obj *crd.MongoDbStandalone) (*crd.MongoDbStandalone, error) {
	var result crd.MongoDbStandalone
	err := f.cl.Put().
		Namespace(f.ns).Resource("mongodbstandalones").
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) DeleteMongoDbStandalone(name string, options *metav1.DeleteOptions) error {
	return f.cl.Delete().
		Namespace(f.ns).Resource("mongodbstandalones").
		Name(name).Body(options).Do().
		Error()
}

func (f *Crdclient) GetMongoDbStandalone(name string) (*crd.MongoDbStandalone, error) {
	var result crd.MongoDbStandalone
	err := f.cl.Get().
		Namespace(f.ns).Resource("mongodbstandalones").
		Name(name).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) ListMongoDbStandalone(opts metav1.ListOptions) (*crd.MongoDbStandalone, error) {
	var result crd.MongoDbStandalone
	err := f.cl.Get().
		Namespace(f.ns).Resource("mongodbstandalones").
		VersionedParams(&opts, f.codec).
		Do().Into(&result)
	return &result, err
}

//Replica Sets

func (f *Crdclient) CreateMongoDbReplicaSet(obj *crd.MongoDbReplicaSet) (*crd.MongoDbReplicaSet, error) {
	var result crd.MongoDbReplicaSet
	err := f.cl.Post().
		Namespace(f.ns).Resource("mongodbreplicasets").
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) UpdateMongoDbReplicaSet(obj *crd.MongoDbReplicaSet) (*crd.MongoDbReplicaSet, error) {
	var result crd.MongoDbReplicaSet
	err := f.cl.Put().
		Namespace(f.ns).Resource("mongodbreplicasets").
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) DeleteMongoDbReplicaSet(name string, options *metav1.DeleteOptions) error {
	return f.cl.Delete().
		Namespace(f.ns).Resource("mongodbreplicasets").
		Name(name).Body(options).Do().
		Error()
}

func (f *Crdclient) GetMongoDbReplicaSet(name string) (*crd.MongoDbReplicaSet, error) {
	var result crd.MongoDbReplicaSet
	err := f.cl.Get().
		Namespace(f.ns).Resource("mongodbreplicasets").
		Name(name).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) ListMongoDbReplicaSet(opts metav1.ListOptions) (*crd.MongoDbReplicaSet, error) {
	var result crd.MongoDbReplicaSet
	err := f.cl.Get().
		Namespace(f.ns).Resource("mongodbreplicasets").
		VersionedParams(&opts, f.codec).
		Do().Into(&result)
	return &result, err
}

//Sharded clusters

func (f *Crdclient) CreateMongoDbShardedCluster(obj *crd.MongoDbShardedCluster) (*crd.MongoDbShardedCluster, error) {
	var result crd.MongoDbShardedCluster
	err := f.cl.Post().
		Namespace(f.ns).Resource("mongodbshardedclusters").
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) UpdateMongoDbShardedCluster(obj *crd.MongoDbShardedCluster) (*crd.MongoDbShardedCluster, error) {
	var result crd.MongoDbShardedCluster
	err := f.cl.Put().
		Namespace(f.ns).Resource("mongodbshardedclusters").
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) DeleteMongoDbShardedCluster(name string, options *metav1.DeleteOptions) error {
	return f.cl.Delete().
		Namespace(f.ns).Resource("mongodbshardedclusters").
		Name(name).Body(options).Do().
		Error()
}

func (f *Crdclient) GetMongoDbShardedCluster(name string) (*crd.MongoDbShardedCluster, error) {
	var result crd.MongoDbShardedCluster
	err := f.cl.Get().
		Namespace(f.ns).Resource("mongodbshardedclusters").
		Name(name).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) ListMongoDbShardedCluster(opts metav1.ListOptions) (*crd.MongoDbShardedCluster, error) {
	var result crd.MongoDbShardedCluster
	err := f.cl.Get().
		Namespace(f.ns).Resource("mongodbshardedclusters").
		VersionedParams(&opts, f.codec).
		Do().Into(&result)
	return &result, err
}

// Create a new List watch for our TPR
//func (f *Crdclient) NewListWatch() *cache.ListWatch {
//	return cache.NewListWatchFromClient(f.cl, f.plural, f.ns, fields.Everything())
//}
