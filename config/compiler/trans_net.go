package compiler

import (
	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func transformNetwork(network ...string) Transform {
	return func(dst *engine.Step, src *yaml.Container, conf *config.Config) {
		for _, name := range network {
			dst.Networks = append(dst.Networks, &engine.NetworkMapping{Name: name})
		}
	}
}
