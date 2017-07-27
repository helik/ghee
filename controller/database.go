package controller

import (
	"github.com/asdine/storm"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	bucketName = "clusters"
	dbPath     = "ghee.db"
)

type Cluster struct {
	ID                   int    `storm:"id,increment"`
	Name                 string `storm:"unique"`
	Address              string
	CertificateAuthority []byte
	ClientCertificate    []byte
	ClientKey            []byte
}

// Create and persist a new Kubernetes cluster.
// name: unique, friendly cluster name
// address: URL of cluster API
// certAuthority: PEM-encoded certificate authority
func AddCluster(name string, address string, certAuthority []byte) error {
	//c := api.NewCluster()

	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	cluster := Cluster{
		Name:                 name,
		Address:              address,
		CertificateAuthority: certAuthority,
		ClientCertificate:    []byte("dfd"),
		ClientKey:            []byte("dfd"),
	}
	err = db.Save(&cluster)

	return err
}

// Get one cluster by its name.
func GetCluster(name string) (Cluster, error) {
	cluster := Cluster{}
	db, err := storm.Open(dbPath)
	if err != nil {
		return cluster, err
	}

	err = db.One("Name", name, &cluster)
	if err != nil {
		return cluster, err
	}
	return cluster, nil
}

// Get a slice of all clusters.
func GetClusters() ([]Cluster, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	clusters := []Cluster{}
	db.All(&clusters)
	return clusters, nil
}

func (c Cluster) restConfig() (*rest.Config, error) {
	config := *api.NewConfig()

	// make k8s cluster struct
	cluster := api.NewCluster()
	cluster.Server = c.Address
	cluster.CertificateAuthorityData = c.CertificateAuthority
	config.Clusters[c.Name] = cluster

	// make k8s authinfo struct
	authInfo := api.NewAuthInfo()
	authInfo.ClientCertificateData = c.ClientCertificate
	authInfo.ClientKeyData = c.ClientKey
	config.AuthInfos[c.Name] = authInfo

	// make k8s context struct
	context := api.NewContext()
	context.Cluster = c.Name
	context.AuthInfo = c.Name
	context.Namespace = "default" // TODO: probably needs to be configurable
	config.Contexts[c.Name] = context

	// set current context to the previous (only) context
	config.CurrentContext = c.Name

	// get rest.Config
	clientConfig := clientcmd.NewDefaultClientConfig(config, &clientcmd.ConfigOverrides{})
	restConfig, err := clientConfig.ClientConfig()
	return restConfig, err
}
