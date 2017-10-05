package compiler

import (
	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

// Transform transforms container configuration to runtime configuration.
type Transform func(*engine.Step, *yaml.Container, *config.Config)
