package yaml

import (
	"reflect"
	"testing"

	"github.com/kr/pretty"
	"gopkg.in/yaml.v2"
)

func TestUnmarshalReport(t *testing.T) {
	testdata := []struct {
		from string
		want Reports
	}{
		{
			from: "{ coverage: path/to/coverage.out }",
			want: Reports{
				Coverage: &Report{
					Source: "path/to/coverage.out",
				},
			},
		},
		{
			from: "{ coverage: { source: path/to/coverage.out, format: gocov } }",
			want: Reports{
				Coverage: &Report{
					Source: "path/to/coverage.out",
					Format: "gocov",
				},
			},
		},
	}

	for _, test := range testdata {
		in := []byte(test.from)
		got := Reports{}
		err := yaml.Unmarshal(in, &got)
		if err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(test.want, got) {
			t.Errorf("problem parsing report %q", test.from)
			pretty.Ldiff(t, test.want, got)
		}
	}
}
