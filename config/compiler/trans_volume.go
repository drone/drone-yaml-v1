package compiler

import (
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func transformVolume(volume ...string) Transform {
	return func(dst *engine.Step, src *yaml.Container, conf *config.Config) {
		for _, mapping := range volume {
			var (
				name   string
				source string
				target string
			)
			parts := strings.Split(mapping, ":")
			if len(parts) != 2 {
				continue
			}
			if strings.HasPrefix(parts[0], "/") {
				source = parts[0]
				target = parts[1]
			} else {
				name = parts[0]
				target = parts[1]
			}
			dst.Volumes = append(dst.Volumes, &engine.VolumeMapping{
				Name:   name,
				Source: source,
				Target: target,
			})
		}
	}
}
