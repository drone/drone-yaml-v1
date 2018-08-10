package compiler

import (
	"reflect"
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

func Test_splitVolumeParts(t *testing.T) {
	testdata := []struct {
		from    string
		to      []string
		success bool
	}{
		{
			from:    `Z::Z::rw`,
			to:      []string{`Z:`, `Z:`, `rw`},
			success: true,
		},
		{
			from:    `Z:\:Z:\:rw`,
			to:      []string{`Z:\`, `Z:\`, `rw`},
			success: true,
		},
		{
			from:    `Z:\git\refs:Z:\git\refs:rw`,
			to:      []string{`Z:\git\refs`, `Z:\git\refs`, `rw`},
			success: true,
		},
		{
			from:    `Z:\git\refs:Z:\git\refs`,
			to:      []string{`Z:\git\refs`, `Z:\git\refs`},
			success: true,
		},
		{
			from:    `Z:/:Z:/:rw`,
			to:      []string{`Z:/`, `Z:/`, `rw`},
			success: true,
		},
		{
			from:    `Z:/git/refs:Z:/git/refs:rw`,
			to:      []string{`Z:/git/refs`, `Z:/git/refs`, `rw`},
			success: true,
		},
		{
			from:    `Z:/git/refs:Z:/git/refs`,
			to:      []string{`Z:/git/refs`, `Z:/git/refs`},
			success: true,
		},
		{
			from:    `/test:/test`,
			to:      []string{`/test`, `/test`},
			success: true,
		},
		{
			from:    `test:/test`,
			to:      []string{`test`, `/test`},
			success: true,
		},
		{
			from:    `test:test`,
			to:      []string{`test`, `test`},
			success: true,
		},
	}
	for _, test := range testdata {
		results := splitVolumeParts(test.from)

		if reflect.DeepEqual(results, test.to) != test.success {
			t.Errorf("Expect %q matches %q is %v", test.from, results, test.to)
		}
	}
}
