package controller

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Gheefile []GheeResource

type GheeResource struct {
	Manifest interface{}
	Clusters []string
	Replicas map[string]int
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
