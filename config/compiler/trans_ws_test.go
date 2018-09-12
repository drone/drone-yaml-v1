package compiler

import "testing"

func Test_normalizeWorkdirWindows(t *testing.T) {
	tests := []struct {
		before string
		after  string
	}{
		{
			before: `\gopath\src\github.com\drone\envsubst`,
			after:  `\gopath\src\github.com\drone\envsubst`,
		},
		{
			before: `c:\gopath\src\github.com\drone\envsubst`,
			after:  `\c\gopath\src\github.com\drone\envsubst`,
		},
		{
			before: `c:\gopath/src\github.com\drone\envsubst`,
			after:  `\c\gopath/src\github.com\drone\envsubst`,
		},
	}

	for _, test := range tests {
		after := normalizeWorkdirWindows(test.before)
		if got, want := after, test.after; got != want {
			t.Errorf("Want normalized path %s, got %s", want, got)
		}
	}
}

func TestToWindows(t *testing.T) {
	got := toWindowsDrive("/go/src/github.com/octocat/hello-world")
	want := "c:\\go\\src\\github.com\\octocat\\hello-world"
	if got != want {
		t.Errorf("Want windows drive %q, got %q", want, got)
	}
}
