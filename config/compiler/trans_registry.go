package compiler

import (
	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func transformRegistry(registries ...Registry) Transform {
	return func(dst *engine.Step, _ *yaml.Container, _ *config.Config) {
		for _, registry := range registries {
			if matchHostname(dst.Image, registry.Hostname) {
				dst.AuthConfig.Username = registry.Username
				dst.AuthConfig.Password = registry.Password
				dst.AuthConfig.Email = registry.Email
				break
			}
		}
	}
}
