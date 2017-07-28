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

func (gf *Gheefile) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux []struct {
		Manifests []interface{}
		Clusters  []string
		Replicas  map[string]int32
	}
	if err := unmarshal(&aux); err != nil {
		return err
	}
	for _, auxResource := range aux {
		resource := GheeResource{
			Manifests: [][]byte{},
			Clusters:  auxResource.Clusters,
			Replicas:  auxResource.Replicas,
		}
		for _, manifest := range auxResource.Manifests {
			byteManifest, err := yaml.Marshal(manifest)
			if err != nil {
				return err
			}
			resource.Manifests = append(resource.Manifests, byteManifest)
		}
		*gf = append(*gf, resource)
	}
	return nil
}
