package kubernetes_api

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type kubeClient struct {
	corev1    *kubernetes.Clientset
	namespace string
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
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return kubeClient{corev1: clientset, namespace: nsName}
}

func (client kubeClient) CreateEnvironment() {
	var err error
	getOpts := metav1.GetOptions{}

	_, err = client.corev1.CoreV1().Namespaces().Get(client.namespace, getOpts)
	if err != nil {
		log.Output(0, err.Error())
		client.createNamespace(client.namespace)
	}

	_, err = client.corev1.RbacV1().ClusterRoles().Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		log.Output(0, err.Error())
		client.createClusterRole()
	}

	_, err = client.corev1.CoreV1().ServiceAccounts(client.namespace).Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		log.Output(0, err.Error())
		client.createServiceAccount()
	}

	_, err = client.corev1.RbacV1().ClusterRoleBindings().Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		log.Output(0, err.Error())
		client.createClusterRoleBinding()
	}

	//TODO Create config map and secret

	_, err = client.corev1.AppsV1().Deployments(client.namespace).Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		log.Output(0, err.Error())
		client.createOperator()
	}

}

func (client kubeClient) DeleteEnvironment() {
	var err error
	getOpts := metav1.GetOptions{}

	//TODO Delete config map and secret

	_, err = client.corev1.AppsV1().Deployments(client.namespace).Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		log.Output(0, err.Error())
	} else {
		client.deleteOperator()
	}

	_, err = client.corev1.RbacV1().ClusterRoleBindings().Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		log.Output(0, err.Error())
	} else {
		client.deleteClusterRoleBinding()
	}

	_, err = client.corev1.CoreV1().ServiceAccounts(client.namespace).Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		log.Output(0, err.Error())
	} else {
		client.deleteServiceAccount()
	}

	_, err = client.corev1.RbacV1().ClusterRoles().Get("mongodb-enterprise-operator", getOpts)
	if err != nil {
		log.Output(0, err.Error())
	} else {
		client.deleteClusterRole()
	}

	_, err = client.corev1.CoreV1().Namespaces().Get(client.namespace, getOpts)
	if err != nil {
		log.Output(0, err.Error())
	} else {
		client.deleteNamespace(client.namespace)
	}
}

func (client kubeClient) createNamespace(ns string) {
	nsSpec := &apiv1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}
	result, err := client.corev1.CoreV1().Namespaces().Create(nsSpec)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created namespace %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteNamespace(ns string) {
	deleteOptions := metav1.DeleteOptions{}
	err := client.corev1.CoreV1().Namespaces().Delete(ns, &deleteOptions)
	if err != nil {
		log.Output(0, err.Error())
	}
}

func (client kubeClient) createClusterRole() {
	rules := make([]rbacv1.PolicyRule, 0, 4)

	rules = append(rules, rbacv1.PolicyRule{Verbs: []string{"get", "list", "create", "update", "delete"},
		APIGroups: []string{" "},
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
		panic(err.Error())
	}
	fmt.Printf("Created cluster role %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteClusterRole() {
	deleteOptions := metav1.DeleteOptions{}
	err := client.corev1.RbacV1().ClusterRoles().Delete("mongodb-enterprise-operator", &deleteOptions)
	if err != nil {
		log.Output(0, err.Error())
	}
}

func (client kubeClient) createServiceAccount() {
	serviceAccount := apiv1.ServiceAccount{
		TypeMeta:   metav1.TypeMeta{Kind: "ServiceAccount", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "mongodb-enterprise-operator", Namespace: client.namespace},
	}

	result, err := client.corev1.CoreV1().ServiceAccounts(client.namespace).Create(&serviceAccount)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created service account %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteServiceAccount() {
	deleteOptions := metav1.DeleteOptions{}
	err := client.corev1.CoreV1().ServiceAccounts(client.namespace).Delete("mongodb-enterprise-operator", &deleteOptions)
	if err != nil {
		log.Output(0, err.Error())
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
		panic(err.Error())
	}
	fmt.Printf("Created cluster role binding %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteClusterRoleBinding() {
	deleteOptions := metav1.DeleteOptions{}
	err := client.corev1.RbacV1().ClusterRoleBindings().Delete("mongodb-enterprise-operator", &deleteOptions)
	if err != nil {
		log.Output(0, err.Error())
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
							Image:           "quay.io/mongodb/mongodb-enterprise-operator:latest",
							ImagePullPolicy: "Always",
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{
									Name:  "OPERATOR_ENV",
									Value: "local",
								},
								apiv1.EnvVar{
									Name:  "MONGODB_ENTERPRISE_DATABASE_IMAGE",
									Value: "quay.io/mongodb/mongodb-enterprise-database:latest",
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
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func (client kubeClient) deleteOperator() {
	deleteOptions := &metav1.DeleteOptions{}
	err := client.corev1.AppsV1().Deployments(client.namespace).Delete("mongodb-enterprise-operator", deleteOptions)
	if err != nil {
		log.Output(0, err.Error())
	}
}
