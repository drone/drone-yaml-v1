package yaml

import (
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestExternal(t *testing.T) {
	var tests = []struct {
		yaml string
		want External
	}{
		{
			yaml: "external: false",
			want: External{External: false},
		},
		{
			yaml: "external: true",
			want: External{External: true},
		},
		{
			yaml: "external: { name: redis_secret }",
			want: External{External: true, Name: "redis_secret"},
		},
	}

	for _, test := range tests {
		got := struct {
			External External
		}{}
		if err := yaml.Unmarshal([]byte(test.yaml), &got); err != nil {
			t.Errorf("got error unmarshaling %q", test.yaml)
		}
		if !reflect.DeepEqual(got.External, test.want) {
			t.Errorf("got external %v want %v", got.External, test.want)
		}
	}
}
