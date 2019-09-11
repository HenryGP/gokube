package v1

import (
	"fmt"

	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CoreInterface interface {
	CreateConfigMap(projectId string, baseUrl string)
}

type coreClient struct {
	client kubernetes.Interface
	ns     string
}

func (c *coreClient) CreateConfigMap(projectId string, baseUrl string) {
	configMap := &apiv1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dredd-project",
			Namespace: c.ns,
		},
		Data: map[string]string{
			"projectId": projectId,
			"baseUrl":   baseUrl,
		},
	}

	result, err := c.client.CoreV1().ConfigMaps(c.ns).Create(configMap)
	if err != nil {
		zap.S().Warn(err.Error())
	}
	fmt.Printf("Created config map %q.\n", result.GetObjectMeta().GetName())
}
