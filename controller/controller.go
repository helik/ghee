package controller

import (
	"github.com/ghodss/yaml"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

type Controller struct {
	clientSet *kubernetes.Clientset
}

func MakeInCluster() *Controller {
	var config *rest.Config
	var err error
	if config, err = rest.InClusterConfig(); err != nil {
		panic(err.Error())
	}
	return createController(config)
}

func MakeOutOfCluster(configPath string) *Controller {
	var config *rest.Config
	var err error
	if config, err = clientcmd.BuildConfigFromFlags("", configPath); err != nil {
		log.Fatal(err)
	}
	return createController(config)
}

func createController(config *rest.Config) *Controller {
	var clientSet *kubernetes.Clientset
	var err error
	if clientSet, err = kubernetes.NewForConfig(config); err != nil {
		log.Fatal(err)
	}
	return &Controller{clientSet}
}

func (c *Controller) Create(kind string, data []byte) {
	switch kind {
	case "ServiceAccount": //TODO find constants somewhere
		sa := core.ServiceAccount{}
		yaml.Unmarshal(data, &sa)
		newSa, err := c.clientSet.CoreV1().ServiceAccounts(sa.Namespace).Create(&sa)
		if err != nil {
			log.Println(err)
		}
		log.Println("Created service account", newSa.Name)
	case "ClusterRole": //RBAC
		cr := rbac.ClusterRole{}
		yaml.Unmarshal(data, &cr)
		newCr, err := c.clientSet.RbacV1beta1().ClusterRoles().Create(&cr)
		if err != nil {
			log.Println(err)
		}
		log.Println("Created cluster role", newCr.Name)
	case "ClusterRoleBinding":
		crb := rbac.ClusterRoleBinding{}
		yaml.Unmarshal(data, &crb)
		newCrb, err := c.clientSet.RbacV1beta1().ClusterRoleBindings().Create(&crb)
		if err != nil {
			log.Println(err)
		}
		log.Println("Created cluster role binding", newCrb.Name)
	case "Role":
		r := rbac.Role{}
		yaml.Unmarshal(data, &r)
		newR, err := c.clientSet.RbacV1beta1().Roles(r.Namespace).Create(&r)
		if err != nil {
			log.Println(err)
		}
		log.Println("Created cluster role", newR.Name)
	case "RoleBinding":
		rb := rbac.RoleBinding{}
		yaml.Unmarshal(data, &rb)
		newRb, err := c.clientSet.RbacV1beta1().RoleBindings(rb.Namespace).Create(&rb)
		if err != nil {
			log.Println(err)
		}
		log.Println("Created cluster role binding", newRb.Name)
	default:
		log.Println("Cannot create resource '" + kind + "'")
	}
}