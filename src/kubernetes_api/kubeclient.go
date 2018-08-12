package kubernetes_api

import (
	"encoding/base64"
	"flag"
	"fmt"
	"kubernetes_api/mongoclient"
	"log"
	"os"
	"path/filepath"

	mongodeployments "kubernetes_api/mongodeployments"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type kubeClient struct {
	corev1      *kubernetes.Clientset
	namespace   string
	deployments *mongoclient.Crdclient
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func New(nsName string) kubeClient {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		zap.S().Panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		zap.S().Panic(err.Error())
	}

	crdcs, scheme, err := mongodeployments.NewClient(config)
	if err != nil {
		zap.S().Panic(err.Error())
	}

	deploymentsClient := mongoclient.CrdClient(crdcs, scheme, nsName)

	return kubeClient{corev1: clientset, deployments: deploymentsClient, namespace: nsName}
}

func (client kubeClient) CreateEnvironment(project string, apiUser string, apiKey string, baseURL string) {
	var err error
	getOpts := metav1.GetOptions{}

	_, err = client.corev1.CoreV1().Namespaces().Get(client.namespace, getOpts)
	if err != nil {
		zap.S().Warn(err.Error())
		client.createNamespace(client.namespace)
	}

	_, err = client.corev1.RbacV1().ClusterRoles().Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		zap.S().Warn(err.Error())
		client.createClusterRole()
	}

	_, err = client.corev1.CoreV1().ServiceAccounts(client.namespace).Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		zap.S().Warn(err.Error())
		client.createServiceAccount()
	}

	_, err = client.corev1.RbacV1().ClusterRoleBindings().Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		zap.S().Warn(err.Error())
		client.createClusterRoleBinding()
	}

	_, err = client.corev1.CoreV1().ConfigMaps(client.namespace).Get("dredd-project", getOpts)
	if err != nil {
		zap.S().Warn(err.Error())
		client.createConfigMap(project, baseURL)
	}

	_, err = client.corev1.CoreV1().Secrets(client.namespace).Get("dredd-om-credentials", getOpts)
	if err != nil {
		zap.S().Warn(err.Error())
		client.createSecret(apiUser, apiKey)
	}

	_, err = client.corev1.AppsV1().Deployments(client.namespace).Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		zap.S().Warn(err.Error())
		client.createOperator()
	}

}

func (client kubeClient) DeleteEnvironment() {
	var err error
	getOpts := metav1.GetOptions{}

	_, err = client.corev1.CoreV1().Secrets(client.namespace).Get("dredd-om-credentials", getOpts)
	if err != nil {
		zap.S().Error(err.Error())
	} else {
		client.deleteSecret()
	}

	_, err = client.corev1.CoreV1().ConfigMaps(client.namespace).Get("dredd-project", getOpts)
	if err != nil {
		zap.S().Error(err.Error())
	} else {
		client.deleteConfigMap()
	}

	_, err = client.corev1.AppsV1().Deployments(client.namespace).Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		zap.S().Error(err.Error())
	} else {
		client.deleteOperator()
	}

	_, err = client.corev1.RbacV1().ClusterRoleBindings().Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		zap.S().Error(err.Error())
	} else {
		client.deleteClusterRoleBinding()
	}

	_, err = client.corev1.CoreV1().ServiceAccounts(client.namespace).Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		zap.S().Error(err.Error())
	} else {
		client.deleteServiceAccount()
	}

	_, err = client.corev1.RbacV1().ClusterRoles().Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		zap.S().Error(err.Error())
	} else {
		client.deleteClusterRole()
	}

	_, err = client.corev1.CoreV1().Namespaces().Get(client.namespace, getOpts)
	if err != nil {
		zap.S().Error(err.Error())
	} else {
		client.deleteNamespace(client.namespace)
	}

}

