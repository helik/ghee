package controller

import (
	"github.com/ghodss/yaml"
	apps "k8s.io/api/apps/v1beta1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	appsv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	rbacv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"log"
)

const (
	deployment         = "Deployment"
	statefulSet        = "StatefulSet"
	configMap          = "ConfigMap"
	namespace          = "Namespace"
	secret             = "Secret"
	service            = "Service"
	serviceAccount     = "ServiceAccount"
	clusterRole        = "ClusterRole"
	clusterRoleBinding = "ClusterRoleBinding"
	role               = "Role"
	roleBinding        = "RoleBinding"
)

type cluster struct {
	name      string
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
	var t metav1.TypeMeta
	var err error
	yaml.Unmarshal(data, &t)
	switch t.Kind {
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
		log.Println("Cluster", c.name, "error in creating", t.Kind+":", err)
	} else {
		log.Println("Cluster", c.name+":", "created", t.Kind, newObj.GetName())
	}
}

func (c *cluster) deleteMany(resources [][]byte) {
	for _, resource := range resources {
		c.deleteResource(resource)
	}
}

func (c *cluster) deleteResource(data []byte) {
	var obj struct {
		Metadata metav1.ObjectMeta
	}
	var t metav1.TypeMeta
	var err error
	yaml.Unmarshal(data, &obj)
	yaml.Unmarshal(data, &t)
	switch t.Kind {
	case deployment:
		d := apps.Deployment{}
		yaml.Unmarshal(data, &d)
		propagation := metav1.DeletePropagationForeground
		deleteOptions := metav1.DeleteOptions{
			GracePeriodSeconds: d.Spec.Template.Spec.TerminationGracePeriodSeconds,
			PropagationPolicy:  &propagation,
		}
		err = c.apps.Deployments(d.Namespace).Delete(d.Name, &deleteOptions)
	case statefulSet:
		ss := apps.StatefulSet{}
		yaml.Unmarshal(data, &ss)
		// get pods created by the statefulset & get their pvcs
		ssPods := map[string]int{}
		pvcs := []string{}
		// TODO use list options to reduce # of pods returned
		var pods *core.PodList
		pods, err = c.core.Pods(ss.Namespace).List(metav1.ListOptions{})
		if err != nil {
			break
		}
		for _, pod := range pods.Items {
			for _, owner := range pod.OwnerReferences {
				if owner.Kind == statefulSet && owner.Name == ss.Name {
					ssPods[pod.Name] = 1
					for _, vol := range pod.Spec.Volumes {
						if vol.VolumeSource.PersistentVolumeClaim != nil {
							pvcs = append(pvcs, vol.VolumeSource.PersistentVolumeClaim.ClaimName)
						}
					}
					break
				}
			}
		}
		propagation := metav1.DeletePropagationForeground
		deleteOptions := metav1.DeleteOptions{
			GracePeriodSeconds: ss.Spec.Template.Spec.TerminationGracePeriodSeconds,
			PropagationPolicy:  &propagation,
		}
		err = c.apps.StatefulSets(ss.Namespace).Delete(ss.Name, &deleteOptions)
		if err != nil {
			break
		}
		log.Println("Cluster", c.name+":", "deleted StatefulSet", ss.Name)
		// need to wait for pods to be delete before we can delete the pvcs
		ssPodsKeys := []string{}
		for key, _ := range ssPods {
			ssPodsKeys = append(ssPodsKeys, key)
		}
		log.Println("Cluster", c.name+":", "watching for pods", ssPodsKeys, "to be deleted")
		var watcher watch.Interface
		watcher, err = c.core.Pods(ss.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			break
		}
		for {
			if len(ssPods) <= 0 {
				watcher.Stop()
				break
			}
			select {
			case event := <-watcher.ResultChan():
				if event.Type == watch.Deleted {
					name := event.Object.(*core.Pod).Name
					if _, present := ssPods[name]; present {
						delete(ssPods, name)
						log.Println("Cluster", c.name+":", "Pod", name, "deleted")
					}
				}
			}
		}
		// delete pvcs
		for _, pvc := range pvcs {
			err = c.core.PersistentVolumeClaims(ss.Namespace).Delete(pvc, &metav1.DeleteOptions{})
			if err != nil {
				break
			}
			log.Println("Cluster", c.name+":", "pvc", pvc, "deleted")
		}
	case clusterRole:
		err = c.rbac.ClusterRoles().Delete(obj.Metadata.Name, &metav1.DeleteOptions{})
	case clusterRoleBinding:
		err = c.rbac.ClusterRoleBindings().Delete(obj.Metadata.Name, &metav1.DeleteOptions{})
	case configMap:
		err = c.core.ConfigMaps(obj.Metadata.Namespace).Delete(obj.Metadata.Name, &metav1.DeleteOptions{})
	case namespace:
		err = c.core.Namespaces().Delete(obj.Metadata.Name, &metav1.DeleteOptions{})
	case role:
		err = c.rbac.Roles(obj.Metadata.Namespace).Delete(obj.Metadata.Name, &metav1.DeleteOptions{})
	case roleBinding:
		err = c.rbac.RoleBindings(obj.Metadata.Namespace).Delete(obj.Metadata.Name, &metav1.DeleteOptions{})
	case secret:
		err = c.core.Secrets(obj.Metadata.Namespace).Delete(obj.Metadata.Name, &metav1.DeleteOptions{})
	case service:
		err = c.core.Services(obj.Metadata.Namespace).Delete(obj.Metadata.Name, &metav1.DeleteOptions{})
	case serviceAccount:
		err = c.core.ServiceAccounts(obj.Metadata.Namespace).Delete(obj.Metadata.Name, &metav1.DeleteOptions{})
	default:
		log.Println("Unknown resource '" + t.Kind + "'")
		return
	}
	if err != nil {
		log.Println("Cluster", c.name, "error in deleting", t.Kind, obj.Metadata.Name+":", err)
	} else {
		log.Println("Cluster", c.name+":", "deleted", t.Kind, obj.Metadata.Name)
	}
}

