package matrix

import (
	"testing"
)

func TestMatrix(t *testing.T) {
	axis, err := ParseString(fakeMatrix)
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := len(axis), 24; got != want {
		t.Errorf("Got %d matrix permutations, want %d", got, want)
	}

	set := map[string]bool{}
	for _, perm := range axis {
		set[perm.String()] = true
	}
	if got, want := len(axis), 24; got != want {
		t.Errorf("Got %d unique matrix permutations, want %d", got, want)
	}
}

func TestMatrixEmpty(t *testing.T) {
	axis, err := ParseString("")
	if err != nil {
		t.Error(err)
		return
	}
	if axis != nil {
		t.Errorf("Got non-nil matrix from empty string")
	}
}

func TestMatrixInclude(t *testing.T) {
	axis, err := ParseString(fakeMatrixInclude)
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := len(axis), 2; got != want {
		t.Errorf("Got %d matrix permutations, want %d", got, want)
	}
	if got, want := axis[0]["go_version"], "1.5"; got != want {
		t.Errorf("Got %s permutation, want %s", got, want)
	}
	if got, want := axis[1]["go_version"], "1.6"; got != want {
		t.Errorf("Got %s permutation, want %s", got, want)
	}
	if got, want := axis[0]["python_version"], "3.4"; got != want {
		t.Errorf("Got %s permutation, want %s", got, want)
	}
	if got, want := axis[1]["python_version"], "3.4"; got != want {
		t.Errorf("Got %s permutation, want %s", got, want)
	}
}

var fakeMatrix = `
matrix:
  go_version:
    - go1
    - go1.2
  python_version:
    - 3.2
    - 3.3
  django_version:
    - 1.7
    - 1.7.1
    - 1.7.2
  redis_version:
    - 2.6
    - 2.8
`

var fakeMatrixInclude = `
matrix:
  include:
    - go_version: 1.5
      python_version: 3.4
    - go_version: 1.6
      python_version: 3.4
`
