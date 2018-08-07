//REFERENCE! https://www.martin-helmich.de/en/blog/kubernetes-crd-client.html
package mongodeployments

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

const (
	CRDGroup   string = "mongodb.com"
	CRDVersion string = "v1"
	CRDPlural  string = ""
)

type Deployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              DeploymentSpec `json:"spec"`
}

type DeploymentSpec struct {
	Persistent           bool   `json:"persistent"`
	Version              string `json:"version"`
	Credentials          string `json:"credentials"`
	Project              string `json:"project"`
	Members              string `json:"members,omitempty"`              //Replica Sets
	ShardCount           string `json:"shardCount,omitempty"`           //Sharded Cluster
	MongodsPerShardCount string `json:"mongodsPerShardCount,omitempty"` //Sharded Cluster
	MongosCount          string `json:"mongosCount,omitempty"`          //Sharded Cluster
	ConfigServerCount    string `json:"configServerCount,omitempty"`    //Sharded Cluster
}

type DeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Deployment `json:"items"`
}

func (in *Deployment) DeepCopyInto(out *Deployment) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = DeploymentSpec{
		Persistent:           in.Spec.Persistent,
		Version:              in.Spec.Version,
		Credentials:          in.Spec.Credentials,
		Project:              in.Spec.Project,
		Members:              in.Spec.Members,
		ShardCount:           in.Spec.ShardCount,
		MongodsPerShardCount: in.Spec.MongodsPerShardCount,
		MongosCount:          in.Spec.MongosCount,
		ConfigServerCount:    in.Spec.ConfigServerCount,
	}
}

func (in *Deployment) DeepCopyObject() runtime.Object {
	out := Deployment{}
	in.DeepCopyInto(&out)
	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *DeploymentList) DeepCopyObject() runtime.Object {
	out := DeploymentList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]Deployment, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return &out
}

// Create a  Rest client with the new CRD Schema
var SchemeGroupVersion = schema.GroupVersion{Group: CRDGroup, Version: CRDVersion}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Deployment{},
		&DeploymentList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func NewClient(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	config := *cfg
	config.GroupVersion = &SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, nil, err
	}
	return client, scheme, nil
}
