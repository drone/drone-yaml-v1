package compiler

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
)

const (
	// defaultTag defines the default tag used when performing images related actions and no tag or digest is specified
	defaultTag = "latest"
	// defaultHostname is the default built-in hostname
	defaultHostname = "docker.io"
	// legacyDefaultHostname is automatically converted to DefaultHostname
	legacyDefaultHostname = "index.docker.io"
	// defaultRepoPrefix is the prefix used for default repositories in default host
	defaultRepoPrefix = "library"
)

// trimImage returns the short image name without tag.
func trimImage(name string) string {
	name = strings.TrimPrefix(name, legacyDefaultHostname+"/")
	name = strings.TrimPrefix(name, defaultHostname+"/")
	name = strings.TrimPrefix(name, defaultRepoPrefix+"/")

	ref, err := reference.ParseNamed(name)
	if err != nil {
		return name
	}
	return reference.TrimNamed(ref).String()
}

// expandImage returns the fully qualified image name.
func expandImage(name string) string {
	ref, err := reference.ParseNamed(name)
	if err != nil {
		return name
	}
	s := ref.String()

	if !strings.Contains(s, ":") {
		s = fmt.Sprintf("%s:%s", s, defaultTag)
	}

	s = strings.TrimPrefix(s, legacyDefaultHostname+"/")
	s = strings.TrimPrefix(s, defaultRepoPrefix+"/")

	return s
}

// matchImage returns true if the image name matches
// an image in the list. Note the image tag is not used
// in the matching logic.
func matchImage(from string, to ...string) bool {
	from = trimImage(from)
	for _, match := range to {
		if from == trimImage(match) {
			return true
		}
	}
	return false
}

// matchHostname returns true if the image hostname
// matches the specified hostname.
func matchHostname(image, hostname string) bool {
	ref, err := reference.ParseNamed(image)
	if err != nil {
		return false
	}
	hn, _ := reference.SplitHostname(ref)

	if hn == "" || hn == defaultRepoPrefix {
		hn = defaultHostname
	}

	if hn == hostname {
		return true
	}

	return false
}
