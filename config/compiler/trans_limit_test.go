package compiler

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/yaml"
)

func Test_transformLimits(t *testing.T) {
	src := new(yaml.Container)
	dst := new(engine.Step)

	limits := Resources{
		MemLimit:     1,
		MemSwapLimit: 2,
		ShmSize:      3,
		CPUQuota:     4,
		CPUSet:       "1,3",
	}

	transformLimits(limits)(dst, src, nil)

	if got, want := dst.MemLimit, limits.MemLimit; got != want {
		t.Errorf("Want MemLimit %v, got %v", want, got)
	}
	if got, want := dst.MemSwapLimit, limits.MemSwapLimit; got != want {
		t.Errorf("Want MemSwapLimit %v, got %v", want, got)
	}
	if got, want := dst.ShmSize, limits.ShmSize; got != want {
		t.Errorf("Want ShmSize %v, got %v", want, got)
	}
	if got, want := dst.CPUQuota, limits.CPUQuota; got != want {
		t.Errorf("Want CPUQuota %v, got %v", want, got)
	}
	if got, want := dst.CPUSet, limits.CPUSet; got != want {
		t.Errorf("Want CPUSet %v, got %v", want, got)
	}
}
