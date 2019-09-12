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
}

type coreClient struct {
	client kubernetes.Interface
	ns     string
}

func (c *coreClient) CreateConfigMap(projectName string, orgId string, baseUrl string) (*apiv1.ConfigMap, error) {
	var err error
	var queryResult *apiv1.ConfigMap

	queryResult, err = c.client.CoreV1().ConfigMaps(c.ns).Get(projectName, metav1.GetOptions{})
	if err != nil {
		zap.S().Warn(err.Error())
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

		if err != nil {
			return nil, err
		}

		zap.S().Infof("Created config map %q.\n", result.GetObjectMeta().GetName())

		return result, nil
	}
	zap.S().Infof("Nothing to create, config map already exists")
	return queryResult, nil
}

func (c *coreClient) CreateSecret(secretName string, apiUser string, apiKey string) (*apiv1.Secret, error) {
	var err error
	var queryResult *apiv1.Secret

	queryResult, err = c.client.CoreV1().Secrets(c.ns).Get(secretName, metav1.GetOptions{})
	if err != nil {
		zap.S().Warn(err.Error())
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
		if err != nil {
			zap.S().Debug("Creation of Secret failed!")
			return nil, err
		}
		zap.S().Infof("Created secret %q.\n", result.GetObjectMeta().GetName())
		return result, nil
	}
	zap.S().Infof("Nothing to create, secret already exists")
	return queryResult, nil
}