func (client kubeClient) createNamespace(ns string) {
	nsSpec := &apiv1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}
	result, err := client.corev1.CoreV1().Namespaces().Create(nsSpec)
	if err != nil {
		zap.S().Panic(err.Error())
	}
	fmt.Printf("Created namespace %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteNamespace(ns string) {
	deleteOptions := metav1.DeleteOptions{}
	err := client.corev1.CoreV1().Namespaces().Delete(ns, &deleteOptions)
	if err != nil {
		zap.S().Warn(err.Error())
	}
}

func (client kubeClient) createClusterRole() {
	rules := make([]rbacv1.PolicyRule, 0, 4)

	rules = append(rules, rbacv1.PolicyRule{Verbs: []string{"get", "list", "create", "update", "delete"},
		APIGroups: []string{""},
		Resources: []string{"configmaps", "secrets", "services"}})

	rules = append(rules, rbacv1.PolicyRule{Verbs: []string{"*"},
		APIGroups: []string{"apps"},
		Resources: []string{"statefulsets"}})

	rules = append(rules, rbacv1.PolicyRule{Verbs: []string{"get", "list", "watch", "create", "delete"},
		APIGroups: []string{"apiextensions.k8s.io"},
		Resources: []string{"customresourcedefinitions"}})

	rules = append(rules, rbacv1.PolicyRule{Verbs: []string{"*"},
		APIGroups: []string{"mongodb.com"},
		Resources: []string{"*"}})

	clusterRoleSpec := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "mongodb-enterprise-operator",
		},
		Rules: rules,
	}

	result, err := client.corev1.RbacV1().ClusterRoles().Create(clusterRoleSpec)
	if err != nil {
		zap.S().Panic(err.Error())
	}
	fmt.Printf("Created cluster role %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteClusterRole() {
	deleteOptions := metav1.DeleteOptions{}
	err := client.corev1.RbacV1().ClusterRoles().Delete("mongodb-enterprise-operator", &deleteOptions)
	if err != nil {
		zap.S().Warn(err.Error())
	}
}

func (client kubeClient) createServiceAccount() {
	serviceAccount := apiv1.ServiceAccount{
		TypeMeta:   metav1.TypeMeta{Kind: "ServiceAccount", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "mongodb-enterprise-operator", Namespace: client.namespace},
	}

	result, err := client.corev1.CoreV1().ServiceAccounts(client.namespace).Create(&serviceAccount)
	if err != nil {
		zap.S().Panic(err.Error())
	}
	fmt.Printf("Created service account %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteServiceAccount() {
	deleteOptions := metav1.DeleteOptions{}
	err := client.corev1.CoreV1().ServiceAccounts(client.namespace).Delete("mongodb-enterprise-operator", &deleteOptions)
	if err != nil {
		zap.S().Warn(err.Error())
	}
}

func (client kubeClient) createClusterRoleBinding() {

	clusterRoleBindingSpec := rbacv1.ClusterRoleBinding{
		TypeMeta:   metav1.TypeMeta{Kind: "ClusterRoleBinding", APIVersion: "rbac.authorization.k8s.io/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "mongodb-enterprise-operator", Namespace: client.namespace},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Namespace: client.namespace,
				Name:      "mongodb-enterprise-operator",
				Kind:      "ServiceAccount",
			},
		},
		RoleRef: rbacv1.RoleRef{APIGroup: "rbac.authorization.k8s.io", Kind: "ClusterRole", Name: "mongodb-enterprise-operator"},
	}

	result, err := client.corev1.RbacV1().ClusterRoleBindings().Create(&clusterRoleBindingSpec)
	if err != nil {
		zap.S().Panic(err.Error())
	}
	fmt.Printf("Created cluster role binding %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteClusterRoleBinding() {
	deleteOptions := metav1.DeleteOptions{}
	err := client.corev1.RbacV1().ClusterRoleBindings().Delete("mongodb-enterprise-operator", &deleteOptions)
	if err != nil {
		zap.S().Warn(err.Error())
	}
}

func (client kubeClient) createConfigMap(projectId string, baseUrl string) {
	configMap := &apiv1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dredd-project",
			Namespace: client.namespace,
		},
		Data: map[string]string{
			"projectId": projectId,
			"baseUrl":   baseUrl,
		},
	}

	result, err := client.corev1.CoreV1().ConfigMaps(client.namespace).Create(configMap)
	if err != nil {
		zap.S().Warn(err.Error())
	}
	fmt.Printf("Created config map %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteConfigMap() {
	deleteOptions := metav1.DeleteOptions{}
	err := client.corev1.Core().ConfigMaps(client.namespace).Delete("dredd-project", &deleteOptions)
	if err != nil {
		zap.S().Warn(err.Error())
	}
}

