package compiler

import (
	"regexp"
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func transformVolume(volume ...string) Transform {
	return func(dst *engine.Step, src *yaml.Container, conf *config.Config) {
		for _, mapping := range volume {
			var (
				name   string
				source string
				target string
			)
			var parts []string
			if conf.Platform.Name == "windows/amd64" {
				parts = splitVolumeParts(mapping)
			} else {
				parts = strings.Split(mapping, ":")
			}

			if len(parts) != 2 {
				continue
			}
			if strings.HasPrefix(parts[0], "/") {
				source = parts[0]
				target = parts[1]
			} else {
				name = parts[0]
				target = parts[1]
			}
			dst.Volumes = append(dst.Volumes, &engine.VolumeMapping{
				Name:   name,
				Source: source,
				Target: target,
			})
		}
	}
}

var volumeRE = regexp.MustCompile(`^((?:[\w]\:)?[^\:]*)\:((?:[\w]\:)?[^\:]*)(?:\:([rwom]*))?`)

func splitVolumeParts(volumeParts string) []string {
	if volumeRE.MatchString(volumeParts) {
		results := volumeRE.FindStringSubmatch(volumeParts)[1:]
		cleanResults := []string{}
		for _, item := range results {
			if item != "" {
				cleanResults = append(cleanResults, item)
			}
		}
		return cleanResults
	} else {
		return strings.Split(volumeParts, ":")
	}
}
