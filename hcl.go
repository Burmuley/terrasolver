package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type TgDependenciesDoc struct {
	Dependencies []Dependency `hcl:"dependency,block"`
	Remain       interface{}  `hcl:",remain"`
}

type Dependency struct {
	Name       string      `hcl:"name,label"`
	ConfigPath string      `hcl:"config_path"`
	Remain     interface{} `hcl:",remain"`
}

func ParseDependencies(f string) ([]string, error) {
	deps := make([]string, 0)

	var tgDepsDoc TgDependenciesDoc
	err := hclsimple.DecodeFile(f, nil, &tgDepsDoc)

	if err != nil {
		return nil, fmt.Errorf("error parsing HCL: %s", err)
	}

	for _, dep := range tgDepsDoc.Dependencies {
		deps = append(deps, dep.ConfigPath)
	}

	return deps, nil
}
