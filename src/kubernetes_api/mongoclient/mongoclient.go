package mongoclient

import (
	crd "kubernetes_api/mongodeployments"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// This file implement all the (CRUD) client methods we need to access our CRD object
//TODO NEED TO SET THE CRDPLURAL DEPENDING ON THE TYPE OF MONGODB INSTANCE!!!
func CrdClient(cl *rest.RESTClient, scheme *runtime.Scheme, namespace string) *Crdclient {
	return &Crdclient{cl: cl, ns: namespace, plural: crd.CRDPlural,
		codec: runtime.NewParameterCodec(scheme)}
}

type Crdclient struct {
	cl     *rest.RESTClient
	ns     string
	plural string
	codec  runtime.ParameterCodec
}

func (f *Crdclient) Create(obj *crd.Deployment) (*crd.Deployment, error) {
	var result crd.Deployment
	err := f.cl.Post().
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) Update(obj *crd.Deployment) (*crd.Deployment, error) {
	var result crd.Deployment
	err := f.cl.Put().
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) Delete(name string, options *metav1.DeleteOptions) error {
	return f.cl.Delete().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Body(options).Do().
		Error()
}

func (f *Crdclient) Get(name string) (*crd.Deployment, error) {
	var result crd.Deployment
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) List(opts metav1.ListOptions) (*crd.DeploymentList, error) {
	var result crd.DeploymentList
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		VersionedParams(&opts, f.codec).
		Do().Into(&result)
	return &result, err
}

// Create a new List watch for our TPR
func (f *Crdclient) NewListWatch() *cache.ListWatch {
	return cache.NewListWatchFromClient(f.cl, f.plural, f.ns, fields.Everything())
}
