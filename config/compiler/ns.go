package compiler

import (
	"fmt"

	"github.com/dchest/uniuri"
	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/gosimple/slug"
)

// namespace is responsible for namespacing the container resources to
// prevent name conflicts on a shared host.
func namespace(dst *engine.Config, src *config.Config) {
	ns := uniuri.New()
	for i, stage := range dst.Stages {
		for ii, step := range stage.Steps {
			step.Alias = step.Name
			step.Name = fmt.Sprintf("%v%d%d_%s", ns, i, ii, step.Name)
			step.Name = slug.Make(step.Name)

			for _, volume := range step.Volumes {
				if volume.Name != "" {
					volume.Name = fmt.Sprintf("%v_%s", ns, volume.Name)
				}
			}
			for _, network := range step.Networks {
				network.Name = fmt.Sprintf("%v_%s", ns, network.Name)
			}
		}
	}

	for _, volume := range dst.Volumes {
		volume.Name = fmt.Sprintf("%v_%s", ns, volume.Name)
		volume.Name = slug.Make(volume.Name)
	}

	for _, network := range dst.Networks {
		network.Name = fmt.Sprintf("%v_%s", ns, network.Name)
		network.Name = slug.Make(network.Name)
	}
}
