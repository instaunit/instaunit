package testcase

import (
	yaml "gopkg.in/yaml.v3"
)

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

type Authentication struct {
	Type        string `yaml:"method"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Scheme      string `yaml:"scheme"`
	Format      string `yaml:"format"`
	Location    string `yaml:"location"`
}

type AccessControl struct {
	Scopes []string `yaml:"scopes"`
	Roles  []string `yaml:"roles"`
}
