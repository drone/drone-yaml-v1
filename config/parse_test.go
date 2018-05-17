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
		Platform: "linux/amd64",
		Workspace: Workspace{
			Base: "/go",
			Path: "src/github.com/octocat/hello-world",
		},
		Clone: &yaml.Container{
			Image:    "docker:git",
			Commands: []string{"git clone https://github.com/octocat/hello-world.git"},
		},
		Branches: yaml.Constraint{
			Include: []string{"master"},
		},
		Networks: map[string]Network{
			"custom": {Driver: "overlay"},
		},
		Volumes: map[string]Volume{
			"custom": {Driver: "blockbridge"},
		},
		Secrets: map[string]Secret{
			"password": {
				Driver:     "custom",
				DriverOpts: map[string]interface{}{"custom.foo": "bar"},
			},
		},
		Labels: yaml.SliceMap{
			map[string]string{
				"com.example.type": "build",
				"com.example.team": "frontend",
			},
		},
		Services: yaml.Containers{
			Containers: []*yaml.Container{
				{
					Name:  "database",
					Image: "mysql",
				},
			},
		},
		Pipeline: yaml.Pipeline{
			Name: "default",
			Steps: []*yaml.Container{
				{
					Name:     "test",
					Image:    "golang",
					Commands: []string{"go install", "go test"},
				},
				{
					Name:     "build",
					Image:    "golang",
					Commands: []string{"go build"},
				},
				{
					Name:  "notify",
					Image: "plugins/slack",
					Vargs: map[string]interface{}{"channel": "dev"},
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
platform: linux/amd64
workspace:
  path: src/github.com/octocat/hello-world
  base: /go
clone:
  image: docker:git
  commands:
    - git clone https://github.com/octocat/hello-world.git
pipeline:
  test:
    image: golang
    commands:
      - go install
      - go test
  build:
    image: golang
    commands:
      - go build
  notify:
    image: plugins/slack
    channel: dev
services:
  database:
    image: mysql
branches: [ master ]
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
    driver: custom
    driver_opts:
      custom.foo: "bar"
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
		Pipeline: yaml.Pipeline{
			Name: "default",
			Steps: []*yaml.Container{
				{
					Name:  "notify_fail",
					Image: "plugins/slack",
				},
				{
					Name:  "notify_success",
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
  notify_fail: *SLACK
  notify_success:
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
			Platform: "linux/amd64",
			Pipeline: yaml.Pipeline{
				Name: "backend",
				Steps: []*yaml.Container{
					&yaml.Container{
						Commands: []string{"go get", "go build"},
						Image:    "golang",
						Name:     "build",
					},
					&yaml.Container{
						Commands: []string{"go test", "go lint"},
						Image:    "golang",
						Name:     "test",
					},
				},
			},
		},
		&Config{
			Platform: "linux/arm",
			Pipeline: yaml.Pipeline{
				Name: "frontend",
				Steps: []*yaml.Container{
					&yaml.Container{
						Commands: []string{"npm install", "npm test"},
						Image:    "node",
						Name:     "test",
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
platform: linux/amd64
pipeline:
  name: backend
  steps:
    - name: build
      image: golang
      commands:
        - go get
        - go build
    - name: test
      image: golang
      commands:
        - go test
        - go lint
---
platform: linux/arm
pipeline:
  name: frontend
  steps:
    - name: test
      image: node
      commands:
        - npm install
        - npm test
`
