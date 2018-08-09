package compiler

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func Test_transformVolume(t *testing.T) {
	testdatum := []struct {
		volume  string
		name    string
		source  string
		target  string
		skipped bool
	}{
		{
			volume: "/foo/bar:/baz",
			source: "/foo/bar",
			target: "/baz",
		},
		{
			volume: "default:/baz",
			name:   "default",
			target: "/baz",
		},
		{
			volume:  "default",
			skipped: true,
		},
	}

	for _, testdata := range testdatum {
		src := new(yaml.Container)
		dst := new(engine.Step)
		conf := new(config.Config)

		transformVolume(testdata.volume)(dst, src, conf)

		if testdata.skipped {
			if len(dst.Volumes) != 0 {
				t.Errorf("Expect volume %q skipped", testdata.volume)
			}
			continue
		}

		if got, want := dst.Volumes[0].Name, testdata.name; got != want {
			t.Errorf("Got volume name %q, want %q", got, want)
		}
		if got, want := dst.Volumes[0].Source, testdata.source; got != want {
			t.Errorf("Got volume source %q, want %q", got, want)
		}
		if got, want := dst.Volumes[0].Target, testdata.target; got != want {
			t.Errorf("Got volume target %q, want %q", got, want)
		}
	}
}
