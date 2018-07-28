package config

import "github.com/drone/drone-yaml-v1/yaml"

type (
	// Config represents the pipeline configuration.
	Config struct {
		Metadata  Metadata
		Platform  Platform
		Clone     Clone
		Workspace Workspace
		Version   yaml.StringInt
		DependsOn yaml.StringSlice `yaml:"depends_on"`
		Trigger   yaml.Constraints
		Labels    yaml.SliceMap
		Pipeline  []map[string]*yaml.Container
		Services  map[string]*yaml.Container
		Networks  map[string]Network
		Volumes   map[string]Volume
		Secrets   map[string]Secret
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
		Aescbc     string
		Aesgcm     string
		Secretbox  string
		Driver     string
		DriverOpts map[string]interface{} `yaml:"driver_opts,omitempty"`
	}

	// Platform provides platform details
	Platform struct {
		Name string
	}

	// Metadata provides platform details
	Metadata struct {
		Name string
	}

	// Clone provides clone customization
	Clone struct {
		Disabled bool
		Depth    int
	}
)
