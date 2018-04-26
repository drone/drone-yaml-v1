package linter

import (
	"fmt"

	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"

	"github.com/gosimple/slug"
)

// Check returns an error if the configuration is invalid.
type Check func(*config.Config) error

// CheckContainer is an adapter to perform a check for every container
// in the configuration. If a check fails the function halts and returns
// an error.
//
// TODO(bradrydzewski) if check container accepted a slice of checks,
// we could chain the checks together to reduce the number of iterations.
func CheckContainer(check func(*config.Config, *yaml.Container) error) Check {
	return func(conf *config.Config) error {
		if container := conf.Clone; container != nil {
			if err := check(conf, container); err != nil {
				return err
			}
		}
		for _, container := range conf.Pipeline.Containers {
			if err := check(conf, container); err != nil {
				return err
			}
		}
		for _, container := range conf.Services.Containers {
			if err := check(conf, container); err != nil {
				return err
			}
		}
		return nil
	}
}

// CheckPipeline checks the pipeline block is not empty.
func CheckPipeline(conf *config.Config) error {
	if len(conf.Pipeline.Containers) == 0 {
		return fmt.Errorf("Invalid or missing pipeline section")
	}
	return nil
}

// CheckNetworks prevents a configuration from defining custom
// networks in untrusted mode.
func CheckNetworks(trusted bool) Check {
	return func(conf *config.Config) error {
		if trusted {
			return nil
		}
		if len(conf.Networks) != 0 {
			return fmt.Errorf("Insufficient privileges to define custom networks")
		}
		return nil
	}
}

// CheckVolumes limits the configuration to only using custom
// local volumes in untrusted mode.
func CheckVolumes(trusted bool) Check {
	return func(conf *config.Config) error {
		if trusted {
			return nil
		}
		for _, volume := range conf.Volumes {
			if !(volume.Driver == "local" || volume.Driver == "") {
				return fmt.Errorf("Insufficient privileges to define custom volumes")
			}
		}
		return nil
	}
}

// CheckImage checks the container image attribute is not empty.
func CheckImage(conf *config.Config, container *yaml.Container) error {
	if len(container.Image) == 0 {
		return fmt.Errorf("Invalid or missing image")
	}
	return nil
}

// CheckCommands checks the container commands to not conflict with
// the container entrypoint and command blocks.
func CheckCommands(conf *config.Config, container *yaml.Container) error {
	if len(container.Commands) == 0 {
		return nil
	}
	if len(container.Vargs) != 0 {
		return fmt.Errorf("Cannot configure both commands and plugin attributes")
	}
	if len(container.Entrypoint) != 0 {
		return fmt.Errorf("Cannot configure both commands and entrypoint attributes")
	}
	if len(container.Command) != 0 {
		return fmt.Errorf("Cannot configure both commands and command attributes")
	}
	return nil
}

// CheckEntrypoint checks that a container is not overriding the entypoint.
func CheckEntrypoint(conf *config.Config, container *yaml.Container) error {
	if !IsService(conf, container) && len(container.Entrypoint) != 0 {
		return fmt.Errorf("Cannot override container entrypoint")
	}
	return nil
}

// CheckCommand checks that a container is not overriding the entypoint.
func CheckCommand(conf *config.Config, container *yaml.Container) error {
	if !IsService(conf, container) && len(container.Command) != 0 {
		return fmt.Errorf("Cannot override container command")
	}
	return nil
}

// CheckTrusted checks that a container is not using any restricted
// settings that require elevated permissions.
func CheckTrusted(trusted bool) Check {
	return CheckContainer(func(conf *config.Config, container *yaml.Container) error {
		if trusted {
			return nil
		}
		if container.Privileged {
			return fmt.Errorf("Insufficient privileges to use privileged mode")
		}
		if container.ShmSize != 0 {
			return fmt.Errorf("Insufficient privileges to override shm_size")
		}
		if len(container.DNS) != 0 {
			return fmt.Errorf("Insufficient privileges to use custom dns")
		}
		if len(container.DNSSearch) != 0 {
			return fmt.Errorf("Insufficient privileges to use dns_search")
		}
		if len(container.Devices) != 0 {
			return fmt.Errorf("Insufficient privileges to use devices")
		}
		if len(container.ExtraHosts) != 0 {
			return fmt.Errorf("Insufficient privileges to use extra_hosts")
		}
		if len(container.NetworkMode) != 0 && container.NetworkMode != "bridge" {
			return fmt.Errorf("Insufficient privileges to use network_mode")
		}
		if len(container.IpcMode) != 0 {
			return fmt.Errorf("Insufficient privileges to use ipc_mode")
		}
		if len(container.Sysctls.Map) != 0 {
			return fmt.Errorf("Insufficient privileges to use sysctls")
		}
		if container.Networks.Networks != nil && len(container.Networks.Networks) != 0 {
			return fmt.Errorf("Insufficient privileges to use networks")
		}
		if len(container.Volumes) != 0 {
			for _, volume := range container.Volumes {
				if !IsDataVolume(conf, volume) {
					return fmt.Errorf("Insufficient privileges to use volumes")
				}
			}
		}
		return nil
	})
}

// IsDataVolume returns true if the volume mapping is a data volume.
func IsDataVolume(conf *config.Config, volume *yaml.Volume) bool {
	if !slug.IsSlug(volume.Source) {
		return false
	}
	_, ok := conf.Volumes[volume.Source]
	return ok
}

// IsService returns true if the container is a service.
func IsService(conf *config.Config, container *yaml.Container) bool {
	if container.Detached {
		return true
	}
	for _, service := range conf.Services.Containers {
		if service == container {
			return true
		}
	}
	return false
}
