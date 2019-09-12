package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type MongoSpec struct {
	Credentials string `json:"credentials"`
	Project     string `json:"project"`
	Version     string `json:"version"`
	Type        string `json:"type"`

	Members              int `json:"members"`
	ConfigServerCount    int `json:"configServerCount"`
	MongoDsPerShardCount int `json:"mongodsPerShardCount"`
	MongosCount          int `json:"mongosCount"`
	ShardCount           int `json:"shardCount"`

	LogLevel               string               `json:"logLevel,omitempty"`
	Security               SecuritySpec         `json:"security,omitempty"`
	AdditionalMongoDConfig AdditionalParamsSpec `json:"additionalMongodConfig,omitempty"`
	ExposedExternally      bool                 `json:"exposedExternally,omitempty"`
}

type SecuritySpec struct {
	TLS                       TLSSpec `json:"tls,omitempty"`
	ClusterAuthenticationMode string  `json:"clusterAuthenticationMode,omitempty"`
}

type AdditionalParamsSpec struct {
	Net NetSpec `json:"net,omitempty"`
}

type NetSpec struct {
	SSL SSLSpec `json:"ssl,omitempty"`
}

type SSLSpec struct {
	Mode string `json:"mode,omitempty"`
}

type TLSSpec struct {
	Enabled bool   `json:"enabled,omitempty"`
	CA      string `json:"ca,omitempty"`
}

type MongoDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              MongoSpec `json:"spec"`
}

type MongoDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDB `json:"items"`
}
