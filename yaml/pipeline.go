package yaml

// Pipeline represents the pipeline section of the yaml.
type Pipeline struct {
	Name  string
	Steps []*Container
}

// UnmarshalYAML implements the Unmarshaller interface.
func (c *Pipeline) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s1 := struct {
		Name  string
		Steps []*Container
	}{}
	err := unmarshal(&s1)
	if err == nil && len(s1.Steps) != 0 {
		c.Name = s1.Name
		c.Steps = s1.Steps
		return nil
	}

	s2 := struct {
		Containers
	}{}

	err = unmarshal(&s2)
	if err != nil {
		return err
	}
	c.Steps = s2.Containers.Containers
	c.Name = "default"
	return nil
}
