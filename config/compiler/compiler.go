package compiler

import (
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-runtime/version"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

type (
	// Resources represents the container resource limits.
	Resources struct {
		MemSwapLimit int64
		MemLimit     int64
		ShmSize      int64
		CPUQuota     int64
		CPUShares    int64
		CPUSet       string
	}

	// Registry represents registry credentials used to
	// pull private images from a docker registry.
	Registry struct {
		Hostname string
		Username string
		Password string
		Email    string
		Token    string
	}

	// Secret represents a repository secret that should
	// be passed to the container at runtime.
	Secret struct {
		Name  string
		Value string
		Match []string
	}

	// Metadata represents pipeline metadata required to
	// filter pipeline steps.
	Metadata struct {
		Ref         string
		Repo        string
		Platform    string
		Environment string
		Event       string
		Branch      string
		Matrix      map[string]string
	}
)

// A Compiler compiles a pipeline configuration to an intermediate runtime.
type Compiler struct {
	metadata   Metadata
	noclone    bool
	transforms []Transform
}

// New returns a new compiler
func New(opts ...Option) *Compiler {
	c := &Compiler{
		transforms: []Transform{
			transformPlugin,
			transformCommand,
			transformProxy(),
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Compile compiles the parsed yaml configuration and converts to the
// drone runtime intermediate representation.
func (c *Compiler) Compile(conf *config.Config) (*engine.Config, error) {
	spec := new(engine.Config)
	spec.Version = version.VersionMajor

	if _, ok := conf.Networks["default"]; !ok {
		dst := &engine.Network{Driver: "bridge", Name: "default"}
		if conf.Platform.Name == "windows/amd64" {
			dst.Driver = "nat"
		}
		spec.Networks = append(spec.Networks, dst)
	}

	if _, ok := conf.Volumes["default"]; !ok {
		dst := &engine.Volume{Driver: "local", Name: "default"}
		spec.Volumes = append(spec.Volumes, dst)
	}

	for name, src := range conf.Networks {
		dst := &engine.Network{Driver: src.Driver, Name: name}
		spec.Networks = append(spec.Networks, dst)
	}

	for name, src := range conf.Volumes {
		dst := &engine.Volume{Driver: src.Driver, Name: name}
		spec.Volumes = append(spec.Volumes, dst)
	}

	if c.noclone == false && conf.Clone.Disable == false {
		image := "drone/git"
		switch conf.Platform.Name {
		case "linux/arm":
			image = "drone/git:linux-arm"
		case "linux/arm64":
			image = "drone/git:linux-arm64"
		case "windows/amd64":
			image = "drone/git:windows-1803"
		}

		dst := &engine.Step{}
		src := &yaml.Container{Name: "clone", Image: image}
		copyContainer(dst, src)
		for _, t := range c.transforms {
			t(dst, src, conf)
		}
		stage := new(engine.Stage)
		spec.Stages = append(spec.Stages, stage)
		stage.Steps = append(stage.Steps, dst)
	}

	if len(conf.Services) != 0 {
		stage := new(engine.Stage)
		for name, src := range conf.Services {
			src.Name = name
			dst := new(engine.Step)
			copyService(dst, src)
			for _, t := range c.transforms {
				t(dst, src, conf)
			}
			if calcSkip(src, c.metadata) {
				continue
			}
			stage.Steps = append(stage.Steps, dst)
		}
		if len(stage.Steps) != 0 {
			spec.Stages = append(spec.Stages, stage)
		}
	}

	for _, group := range conf.Pipeline {
		stage := new(engine.Stage)
		for name, src := range group {
			src.Name = name
			dst := new(engine.Step)
			copyContainer(dst, src)
			for _, t := range c.transforms {
				t(dst, src, conf)
			}
			if calcSkip(src, c.metadata) {
				continue
			}
			stage.Steps = append(stage.Steps, dst)
		}
		if len(stage.Steps) != 0 {
			spec.Stages = append(spec.Stages, stage)
		}
	}

	namespace(spec, conf)
	return spec, nil
}

// helper function copies the service contianer configuration from the
// yaml container to the engine container representation.
func copyService(dst *engine.Step, src *yaml.Container) {
	copyContainer(dst, src)
	dst.Detached = true
}

// helper function copies the contianer configuration from the yaml
// container to the engine container representation.
func copyContainer(dst *engine.Step, src *yaml.Container) {
	dst.Name = src.Name
	dst.Image = expandImage(src.Image)
	dst.Pull = src.Pull
	dst.Detached = src.Detached
	dst.Privileged = src.Privileged
	dst.Environment = src.Environment.Map
	dst.Labels = src.Labels.Map
	dst.Entrypoint = src.Entrypoint
	dst.Command = src.Command
	dst.ExtraHosts = src.ExtraHosts
	dst.Tmpfs = src.Tmpfs
	dst.DNS = src.DNS
	dst.DNSSearch = src.DNSSearch
	dst.NetworkMode = src.NetworkMode
	dst.IpcMode = src.IpcMode
	dst.Sysctls = src.Sysctls.Map
	dst.ErrIgnore = src.ErrIgnore
	dst.OnSuccess = calcOnSucess(src)
	dst.OnFailure = calcOnFailure(src)
	if dst.Environment == nil {
		dst.Environment = map[string]string{}
	}
	dst.Environment["DRONE_STEP"] = dst.Name

	defaultContainerNetwork(dst, src)
	copyContainerReports(dst, src)
	copyContainerVolume(dst, src)
	copyContainerNetwork(dst, src)
	copyContainerDevices(dst, src)
}

// helper function set the default container network.
func defaultContainerNetwork(dst *engine.Step, src *yaml.Container) {
	network := &engine.NetworkMapping{
		Name:    "default",
		Aliases: []string{dst.Name},
	}
	dst.Networks = append(dst.Networks, network)
}

// helper function copies the volume configuration from the yaml
// container to the engine container representation.
func copyContainerVolume(dst *engine.Step, src *yaml.Container) {
	for _, vol := range src.Volumes {
		volume := &engine.VolumeMapping{}
		volume.Target = vol.Destination
		if strings.HasPrefix(vol.Source, "/") {
			volume.Source = vol.Source
		} else {
			volume.Name = vol.Source
		}
		dst.Volumes = append(dst.Volumes, volume)
	}
}

// helper function copies the network configuration from the yaml
// container to the engine container representation.
func copyContainerNetwork(dst *engine.Step, src *yaml.Container) {
	for _, net := range src.Networks.Networks {
		network := &engine.NetworkMapping{}
		network.Name = net.Name
		network.Aliases = net.Aliases
		if len(network.Aliases) == 0 {
			network.Aliases = []string{dst.Name}
		}
		dst.Networks = append(dst.Networks, network)
	}
}

func copyContainerReports(dst *engine.Step, src *yaml.Container) {
	if src.Reports.Coverage != nil {
		export := new(engine.File)
		export.Path = src.Reports.Coverage.Source
		export.Mime = src.Reports.Coverage.Format
		dst.Exports = append(dst.Exports, export)
	}
}

func copyContainerDevices(dst *engine.Step, src *yaml.Container) {
	for _, dev := range src.Devices {
		parts := strings.Split(dev, ":")
		if len(parts) != 2 {
			continue
		}
		dst.Devices = append(dst.Devices, &engine.DeviceMapping{
			Source: parts[0],
			Target: parts[1],
		})
	}
}

// calculate container on_success
func calcOnSucess(src *yaml.Container) bool {
	return src.Constraints.Status.Match("success")
}

// calculate container on_failure
func calcOnFailure(src *yaml.Container) bool {
	return (len(src.Constraints.Status.Include)+
		len(src.Constraints.Status.Exclude) != 0) &&
		src.Constraints.Status.Match("failure")
}

func calcSkip(src *yaml.Container, metadata Metadata) bool {
	return !(src.Constraints.Platform.Match(metadata.Platform) &&
		src.Constraints.Environment.Match(metadata.Environment) &&
		src.Constraints.Event.Match(metadata.Event) &&
		src.Constraints.Branch.Match(metadata.Branch) &&
		src.Constraints.Repo.Match(metadata.Repo) &&
		src.Constraints.Ref.Match(metadata.Ref) &&
		src.Constraints.Matrix.Match(metadata.Matrix))
}
