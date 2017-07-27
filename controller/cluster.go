package controller

import (
	"github.com/ghodss/yaml"
	apps "k8s.io/api/apps/v1beta1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	rbacv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"log"
)

const (
	clusterRole        = "ClusterRole"
	clusterRoleBinding = "ClusterRoleBinding"
	configMap          = "ConfigMap"
	deployment         = "Deployment"
	namespace          = "Namespace"
	role               = "Role"
	roleBinding        = "RoleBinding"
	secret             = "Secret"
	service            = "Service"
	serviceAccount     = "ServiceAccount"
	statefulSet        = "StatefulSet"
)

type cluster struct {
	clientSet *kubernetes.Clientset
	apps      appsv1beta1.AppsV1beta1Interface
	core      corev1.CoreV1Interface
	rbac      rbacv1beta1.RbacV1beta1Interface
}

func (c *cluster) createMany(resources [][]byte, replicaCount int32) {
	for _, resource := range resources {
		c.createResource(resource, replicaCount)
	}
}

func (c *cluster) createResource(data []byte, replicaCount int32) {
	var newObj metav1.Object
	var err error
	var t struct {
		Kind string
	}
	yaml.Unmarshal(data, &t)
	switch t.Kind {
	// have a replica count
	case deployment:
		d := apps.Deployment{}
		yaml.Unmarshal(data, &d)
		d.Spec.Replicas = &replicaCount
		newObj, err = c.apps.Deployments(d.Namespace).Create(&d)
	case statefulSet:
		ss := apps.StatefulSet{}
		yaml.Unmarshal(data, &ss)
		ss.Spec.Replicas = &replicaCount
		newObj, err = c.apps.StatefulSets(ss.Namespace).Create(&ss)
		// do not have a replica count
	case clusterRole:
		cr := rbac.ClusterRole{}
		yaml.Unmarshal(data, &cr)
		newObj, err = c.rbac.ClusterRoles().Create(&cr)
	case clusterRoleBinding:
		crb := rbac.ClusterRoleBinding{}
		yaml.Unmarshal(data, &crb)
		newObj, err = c.rbac.ClusterRoleBindings().Create(&crb)
	case configMap:
		cm := core.ConfigMap{}
		yaml.Unmarshal(data, &cm)
		newObj, err = c.core.ConfigMaps(cm.Namespace).Create(&cm)
	case namespace:
		ns := core.Namespace{}
		yaml.Unmarshal(data, &ns)
		newObj, err = c.core.Namespaces().Create(&ns)
	case role:
		r := rbac.Role{}
		yaml.Unmarshal(data, &r)
		newObj, err = c.rbac.Roles(r.Namespace).Create(&r)
	case roleBinding:
		rb := rbac.RoleBinding{}
		yaml.Unmarshal(data, &rb)
		newObj, err = c.rbac.RoleBindings(rb.Namespace).Create(&rb)
	case secret:
		s := core.Secret{}
		yaml.Unmarshal(data, &s)
		newObj, err = c.core.Secrets(s.Namespace).Create(&s)
	case service:
		s := core.Service{}
		yaml.Unmarshal(data, &s)
		newObj, err = c.core.Services(s.Namespace).Create(&s)
	case serviceAccount:
		sa := core.ServiceAccount{}
		yaml.Unmarshal(data, &sa)
		newObj, err = c.core.ServiceAccounts(sa.Namespace).Create(&sa)
	default:
		log.Println("Unknown resource '" + t.Kind + "'")
		return
	}
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Created", newObj.GetName(), "in cluster", newObj.GetClusterName())
	}
}
