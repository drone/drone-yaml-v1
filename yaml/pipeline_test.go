package yaml

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestPipeline(t *testing.T) {
	p := new(Pipeline)
	err := yaml.Unmarshal(samplePipeline, p)
	if err != nil {
		t.Error(err)
	}
	if got, want := p.Name, "default"; got != want {
		t.Errorf("Want pipeline name %q, got %q", want, got)
	}
	if got, want := len(p.Steps), 2; got != want {
		t.Errorf("Want %d pipeline steps, got %d", want, got)
	}
	if got, want := p.Steps[0].Name, "build"; got != want {
		t.Errorf("Want pipeline name %q, got %q", want, got)
	}
	if got, want := p.Steps[1].Name, "test"; got != want {
		t.Errorf("Want pipeline name %q, got %q", want, got)
	}
}

func TestPipeline_Legacy(t *testing.T) {
	p := new(Pipeline)
	err := yaml.Unmarshal(samplePipeline, p)
	if err != nil {
		t.Error(err)
	}
	if got, want := p.Name, "default"; got != want {
		t.Errorf("Want pipeline name %q, got %q", want, got)
	}
	if got, want := len(p.Steps), 2; got != want {
		t.Errorf("Want %d pipeline steps, got %d", want, got)
	}
	if got, want := p.Steps[0].Name, "build"; got != want {
		t.Errorf("Want pipeline name %q, got %q", want, got)
	}
	if got, want := p.Steps[1].Name, "test"; got != want {
		t.Errorf("Want pipeline name %q, got %q", want, got)
	}
}

var samplePipeline = []byte(`
name: default
steps:
  - name: build
    commands:
      - go get
      - go build
  - name: test
    commands:
      - go lint
      - go test
`)

var legacyPipeline = []byte(`
build:
  commands:
    - go get
    - go build
test:
  commands:
    - go lint
    - go test
`)
