package yaml

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestCommand(t *testing.T) {
	var tests = []struct {
		yaml string
		want []string
	}{
		{
			yaml: "command: echo hello world",
			want: []string{"echo", "hello", "world"},
		},
		{
			yaml: "command: echo 'hello world'",
			want: []string{"echo", "hello world"},
		},
		{
			yaml: "command: [ echo, 'hello world' ]",
			want: []string{"echo", "hello world"},
		},
		{
			yaml: "{ entrypoint: [/bin/sh, -c], command: ['echo hello world'] }",
			want: []string{"echo hello world"},
		},
	}

	for _, test := range tests {
		got := struct {
			Command Command
		}{}
		if err := yaml.Unmarshal([]byte(test.yaml), &got); err != nil {
			t.Errorf("got error unmarshaling %q", test.yaml)
		}
		if !reflect.DeepEqual([]string(got.Command), test.want) {
			t.Errorf("got command %v want %v", got.Command, test.want)
		}
	}
}
