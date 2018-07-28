package config

import (
	"testing"

	"github.com/drone/drone-yaml-v1/yaml"
	"github.com/kr/pretty"
)

func TestParse(t *testing.T) {
	got, err := ParseString(sampleYaml)
	if err != nil {
		t.Error(err)
	}
	want := &Config{
		Version: 1,
		Platform: Platform{
			Name: "linux/amd64",
		},
		Workspace: Workspace{
			Base: "/go",
			Path: "src/github.com/octocat/hello-world",
		},
		Clone: Clone{
			Depth: 50,
		},
		Networks: map[string]Network{
			"custom": {Driver: "overlay"},
		},
		Volumes: map[string]Volume{
			"custom": {Driver: "blockbridge"},
		},
		Secrets: map[string]Secret{
			"password": {
				External: yaml.External{
					External: true,
				},
			},
		},
		Labels: yaml.SliceMap{
			map[string]string{
				"com.example.type": "build",
				"com.example.team": "frontend",
			},
		},
		Services: map[string]*yaml.Container{
			"database": &yaml.Container{Image: "mysql"},
		},
		DependsOn: yaml.StringSlice{"frontend", "backend"},
		Pipeline: []map[string]*yaml.Container{
			map[string]*yaml.Container{
				"test": &yaml.Container{
					Image:    "golang",
					Commands: []string{"go install", "go test"},
				},
			},
			map[string]*yaml.Container{
				"build": &yaml.Container{
					Image:    "golang",
					Commands: []string{"go build"},
				},
			},
			map[string]*yaml.Container{
				"slack": &yaml.Container{
					Image: "plugins/slack",
					Vargs: map[string]interface{}{"channel": "dev"},
				},
				"gitter": &yaml.Container{
					Image: "plugins/gitter",
				},
			},
		},
	}
	diff := pretty.Diff(got, want)
	if len(diff) != 0 {
		t.Errorf("Failed to parse yaml with anchors. Diff %s", diff)
	}
}

var sampleYaml = `
version: 1

platform:
  name: linux/amd64

workspace:
  path: src/github.com/octocat/hello-world
  base: /go

clone:
  depth: 50

pipeline:
  - test:
      image: golang
      commands:
      - go install
      - go test
  - build:
      image: golang
      commands:
      - go build
  - slack:
      image: plugins/slack
      channel: dev
    gitter:
      image: plugins/gitter

services:
  database:
    image: mysql

networks:
  custom:
    driver: overlay

volumes:
  custom:
    driver: blockbridge

labels:
  com.example.type: "build"
  com.example.team: "frontend"

secrets:
  password:
    external: true

depends_on:
- frontend
- backend
`

//
// the purpose behind this anchor test is to ensure we are using
// a patched version of go-yaml
//

func TestParseAnchor(t *testing.T) {
	got, err := ParseString(sampleYamlAnchors)
	if err != nil {
		t.Error(err)
	}
	want := &Config{
		Pipeline: []map[string]*yaml.Container{
			map[string]*yaml.Container{
				"notify_fail": &yaml.Container{
					Image: "plugins/slack",
				},
			},
			map[string]*yaml.Container{
				"notify_success": &yaml.Container{
					Image: "plugins/slack",
					Constraints: yaml.Constraints{
						Status: yaml.Constraint{
							Include: []string{"success"},
						},
					},
				},
			},
		},
	}

	diff := pretty.Diff(got, want)
	if len(diff) != 0 {
		t.Errorf("Failed to parse yaml with anchors. Diff %s", diff)
	}
}

var sampleYamlAnchors = `
_slack: &SLACK
  image: plugins/slack
pipeline:
  - notify_fail: *SLACK
  - notify_success:
      << : *SLACK
      when:
        status: success
`

//
// the purpose behind this anchor test is to ensure we are using
// a patched version of go-yaml
//

func TestParseMulti(t *testing.T) {
	got, err := ParseMultiString(sampleYamlMulti)
	if err != nil {
		t.Error(err)
	}
	want := []*Config{
		&Config{
			Metadata: Metadata{
				Name: "backend",
			},
			Platform: Platform{
				Name: "linux/amd64",
			},

			Pipeline: []map[string]*yaml.Container{
				map[string]*yaml.Container{
					"build": &yaml.Container{
						Commands: []string{"go get", "go build"},
						Image:    "golang",
					},
				},
				map[string]*yaml.Container{
					"test": &yaml.Container{
						Commands: []string{"go test", "go lint"},
						Image:    "golang",
					},
				},
			},
		},
		&Config{
			Metadata: Metadata{
				Name: "frontend",
			},
			Platform: Platform{
				Name: "linux/arm",
			},
			Pipeline: []map[string]*yaml.Container{
				map[string]*yaml.Container{
					"test": &yaml.Container{
						Commands: []string{"npm install", "npm test"},
						Image:    "node",
					},
				},
			},
		},
	}
	if diff := pretty.Diff(got, want); len(diff) != 0 {
		t.Errorf("Failed to parse multi-yaml document. Diff %s", diff)
	}
}

var sampleYamlMulti = `---
metadata:
  name: backend
platform:
  name: linux/amd64
pipeline:
  - build:
      image: golang
      commands:
        - go get
        - go build
  - test:
      image: golang
      commands:
        - go test
        - go lint
---
metadata:
	name: frontend

platform:
	name: linux/arm

pipeline:
  - test:
      image: node
      commands:
        - npm install
        - npm test
`