func (client kubeClient) createSecret(username string, key string) {

	encodedUser := base64.StdEncoding.EncodeToString([]byte(username))
	encodedKey := base64.StdEncoding.EncodeToString([]byte(key))

	decodedUser, err := base64.StdEncoding.DecodeString(encodedUser)
	if err != nil {
		log.Output(0, err.Error())
		return
	}

	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		log.Output(0, err.Error())
		return
	}

	secret := &apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dredd-om-credentials",
			Namespace: client.namespace,
		},
		Type: "from-literal",
		Data: map[string][]byte{
			"user":         decodedUser,
			"publicApiKey": decodedKey,
		},
	}

	result, err := client.corev1.CoreV1().Secrets(client.namespace).Create(secret)
	if err != nil {
		zap.S().Panic(err.Error())
	}
	fmt.Printf("Created secret %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteSecret() {
	deleteOptions := metav1.DeleteOptions{}
	err := client.corev1.Core().Secrets(client.namespace).Delete("dredd-om-credentials", &deleteOptions)
	if err != nil {
		zap.S().Warn(err.Error())
	}
}

func (client kubeClient) ListNamespaces() {
	nsList, err := client.corev1.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		log.Output(0, err.Error())
	}
	for i := range nsList.Items {
		fmt.Println(nsList.Items[i].Name)
	}
}

func int32Ptr(i int32) *int32 { return &i }

