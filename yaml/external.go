package yaml

// External represent an external resource.
type External struct {
	External bool
	Name     string
}

// UnmarshalYAML unmarshals the extneral resource.
func (e *External) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var boolType bool
	if err := unmarshal(&boolType); err == nil {
		e.External = boolType
		return nil
	}

	var structType = struct {
		Name string
	}{}
	if err := unmarshal(&structType); err != nil {
		return err
	}
	if structType.Name != "" {
		e.External = true
		e.Name = structType.Name
	}
	return nil
}
