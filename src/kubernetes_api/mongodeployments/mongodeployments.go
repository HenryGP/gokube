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
)

type MongoDbStandalone struct {
	metav1.TypeMeta       `json:",inline"`
	metav1.ObjectMeta     `json:"metadata"`
	MongoDbStandaloneSpec `json:"spec"`
}

type MongoDbStandaloneSpec struct {
	Persistent  bool   `json:"persistent"`
	Version     string `json:"version"`
	Credentials string `json:"credentials"`
	Project     string `json:"project"`
}

type MongoDbStandaloneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []MongoDbStandalone `json:"items"`
}

type MongoDbReplicaSet struct {
	metav1.TypeMeta       `json:",inline"`
	metav1.ObjectMeta     `json:"metadata"`
	MongoDbReplicaSetSpec `json:"spec"`
}

type MongoDbReplicaSetSpec struct {
	Persistent  bool   `json:"persistent"`
	Version     string `json:"version"`
	Credentials string `json:"credentials"`
	Project     string `json:"project"`
	Members     int    `json:"members"`
}

type MongoDbReplicaSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []MongoDbReplicaSet `json:"items"`
}

//Cluster

type MongoDbShardedCluster struct {
	metav1.TypeMeta           `json:",inline"`
	metav1.ObjectMeta         `json:"metadata"`
	MongoDbShardedClusterSpec `json:"spec"`
}

type MongoDbShardedClusterSpec struct {
	Persistent           bool   `json:"persistent"`
	Version              string `json:"version"`
	Credentials          string `json:"credentials"`
	Project              string `json:"project"`
	ConfigServerCount    int    `json:"configServerCount"`
	ShardCount           int    `json:"shardCount"`
	MongosCount          int    `json:"mongosCount"`
	MongodsPerShardCount int    `json:"mongodsPerShardCount"`
}

type MongoDbShardedClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []MongoDbShardedCluster `json:"items"`
}

func (in *MongoDbShardedCluster) DeepCopyInto(out *MongoDbShardedCluster) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.MongoDbShardedClusterSpec = MongoDbShardedClusterSpec{
		Persistent:           in.MongoDbShardedClusterSpec.Persistent,
		Version:              in.MongoDbShardedClusterSpec.Version,
		Credentials:          in.MongoDbShardedClusterSpec.Credentials,
		Project:              in.MongoDbShardedClusterSpec.Project,
		ConfigServerCount:    in.MongoDbShardedClusterSpec.ConfigServerCount,
		ShardCount:           in.MongoDbShardedClusterSpec.ShardCount,
		MongosCount:          in.MongoDbShardedClusterSpec.MongosCount,
		MongodsPerShardCount: in.MongoDbShardedClusterSpec.MongodsPerShardCount,
	}
}

func (in *MongoDbShardedCluster) DeepCopyObject() runtime.Object {
	out := MongoDbShardedCluster{}
	in.DeepCopyInto(&out)
	return &out
}

func (in *MongoDbShardedClusterList) DeepCopyObject() runtime.Object {
	out := MongoDbShardedClusterList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]MongoDbShardedCluster, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return &out
}

//

func (in *MongoDbReplicaSet) DeepCopyInto(out *MongoDbReplicaSet) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.MongoDbReplicaSetSpec = MongoDbReplicaSetSpec{
		Persistent:  in.MongoDbReplicaSetSpec.Persistent,
		Version:     in.MongoDbReplicaSetSpec.Version,
		Credentials: in.MongoDbReplicaSetSpec.Credentials,
		Project:     in.MongoDbReplicaSetSpec.Project,
		Members:     in.MongoDbReplicaSetSpec.Members,
	}
}

func (in *MongoDbReplicaSet) DeepCopyObject() runtime.Object {
	out := MongoDbReplicaSet{}
	in.DeepCopyInto(&out)
	return &out
}

func (in *MongoDbReplicaSetList) DeepCopyObject() runtime.Object {
	out := MongoDbReplicaSetList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]MongoDbReplicaSet, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return &out
}

func (in *MongoDbStandalone) DeepCopyInto(out *MongoDbStandalone) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.MongoDbStandaloneSpec = MongoDbStandaloneSpec{
		Persistent:  in.MongoDbStandaloneSpec.Persistent,
		Version:     in.MongoDbStandaloneSpec.Version,
		Credentials: in.MongoDbStandaloneSpec.Credentials,
		Project:     in.MongoDbStandaloneSpec.Project,
	}
}

func (in *MongoDbStandalone) DeepCopyObject() runtime.Object {
	out := MongoDbStandalone{}
	in.DeepCopyInto(&out)
	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *MongoDbStandaloneList) DeepCopyObject() runtime.Object {
	out := MongoDbStandaloneList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]MongoDbStandalone, len(in.Items))
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
		&MongoDbStandalone{},
		&MongoDbStandaloneList{},
		&MongoDbReplicaSet{},
		&MongoDbReplicaSetList{},
		&MongoDbShardedCluster{},
		&MongoDbShardedClusterList{},
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
