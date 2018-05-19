package yaml

import (
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestConstraintMatch(t *testing.T) {
	testdata := []struct {
		conf string
		with string
		want bool
	}{
		// string value
		{
			conf: "master",
			with: "develop",
			want: false,
		},
		{
			conf: "master",
			with: "master",
			want: true,
		},
		{
			conf: "feature/*",
			with: "feature/foo",
			want: true,
		},
		// slice value
		{
			conf: "[ master, feature/* ]",
			with: "develop",
			want: false,
		},
		{
			conf: "[ master, feature/* ]",
			with: "master",
			want: true,
		},
		{
			conf: "[ master, feature/* ]",
			with: "feature/foo",
			want: true,
		},
		// includes block
		{
			conf: "include: master",
			with: "develop",
			want: false,
		},
		{
			conf: "include: master",
			with: "master",
			want: true,
		},
		{
			conf: "include: feature/*",
			with: "master",
			want: false,
		},
		{
			conf: "include: feature/*",
			with: "feature/foo",
			want: true,
		},
		{
			conf: "include: [ master, feature/* ]",
			with: "develop",
			want: false,
		},
		{
			conf: "include: [ master, feature/* ]",
			with: "master",
			want: true,
		},
		{
			conf: "include: [ master, feature/* ]",
			with: "feature/foo",
			want: true,
		},
		// excludes block
		{
			conf: "exclude: master",
			with: "develop",
			want: true,
		},
		{
			conf: "exclude: master",
			with: "master",
			want: false,
		},
		{
			conf: "exclude: feature/*",
			with: "master",
			want: true,
		},
		{
			conf: "exclude: feature/*",
			with: "feature/foo",
			want: false,
		},
		{
			conf: "exclude: [ master, develop ]",
			with: "master",
			want: false,
		},
		{
			conf: "exclude: [ feature/*, bar ]",
			with: "master",
			want: true,
		},
		{
			conf: "exclude: [ feature/*, bar ]",
			with: "feature/foo",
			want: false,
		},
		// include and exclude blocks
		{
			conf: "{ include: [ master, feature/* ], exclude: [ develop ] }",
			with: "master",
			want: true,
		},
		{
			conf: "{ include: [ master, feature/* ], exclude: [ feature/bar ] }",
			with: "feature/bar",
			want: false,
		},
		{
			conf: "{ include: [ master, feature/* ], exclude: [ master, develop ] }",
			with: "master",
			want: false,
		},
		// empty blocks
		{
			conf: "",
			with: "master",
			want: true,
		},
		// double star
		{
			conf: "foo/**",
			with: "foo/bar/baz/qux",
			want: true,
		},
		{
			conf: "foo/**/qux",
			with: "foo/bar/baz/qux",
			want: true,
		},
	}
	for _, test := range testdata {
		c := parseConstraint(test.conf)
		got, want := c.Match(test.with), test.want
		if got != want {
			t.Errorf("Expect %q matches %q is %v", test.with, test.conf, want)
		}
	}
}

func TestConstraintMatchAny(t *testing.T) {
	testdata := []struct {
		conf string
		with []string
		want bool
	}{
		{
			conf: "foo/bar",
			with: []string{"foo/bar"},
			want: true,
		},
		{
			conf: "foo/*",
			with: []string{"foo/bar"},
			want: true,
		},
		{
			conf: "foo/*",
			with: []string{"foo/baz", "/foo/bar"},
			want: true,
		},
		{
			conf: "foo/**",
			with: []string{"foo/bar/baz/qux"},
			want: true,
		},
		{
			conf: "foo/**",
			with: []string{"bar/baz/qux/foo"},
			want: false,
		},
		{
			conf: "",
			with: []string{},
			want: true,
		},
	}
	for _, test := range testdata {
		c := parseConstraint(test.conf)
		got, want := c.MatchAny(test.with), test.want
		if got != want {
			t.Errorf("Expect %+v matches %q is %v", test.with, test.conf, want)
		}
	}
}

func TestConstraintMap(t *testing.T) {
	testdata := []struct {
		conf string
		with map[string]string
		want bool
	}{
		{
			conf: "GOLANG: 1.7",
			with: map[string]string{"GOLANG": "1.7"},
			want: true,
		},
		{
			conf: "GOLANG: tip",
			with: map[string]string{"GOLANG": "1.7"},
			want: false,
		},
		{
			conf: "{ GOLANG: 1.7, REDIS: 3.1 }",
			with: map[string]string{"GOLANG": "1.7", "REDIS": "3.1", "MYSQL": "5.6"},
			want: true,
		},
		{
			conf: "{ GOLANG: 1.7, REDIS: 3.1 }",
			with: map[string]string{"GOLANG": "1.7", "REDIS": "3.0"},
			want: false,
		},
		{
			conf: "{ GOLANG: 1.7, REDIS: 3.* }",
			with: map[string]string{"GOLANG": "1.7", "REDIS": "3.0"},
			want: false,
		},
		// include syntax
		{
			conf: "include: { GOLANG: 1.7 }",
			with: map[string]string{"GOLANG": "1.7"},
			want: true,
		},
		{
			conf: "include: { GOLANG: tip }",
			with: map[string]string{"GOLANG": "1.7"},
			want: false,
		},
		{
			conf: "include: { GOLANG: 1.7, REDIS: 3.1 }",
			with: map[string]string{"GOLANG": "1.7", "REDIS": "3.1", "MYSQL": "5.6"},
			want: true,
		},
		{
			conf: "include: { GOLANG: 1.7, REDIS: 3.1 }",
			with: map[string]string{"GOLANG": "1.7", "REDIS": "3.0"},
			want: false,
		},
		// exclude syntax
		{
			conf: "exclude: { GOLANG: 1.7 }",
			with: map[string]string{"GOLANG": "1.7"},
			want: false,
		},
		{
			conf: "exclude: { GOLANG: tip }",
			with: map[string]string{"GOLANG": "1.7"},
			want: true,
		},
		{
			conf: "exclude: { GOLANG: 1.7, REDIS: 3.1 }",
			with: map[string]string{"GOLANG": "1.7", "REDIS": "3.1", "MYSQL": "5.6"},
			want: false,
		},
		{
			conf: "exclude: { GOLANG: 1.7, REDIS: 3.1 }",
			with: map[string]string{"GOLANG": "1.7", "REDIS": "3.0"},
			want: true,
		},
		// exclude AND include values
		{
			conf: "{ include: { GOLANG: 1.7 }, exclude: { GOLANG: 1.7 } }",
			with: map[string]string{"GOLANG": "1.7"},
			want: false,
		},
		// blanks
		{
			conf: "",
			with: map[string]string{"GOLANG": "1.7", "REDIS": "3.0"},
			want: true,
		},
		{
			conf: "GOLANG: 1.7",
			with: map[string]string{},
			want: false,
		},
		{
			conf: "{ GOLANG: 1.7, REDIS: 3.0 }",
			with: map[string]string{},
			want: false,
		},
		{
			conf: "include: { GOLANG: 1.7, REDIS: 3.1 }",
			with: map[string]string{},
			want: false,
		},
		{
			conf: "exclude: { GOLANG: 1.7, REDIS: 3.1 }",
			with: map[string]string{},
			want: true,
		},
	}
	for _, test := range testdata {
		c := parseConstraintMap(test.conf)
		got, want := c.Match(test.with), test.want
		if got != want {
			t.Errorf("Expect %q matches %q is %v", test.with, test.conf, want)
		}
	}
}

func parseConstraint(s string) *Constraint {
	c := &Constraint{}
	yaml.Unmarshal([]byte(s), c)
	return c
}

func parseConstraintMap(s string) *ConstraintMap {
	c := &ConstraintMap{}
	yaml.Unmarshal([]byte(s), c)
	return c
}
