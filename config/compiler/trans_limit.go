package compiler

import (
	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func transformLimits(limits Resources) Transform {
	return func(dst *engine.Step, src *yaml.Container, conf *config.Config) {
		dst.MemSwapLimit = limits.MemSwapLimit
		dst.MemLimit = limits.MemLimit
		dst.ShmSize = limits.ShmSize
		dst.CPUQuota = limits.CPUQuota
		dst.CPUShares = limits.CPUShares
		dst.CPUSet = limits.CPUSet
	}
}
