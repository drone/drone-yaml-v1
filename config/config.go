package config

import "github.com/drone/drone-yaml-v1/yaml"

type (
	// Config represents the pipeline configuration.
	Config struct {
		Platform  string
		Version   yaml.StringInt
		Branches  yaml.Constraint
		Clone     yaml.Containers
		Pipeline  yaml.Containers
		Services  yaml.Containers
		Labels    yaml.SliceMap
		Networks  map[string]Network
		Volumes   map[string]Volume
		Workspace Workspace
	}

	// Workspace represents the pipeline workspace configuraiton.
	Workspace struct {
		Base string
		Path string
	}

	// Volume represents the container volume configuration.
	Volume struct {
		Driver     string
		DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	}

	// Network represents the container network configuration.
	Network struct {
		Driver     string
		DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	}
)
