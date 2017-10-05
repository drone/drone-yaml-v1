package yaml

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestStringInt(t *testing.T) {
	var tests = []struct {
		yaml string
		want int64
	}{
		{
			yaml: "123",
			want: 123,
		},
		{
			yaml: "'123'",
			want: 123,
		},
	}

	for _, test := range tests {
		var got StringInt

		if err := yaml.Unmarshal([]byte(test.yaml), &got); err != nil {
			t.Error(err)
		}

		if test.want != int64(got) {
			t.Errorf("got int64 %v want %v", got, test.want)
		}
	}
}

func TestStringIntError(t *testing.T) {
	var tests = []struct {
		yaml string
		want string
	}{
		{
			yaml: "hello world",
			want: "strconv.ParseInt: parsing \"hello world\": invalid syntax",
		},
		{
			yaml: "'hello world'",
			want: "strconv.ParseInt: parsing \"hello world\": invalid syntax",
		},
		{
			yaml: "{}",
			want: "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!map into string",
		},
	}

	for _, test := range tests {
		var got StringInt

		err := yaml.Unmarshal([]byte(test.yaml), &got)
		if err == nil {
			t.Errorf("Want error unmarshaling integer %q", test.yaml)
		}
		if err.Error() != test.want {
			t.Errorf("Want error %q, got %q", test.want, err)
		}
	}
}
