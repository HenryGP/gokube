package v1

import (
	"encoding/base64"

	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CoreInterface interface {
	CreateConfigMap(projectName string, projectId string, baseUrl string) (*apiv1.ConfigMap, error)
	CreateSecret(projectName string, apiUser string, apiKey string) (*apiv1.Secret, error)
	DeleteConfigMap(projectName string) error
	DeleteSecret(secretName string) error
	GetConfigMaps() (*apiv1.ConfigMapList, error)
	GetConfigMap(projectName string) (*apiv1.ConfigMap, error)
	GetSecret(secretName string) (*apiv1.Secret, error)
	GetSecrets() (*apiv1.SecretList, error)
}

type coreClient struct {
	client kubernetes.Interface
	ns     string
}

func (c *coreClient) DeleteConfigMap(projectName string) error {
	err := c.client.CoreV1().ConfigMaps(c.ns).Delete(projectName, &metav1.DeleteOptions{})
	return err
}

func (c *coreClient) GetConfigMap(projectName string) (*apiv1.ConfigMap, error) {
	cfgMap, err := c.client.CoreV1().ConfigMaps(c.ns).Get(projectName, metav1.GetOptions{})
	return cfgMap, err
}

func (c *coreClient) GetConfigMaps() (*apiv1.ConfigMapList, error) {
	cfgMaps, err := c.client.CoreV1().ConfigMaps(c.ns).List(metav1.ListOptions{})
	return cfgMaps, err
}

func (c *coreClient) CreateConfigMap(projectName string, orgId string, baseUrl string) (*apiv1.ConfigMap, error) {
	configMap := &apiv1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      projectName,
			Namespace: c.ns,
		},
		Data: map[string]string{
			"projectName": projectName,
			"orgId":       orgId,
			"baseUrl":     baseUrl,
		},
	}
	result, err := c.client.CoreV1().ConfigMaps(c.ns).Create(configMap)
	return result, err
}

func (c *coreClient) DeleteSecret(secretName string) error {
	err := c.client.CoreV1().Secrets(c.ns).Delete(secretName, &metav1.DeleteOptions{})
	return err
}

func (c *coreClient) CreateSecret(secretName string, apiUser string, apiKey string) (*apiv1.Secret, error) {
	encodedUser := base64.StdEncoding.EncodeToString([]byte(apiUser))
	encodedKey := base64.StdEncoding.EncodeToString([]byte(apiKey))
	decodedUser, err := base64.StdEncoding.DecodeString(encodedUser)
	if err != nil {
		zap.S().Debug("Decoding String for user failed!")
		return nil, err
	}
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		zap.S().Debug("Decoding String for key failed!")
		return nil, err
	}
	secret := &apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: c.ns,
		},
		Type: "from-literal",
		Data: map[string][]byte{
			"user":         decodedUser,
			"publicApiKey": decodedKey,
		},
	}
	result, err := c.client.CoreV1().Secrets(c.ns).Create(secret)
	return result, err
}

func (c *coreClient) GetSecret(secretName string) (*apiv1.Secret, error) {
	secret, err := c.client.CoreV1().Secrets(c.ns).Get(secretName, metav1.GetOptions{})
	return secret, err
}

func (c *coreClient) GetSecrets() (*apiv1.SecretList, error) {
	secrets, err := c.client.CoreV1().Secrets(c.ns).List(metav1.ListOptions{})
	return secrets, err
}
