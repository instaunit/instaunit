package test

import yaml "gopkg.in/yaml.v3"

type parameter struct {
	Location    string `yaml:"location"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
}

type Parameter struct {
	parameter
}

func (p *Parameter) UnmarshalYAML(value *yaml.Node) error {
	if l := len(value.Content); l == 0 {
		*p = Parameter{parameter{Description: value.Value}}
		return nil
	}
	var x parameter
	err := value.Decode(&x)
	if err != nil {
		return err
	}
	*p = Parameter{x}
	return nil
}
