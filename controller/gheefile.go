package controller

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Gheefile []GheeResource

type GheeResource struct {
	Manifests [][]byte
	Clusters  []string
	Replicas  map[string]int32
}

func ReadGheefile(filepath string) (Gheefile, error) {
	g := Gheefile{}

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

func ReadGheeClusterFile(filepath string) (Cluster, error) {
	c := Cluster{}

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
