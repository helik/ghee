package controller

import (
	"github.com/asdine/storm"
)

const (
	bucketName = "clusters"
	dbPath     = "ghee.db"
)

type Cluster struct {
	ID            int    `storm:"id,increment"`
	Name          string `storm:"unique"`
	Address       string
	CertAuthority []byte
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
		Name:          name,
		Address:       address,
		CertAuthority: certAuthority,
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
