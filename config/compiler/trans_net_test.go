package compiler

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/yaml"
)

func Test_transformNetwork(t *testing.T) {
	src := new(yaml.Container)
	dst := new(engine.Step)

	transformNetwork("foo", "bar")(dst, src, nil)

	if got, want := dst.Networks[0].Name, "foo"; got != want {
		t.Errorf("Got network name %q, want %q", got, want)
	}
	if got, want := dst.Networks[1].Name, "bar"; got != want {
		t.Errorf("Got network name %q, want %q", got, want)
	}
}