func (c *cluster) updateMany(resources [][]byte, replicaCount int32) {
	for _, resource := range resources {
		c.updateResource(resource, replicaCount)
	}
}

func (c *cluster) updateResource(data []byte, replicaCount int32) {
	var newObj metav1.Object
	var t metav1.TypeMeta
	var err error
	yaml.Unmarshal(data, &t)
	switch t.Kind {
	case deployment:
		d := apps.Deployment{}
		yaml.Unmarshal(data, &d)
		d.Spec.Replicas = &replicaCount
		newObj, err = c.apps.Deployments(d.Namespace).Update(&d)
	case statefulSet:
		var ss apps.StatefulSet
		yaml.Unmarshal(data, &ss)
		var currentSS *apps.StatefulSet
		currentSS, err = c.apps.StatefulSets(ss.Namespace).Get(ss.Name, metav1.GetOptions{})
		remaining := *currentSS.Spec.Replicas - replicaCount
		ss.Spec.Replicas = &replicaCount
		newObj, err = c.apps.StatefulSets(ss.Namespace).Update(&ss)
		if err != nil {
			break
		}
		// watch for deleted pods, get their pvcs, delete them
		var watcher watch.Interface
		watcher, err = c.core.Pods(ss.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			break
		}
		if remaining > 0 {
			log.Println("Cluster", c.name+":", "watching for", remaining, "pods to be deleted")
		}
		pvcs := []string{}
		for {
			// assumption: after the-difference-in-replica-count pods are deleted, the statefulset is up to date
			if remaining <= 0 {
				watcher.Stop()
				break
			}
			select {
			case event := <-watcher.ResultChan():
				if event.Type == watch.Deleted {
					pod := event.Object.(*core.Pod)
					for _, owner := range pod.OwnerReferences {
						if owner.Kind == statefulSet && owner.Name == ss.Name {
							remaining--
							log.Println("Cluster", c.name+":", "Pod", pod.Name, "deleted")
							for _, vol := range pod.Spec.Volumes {
								if vol.VolumeSource.PersistentVolumeClaim != nil {
									pvcs = append(pvcs, vol.VolumeSource.PersistentVolumeClaim.ClaimName)
								}
							}
							break
						}
					}
				}
			}
		}
		for _, pvc := range pvcs {
			c.core.PersistentVolumeClaims(ss.Namespace).Delete(pvc, &metav1.DeleteOptions{})
			if err != nil {
				break
			}
		}
	case clusterRole:
		cr := rbac.ClusterRole{}
		yaml.Unmarshal(data, &cr)
		newObj, err = c.rbac.ClusterRoles().Update(&cr)
	case clusterRoleBinding:
		crb := rbac.ClusterRoleBinding{}
		yaml.Unmarshal(data, &crb)
		newObj, err = c.rbac.ClusterRoleBindings().Update(&crb)
	case configMap:
		cm := core.ConfigMap{}
		yaml.Unmarshal(data, &cm)
		newObj, err = c.core.ConfigMaps(cm.Namespace).Update(&cm)
	case namespace:
		ns := core.Namespace{}
		yaml.Unmarshal(data, &ns)
		newObj, err = c.core.Namespaces().Update(&ns)
	case role:
		r := rbac.Role{}
		yaml.Unmarshal(data, &r)
		newObj, err = c.rbac.Roles(r.Namespace).Update(&r)
	case roleBinding:
		rb := rbac.RoleBinding{}
		yaml.Unmarshal(data, &rb)
		newObj, err = c.rbac.RoleBindings(rb.Namespace).Update(&rb)
	case secret:
		s := core.Secret{}
		yaml.Unmarshal(data, &s)
		newObj, err = c.core.Secrets(s.Namespace).Update(&s)
	case service:
		s := core.Service{}
		yaml.Unmarshal(data, &s)
		var currentService *core.Service
		currentService, err = c.core.Services(s.Namespace).Get(s.Name, metav1.GetOptions{})
		newObj, err = c.core.Services(currentService.Namespace).Update(currentService)
	case serviceAccount:
		sa := core.ServiceAccount{}
		yaml.Unmarshal(data, &sa)
		newObj, err = c.core.ServiceAccounts(sa.Namespace).Update(&sa)
	default:
		log.Println("Unknown resource '" + t.Kind + "'")
		return
	}
	if err != nil {
		log.Println("Cluster", c.name, "error in updating", t.Kind+":", err)
	} else {
		log.Println("Cluster", c.name+":", "updated", t.Kind, newObj.GetName())
	}
}