func (client kubeClient) createOperator() {

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mongodb-enterprise-operator",
			Namespace: client.namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "mongodb-enterprise-operator",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "mongodb-enterprise-operator",
					},
				},
				Spec: apiv1.PodSpec{
					ServiceAccountName: "mongodb-enterprise-operator",
					Containers: []apiv1.Container{
						apiv1.Container{
							Name:            "mongodb-enterprise-operator",
							Image:           "quay.io/mongodb/mongodb-enterprise-operator:0.2",
							ImagePullPolicy: "Always",
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{
									Name:  "OPERATOR_ENV",
									Value: "local",
								},
								apiv1.EnvVar{
									Name:  "MONGODB_ENTERPRISE_DATABASE_IMAGE",
									Value: "quay.io/mongodb/mongodb-enterprise-database:0.2",
								},
								apiv1.EnvVar{
									Name:  "IMAGE_PULL_POLICY",
									Value: "Always",
								},
								apiv1.EnvVar{
									Name:  "IMAGE_PULL_SECRETS",
									Value: "",
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := client.corev1.AppsV1().Deployments(client.namespace).Create(deployment)
	if err != nil {
		zap.S().Panic(err.Error())
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteOperator() {
	deleteOptions := &metav1.DeleteOptions{}
	err := client.corev1.AppsV1().Deployments(client.namespace).Delete("mongodb-enterprise-operator", deleteOptions)
	if err != nil {
		zap.S().Error(err.Error())
	}
}

func (client kubeClient) CreateStandalone(name string, mongoVersion string) {
	deployment := &mongodeployments.MongoDbStandalone{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MongoDbStandalone",
			APIVersion: "mongodb.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: client.namespace,
		},
		MongoDbStandaloneSpec: mongodeployments.MongoDbStandaloneSpec{
			Persistent:  false,
			Version:     mongoVersion,
			Credentials: "dredd-om-credentials",
			Project:     "dredd-project",
		},
	}

	result, err := client.deployments.CreateMongoDbStandalone(deployment)
	if err != nil {
		zap.S().Error(err.Error())
	} else if err == nil {
		zap.S().Infof("Created: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		zap.S().Warnf("Already exists: %#v\n", result)
	} else {
		zap.S().Panic(err.Error())
	}

}

func (client kubeClient) DeleteStandalone(name string) {
	deleteOptions := &metav1.DeleteOptions{}
	err := client.deployments.DeleteMongoDbStandalone(name, deleteOptions)
	if err != nil {
		zap.S().Error(err.Error())
	} else if err == nil {
		zap.S().Infof("Deleted: %#v\n", name)
	} else if apierrors.IsGone(err) {
		zap.S().Warnf("Doesn't exists: %#v\n", name)
	} else {
		zap.S().Panic(err.Error())
	}
}

func (client kubeClient) CreateReplicaSet(name string, mongoVersion string, totalMembers int) {
	deployment := &mongodeployments.MongoDbReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MongoDbReplicaSet",
			APIVersion: "mongodb.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: client.namespace,
		},
		MongoDbReplicaSetSpec: mongodeployments.MongoDbReplicaSetSpec{
			Persistent:  false,
			Version:     mongoVersion,
			Credentials: "dredd-om-credentials",
			Project:     "dredd-project",
			Members:     totalMembers,
		},
	}

	result, err := client.deployments.CreateMongoDbReplicaSet(deployment)
	if err != nil {
		zap.S().Error(err.Error())
	} else if err == nil {
		zap.S().Infof("Created: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		zap.S().Warnf("Already exists: %#v\n", result)
	} else {
		zap.S().Panic(err.Error())
	}
}
func (client kubeClient) DeleteReplicaSet(name string) {
	deleteOptions := &metav1.DeleteOptions{}
	err := client.deployments.DeleteMongoDbReplicaSet(name, deleteOptions)
	if err != nil {
		zap.S().Error(err.Error())
	} else if err == nil {
		zap.S().Infof("Deleted: %#v\n", name)
	} else if apierrors.IsGone(err) {
		zap.S().Warnf("Doesn't exists: %#v\n", name)
	} else {
		zap.S().Panic(err.Error())
	}
}

func (client kubeClient) CreateShardedCluster(name string, mongoVersion string, totalMembersPerShard int, totalShards int, totalCfgServers int, totalMongos int) {
	deployment := &mongodeployments.MongoDbShardedCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MongoDbReplicaSet",
			APIVersion: "mongodb.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: client.namespace,
		},
		MongoDbShardedClusterSpec: mongodeployments.MongoDbShardedClusterSpec{
			Persistent:           false,
			Version:              mongoVersion,
			Credentials:          "dredd-om-credentials",
			Project:              "dredd-project",
			ConfigServerCount:    totalCfgServers,
			ShardCount:           totalShards,
			MongosCount:          totalMongos,
			MongodsPerShardCount: totalMembersPerShard,
		},
	}

	result, err := client.deployments.CreateMongoDbShardedCluster(deployment)
	if err != nil {
		zap.S().Error(err.Error())
	} else if err == nil {
		zap.S().Infof("Created: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		zap.S().Warnf("Already exists: %#v\n", result)
	} else {
		zap.S().Panic(err.Error())
	}
}
func (client kubeClient) DeleteShardedCluster(name string) {
	deleteOptions := &metav1.DeleteOptions{}
	err := client.deployments.DeleteMongoDbShardedCluster(name, deleteOptions)
	if err != nil {
		zap.S().Error(err.Error())
	} else if err == nil {
		zap.S().Infof("Deleted: %#v\n", name)
	} else if apierrors.IsGone(err) {
		zap.S().Warnf("Doesn't exists: %#v\n", name)
	} else {
		zap.S().Panic(err.Error())
	}
}
