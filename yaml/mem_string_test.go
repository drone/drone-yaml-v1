package yaml

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestStringMem(t *testing.T) {
	var tests = []struct {
		yaml string
		want int64
	}{
		{
			yaml: "1KB",
			want: 1024,
		},
		{
			yaml: "1024",
			want: 1024,
		},
	}

	for _, test := range tests {
		var got MemStringInt

		if err := yaml.Unmarshal([]byte(test.yaml), &got); err != nil {
			t.Error(err)
		}

		if test.want != int64(got) {
			t.Errorf("got int64 %v want %v", got, test.want)
		}
	}

	var got MemStringInt
	if err := yaml.Unmarshal([]byte("{}"), &got); err == nil {
		t.Errorf("Want error unmarshaling invalid memory value.")
	}
}
