package compiler

import (
	"net/url"
	"path/filepath"
)

// Option set a compiler option.
type Option func(*Compiler)

// WithTransform returns a compiler option to set transform t.
func WithTransform(t Transform) Option {
	return func(c *Compiler) {
		c.transforms = append(c.transforms, t)
	}
}

// WithMetadata returns a compiler option to set metadata.
func WithMetadata(m Metadata) Option {
	return func(c *Compiler) {
		c.metadata = m
	}
}

// WithClone returns a compiler option to clone.
func WithClone(clone bool) Option {
	return func(c *Compiler) {
		c.noclone = !clone
	}
}

// WithVolumes configutes the compiler with default volumes that
// are mounted to each container in the pipeline.
func WithVolumes(volumes ...string) Option {
	return WithTransform(
		transformVolume(volumes...),
	)
}

// WithRegistry configures the compiler with registry credentials
// that should be used to download images.
func WithRegistry(registries ...Registry) Option {
	return WithTransform(
		transformRegistry(registries...),
	)
}

// WithSecret configures the compiler with external secrets
// to be injected into the container at runtime.
func WithSecret(secrets ...Secret) Option {
	return WithTransform(
		transformSecret(secrets...),
	)
}

// WithNetrc configures the compiler with netrc authentication
// credentials added by default to every container in the pipeline.
func WithNetrc(username, password, machine string) Option {
	return WithEnviron(
		map[string]string{
			"CI_NETRC_USERNAME":    username,
			"CI_NETRC_PASSWORD":    password,
			"CI_NETRC_MACHINE":     machine,
			"DRONE_NETRC_USERNAME": username,
			"DRONE_NETRC_PASSWORD": password,
			"DRONE_NETRC_MACHINE":  machine,
		},
	)
}

// WithWorkspace configures the compiler with the workspace base
// and path. The workspace base is a volume created at runtime and
// mounted into all containers in the pipeline. The base and path
// are joined to provide the working directory for all build and
// plugin steps in the pipeline.
func WithWorkspace(base, path string) Option {
	return WithTransform(
		transformWorkspace(base, path),
	)
}

// WithWorkspaceFromURL configures the compiler with the workspace
// base and path based on the repository url.
func WithWorkspaceFromURL(base, link string) Option {
	path := "src"
	parsed, err := url.Parse(link)
	if err == nil {
		path = filepath.Join(path, parsed.Hostname(), parsed.Path)
	}
	return WithWorkspace(base, path)
}

// WithPrivileged configures the compiler to automatically execute
// images as privileged containers if the match the given list.
func WithPrivileged(images ...string) Option {
	return WithTransform(
		transformPrivilege(images...),
	)
}

// WithEnviron configures the compiler with environment variables
// added by default to every container in the pipeline.
func WithEnviron(env map[string]string) Option {
	return WithTransform(
		transformEnv(env),
	)
}

// WithNetworks configures the compiler with additionnal networks
// to be connected to build containers
func WithNetworks(networks ...string) Option {
	return WithTransform(
		transformNetwork(networks...),
	)
}

// WithLimits configures the compiler with default resource limits that
// are applied each container in the pipeline.
func WithLimits(limits Resources) Option {
	return WithTransform(
		transformLimits(limits),
	)
}
