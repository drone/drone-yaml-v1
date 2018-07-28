package linter

import (
	"testing"

	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func TestIsService(t *testing.T) {
	conf := new(config.Config)
	conf.Services = map[string]*yaml.Container{}
	container := new(yaml.Container)

	if got, want := IsService(conf, container), false; got != want {
		t.Errorf("Expect pipeline contianer not classified as service")
	}

	container.Detached = true
	if got, want := IsService(conf, container), true; got != want {
		t.Errorf("Expect detached contianer classified as service")
	}

	conf.Services["mysql"] = container
	container.Detached = false
	if got, want := IsService(conf, container), true; got != want {
		t.Errorf("Expect service contianer classified as service")
	}
}

func TestIsDataVolume(t *testing.T) {
	conf := new(config.Config)
	volume := new(yaml.Volume)
	volume.Source = "/foo/bar"

	if got, want := IsDataVolume(conf, volume), false; got != want {
		t.Errorf("Expect volume not classified as a data volume")
	}

	volume.Source = "global"
	if got, want := IsDataVolume(conf, volume), false; got != want {
		t.Errorf("Expect volume not classified as a data volume if not defined in volumes section")
	}

	conf.Volumes = map[string]config.Volume{}
	conf.Volumes["global"] = config.Volume{}
	volume.Source = "global"
	if got, want := IsDataVolume(conf, volume), true; got != want {
		t.Errorf("Expect volume classified as a data volume")
	}
}
