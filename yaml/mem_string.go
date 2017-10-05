package yaml

import (
	units "github.com/docker/go-units"
)

// MemStringInt represents an integer or string followed by a unit symbol.
type MemStringInt int64

// UnmarshalYAML implements the Unmarshaller interface.
func (s *MemStringInt) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var intType int64
	if err := unmarshal(&intType); err == nil {
		*s = MemStringInt(intType)
		return nil
	}

	var stringType string
	if err := unmarshal(&stringType); err != nil {
		return err
	}

	intType, err := units.RAMInBytes(stringType)
	if err == nil {
		*s = MemStringInt(intType)
	}
	return err
}
