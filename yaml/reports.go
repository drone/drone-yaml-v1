package yaml

type (
	// Reports represents a list of report artifacts.
	Reports struct {
		Coverage *Report
	}

	// Report represents a coverage report artifacts.
	Report struct {
		Source string
		Format string
	}
)

// UnmarshalYAML implements the Unmarshaller interface.
func (r *Report) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err == nil {
		r.Source = str
		return nil
	}
	v := struct {
		Source string
		Format string
	}{}
	if err := unmarshal(&v); err != nil {
		return err
	}
	r.Source = v.Source
	r.Format = v.Format
	return nil
}
