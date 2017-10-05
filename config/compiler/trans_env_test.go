package compiler

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/yaml"
)

func Test_transformEnv(t *testing.T) {
	src := new(yaml.Container)
	dst := new(engine.Step)

	env := map[string]string{"foo": "bar", "baz": ""}
	transformEnv(env)(dst, src, nil)

	if got, want := dst.Environment["foo"], "bar"; got != want {
		t.Errorf("Want environment variable foo=%q, got %q", want, got)
	}
	if _, ok := dst.Environment["baz"]; ok {
		t.Errorf("Should not inject empty environment variables")
	}
}
