package controller

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"sync"
)

type Controller struct {
	clusters map[string]*cluster
}

func Make(configs map[string]*rest.Config) *Controller {
	controller := Controller{
		clusters: make(map[string]*cluster),
	}
	for clusterName, config := range configs {
		controller.clusters[clusterName] = createCluster(config)
	}
	return &controller
}

func createCluster(config *rest.Config) *cluster {
	var clientSet *kubernetes.Clientset
	var err error
	if clientSet, err = kubernetes.NewForConfig(config); err != nil {
		log.Fatal(err)
	}
	return &cluster{
		clientSet: clientSet,
		apps:      clientSet.AppsV1beta1(),
		core:      clientSet.CoreV1(),
		rbac:      clientSet.RbacV1beta1(),
	}
}

func (c *Controller) Create(manifest Gheefile) {
	for _, resource := range manifest {
		c.createResource(resource)
	}
}

func (c *Controller) createResource(resource GheeResource) {
	wg := sync.WaitGroup{}
	wg.Add(len(resource.Clusters))
	for _, clusterName := range resource.Clusters {
		go func(cluster *cluster) {
			cluster.createMany(resource.Manifests, getWithDefault(resource.Replicas, clusterName, 1))
			wg.Done()
		}(c.clusters[clusterName])
	}
	wg.Wait()
}

func getWithDefault(replicas map[string]int32, key string, defaultVal int32) int32 {
	if v, present := replicas[key]; present {
		return v
	}
	return defaultVal
}
