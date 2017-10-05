package compiler

import "testing"

func Test_encodeSlice(t *testing.T) {
	testdatum := []struct {
		data interface{}
		text string
	}{
		{
			data: []string{"foo", "bar", "baz"},
			text: "foo,bar,baz",
		},
		{
			data: []int{1, 1, 2, 3, 5, 8},
			text: "1,1,2,3,5,8",
		},
		{
			data: []struct {
				Name string `json:"name"`
			}{
				{"jane"},
				{"john"},
			},
			text: `[{"name":"jane"},{"name":"john"}]`,
		},
	}

	for _, testdata := range testdatum {
		if got, want := encodeSlice(testdata.data), testdata.text; got != want {
			t.Errorf("Want interface{} encoded to %q, got %q", want, got)
		}
	}
}

func Test_encodeMap(t *testing.T) {
	testdatum := []struct {
		data interface{}
		text string
	}{
		{
			data: map[string]string{"foo": "bar"},
			text: `{"foo":"bar"}`,
		},
	}

	for _, testdata := range testdatum {
		if got, want := encodeMap(testdata.data), testdata.text; got != want {
			t.Errorf("Want interface{} encoded to %q, got %q", want, got)
		}
	}
}
