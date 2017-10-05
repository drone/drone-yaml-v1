package yaml

import (
	"strconv"
)

// StringInt represents a string or an integer.
type StringInt int64

// UnmarshalYAML implements the Unmarshaller interface.
func (s *StringInt) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var intType int64
	if err := unmarshal(&intType); err == nil {
		*s = StringInt(intType)
		return nil
	}

	var stringType string
	if err := unmarshal(&stringType); err != nil {
		return err
	}

	intType, err := strconv.ParseInt(stringType, 10, 64)
	if err == nil {
		*s = StringInt(intType)
	}
	return err
}
