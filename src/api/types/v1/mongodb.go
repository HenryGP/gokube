package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type MongoSpec struct {
	Credentials string `json:"credentials"`
	Project     string `json:"project"`
	Version     string `json:"version"`
	Type        string `json:"type"`
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
