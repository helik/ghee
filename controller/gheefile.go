package controller

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/helik/ghee/database"
)

type GheeManifest []GheeResource

type GheeResource struct {
	Manifests [][]byte
	Clusters  []string
	Replicas  map[string]int32
}

func ReadGheeManifest(filepath string) (GheeManifest, error) {
	g := GheeManifest{}

	body, err := ioutil.ReadFile(filepath)
	if err != nil {
		return g, err
	}

	err = yaml.Unmarshal(body, &g)
	if err != nil {
		return g, err
	}

	return g, nil
}

func ReadGheeClusterFile(filepath string) (database.Cluster, error) {
	c := database.Cluster{}

	body, err := ioutil.ReadFile(filepath)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(body, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}
