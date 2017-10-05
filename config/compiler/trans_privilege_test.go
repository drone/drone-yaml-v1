package compiler

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/yaml"
)

func Test_transformPrivilege(t *testing.T) {
	testdatum := []struct {
		image   string
		match   string
		command []string
		matched bool
	}{
		// should never match
		{
			image: "golang",
			match: "plugin/docker",
		},
		{
			image: "golang:latest",
			match: "plugin/docker",
		},
		{
			image: "docker.io/golang:latest",
			match: "plugin/docker",
		},
		{
			image: "docker.io/golang:latest",
			match: "plugin/docker:latest",
		},

		// should match
		{
			image:   "plugins/docker",
			match:   "plugins/docker",
			matched: true,
		},
		{
			image:   "plugins/docker:latest",
			match:   "plugins/docker",
			matched: true,
		},
		{
			image:   "plugins/docker",
			match:   "plugins/docker:latest",
			matched: true,
		},
		{
			image:   "docker.io/plugins/docker:latest",
			match:   "plugins/docker:latest",
			matched: true,
		},

		// should not match when commands set
		{
			image:   "plugins/heroku",
			match:   "plugins/heroku",
			matched: false,
			command: []string{"go build"},
		},
	}
	for _, testdata := range testdatum {
		src := new(yaml.Container)
		src.Commands = testdata.command

		dst := new(engine.Step)
		dst.Image = testdata.image

		transformPrivilege(testdata.match)(dst, src, nil)

		if got, want := dst.Privileged, testdata.matched; got != want {
			t.Errorf("Want image %q matches %q is %v, got %v",
				testdata.image,
				testdata.match,
				want,
				got,
			)
		}
	}
}
