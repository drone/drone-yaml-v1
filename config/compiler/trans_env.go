package compiler

import (
	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func transformEnv(envs map[string]string) Transform {
	return func(dst *engine.Step, src *yaml.Container, conf *config.Config) {
		if dst.Environment == nil {
			dst.Environment = map[string]string{}
		}
		for k, v := range envs {
			if len(v) != 0 {
				dst.Environment[k] = v
			}
		}
	}
}
