package compiler

import (
	"path"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func transformWorkspace(defaultBase, defaultPath string) Transform {
	return func(dst *engine.Step, src *yaml.Container, conf *config.Config) {
		workdirBase := conf.Workspace.Base
		workdirPath := conf.Workspace.Path
		if workdirBase == "" {
			workdirBase = defaultBase
		}
		if workdirPath == "" {
			workdirPath = defaultPath
		}
		if !assertService(dst, src) {
			dst.WorkingDir = path.Join(workdirBase, workdirPath)
		}
		if !assertDefaultVolume(dst) {
			volume := &engine.VolumeMapping{
				Name:   "default",
				Target: workdirBase,
			}
			dst.Volumes = append(dst.Volumes, volume)
		}
	}
}

func assertDefaultVolume(dst *engine.Step) bool {
	for _, volume := range dst.Volumes {
		if volume.Name == "default" {
			return true
		}
	}
	return false
}

func assertService(dst *engine.Step, src *yaml.Container) bool {
	return len(src.Commands) == 0 && dst.Detached
}
