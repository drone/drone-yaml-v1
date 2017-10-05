package compiler

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/yaml"
)

func Test_transformRegistry(t *testing.T) {
	testdatum := []struct {
		image   string
		host    string
		matched bool
	}{
		{
			image:   "golang",
			host:    "docker.io",
			matched: true,
		},
		{
			image:   "docker.io/library/golang",
			host:    "docker.io",
			matched: true,
		},
		{
			image:   "grc.io/golang",
			host:    "docker.io",
			matched: false,
		},
		{
			image:   "golang",
			host:    "gcr.io",
			matched: false,
		},
	}

	for _, testdata := range testdatum {
		src := new(yaml.Container)
		dst := new(engine.Step)
		dst.Image = testdata.image

		registry := Registry{
			Username: "octocat",
			Password: "password",
			Hostname: testdata.host,
		}

		transformRegistry(registry)(dst, src, nil)

		if !testdata.matched {
			if dst.AuthConfig.Username != "" {
				t.Errorf("Expect host %q does not match image %q",
					testdata.host,
					testdata.image,
				)
			}
			continue
		}
		if got, want := dst.AuthConfig.Username, registry.Username; got != want {
			t.Errorf("Expect registry username %s, got %s", want, got)
		}
		if got, want := dst.AuthConfig.Password, registry.Password; got != want {
			t.Errorf("Expect registry password %s, got %s", want, got)
		}
	}
}
