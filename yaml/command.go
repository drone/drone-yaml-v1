package yaml

import (
	"github.com/flynn/go-shlex"
)

// Command represents a docker command, can be a string or an array of strings.
type Command []string

// UnmarshalYAML implements the Unmarshaller interface.
func (s *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var stringType string
	if err := unmarshal(&stringType); err == nil {
		parts, err := shlex.Split(stringType)
		if err != nil {
			return err
		}
		*s = parts
		return nil
	}

	var sliceType []string
	if err := unmarshal(&sliceType); err != nil {
		return err
	}
	*s = sliceType
	return nil
}
