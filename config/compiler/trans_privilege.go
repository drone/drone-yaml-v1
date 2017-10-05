package compiler

import (
	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func transformPrivilege(images ...string) Transform {
	return func(dst *engine.Step, src *yaml.Container, conf *config.Config) {
		for _, image := range images {
			if len(src.Commands) != 0 {
				return
			}
			if matchImage(dst.Image, image) {
				dst.Privileged = true
				dst.Command = []string{}
				dst.Entrypoint = []string{}
			}
		}
	}
}
