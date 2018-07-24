package config

import "github.com/drone/drone-yaml-v1/yaml"

type (
	// Config represents the pipeline configuration.
	Config struct {
		Name      string
		Platform  string
		Version   yaml.StringInt
		DependsOn yaml.StringSlice `yaml:"depends_on"`
		Branches  yaml.Constraint
		Trigger   yaml.Constraints
		Files     yaml.Constraint
		Clone     *yaml.Container
		Pipeline  yaml.Pipeline
		Services  yaml.Containers
		Labels    yaml.SliceMap
		Networks  map[string]Network
		Volumes   map[string]Volume
		Secrets   map[string]Secret
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

	// Secret represents the container secret configuration.
	Secret struct {
		External   yaml.External
		File       string
		Secretbox  string
		Driver     string
		DriverOpts map[string]interface{} `yaml:"driver_opts,omitempty"`
	}
)
