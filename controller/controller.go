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

func Create(manifest Gheefile) {
	c := makeController(manifest)
	for _, resource := range manifest {
		c.createResource(resource)
	}
}

func Delete(manifest Gheefile) {
	c := makeController(manifest)
	for _, resource := range manifest {
		c.deleteResource(resource)
	}
}

func makeController(manifest Gheefile) *Controller {
	controller := Controller{
		clusters: make(map[string]*cluster),
	}
	// add all needed configs to the controller
	for _, resource := range manifest {
		for _, clusterName := range resource.Clusters {
			if _, present := controller.clusters[clusterName]; !present {
				// TODO get config from db
				config := rest.Config{}
				controller.clusters[clusterName] = createCluster(clusterName, &config)
			}
		}
	}
	return &controller
}

func createCluster(name string, config *rest.Config) *cluster {
	var clientSet *kubernetes.Clientset
	var err error
	if clientSet, err = kubernetes.NewForConfig(config); err != nil {
		log.Fatal(err)
	}
	return &cluster{
		name:      name,
		clientSet: clientSet,
		apps:      clientSet.AppsV1beta1(),
		core:      clientSet.CoreV1(),
		rbac:      clientSet.RbacV1beta1(),
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

func (c *Controller) deleteResource(resource GheeResource) {
	wg := sync.WaitGroup{}
	wg.Add(len(resource.Clusters))
	for _, clusterName := range resource.Clusters {
		go func(cluster *cluster) {
			cluster.deleteMany(resource.Manifests)
			wg.Done()
		}(c.clusters[clusterName])
	}
	wg.Wait()
}

// helper functions

func getWithDefault(replicas map[string]int32, key string, defaultVal int32) int32 {
	if v, present := replicas[key]; present {
		return v
	}
	return defaultVal
}
