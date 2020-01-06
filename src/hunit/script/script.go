package script

import (
	"fmt"
	"os"
	"strings"

	"github.com/instaunit/instaunit/hunit/expr"

	"github.com/bww/epl"
	"github.com/robertkrimen/otto"
)

type ErrInvalidTypeError struct {
	source string
	expect string
	result interface{}
}

func (e ErrInvalidTypeError) Error() string {
	return fmt.Sprintf("Invalid type; expected %s in {%s}, got (%T) %v", e.expect, e.source, e.result, e.result)
}

// A script
type Script struct {
	Type   string `yaml:"type"`
	Source string `yaml:"source"`
}

func (s Script) Bool(v expr.Variables) (bool, error) {
	res, err := s.Eval(v)
	if err != nil {
		return false, err
	}
	switch v := res.(type) {
	case bool:
		return v, nil
	case otto.Value:
		return v.ToBoolean()
	default:
		return false, ErrInvalidTypeError{s.Source, "bool", res}
	}
}

func (s Script) Eval(v expr.Variables) (interface{}, error) {
	cxt := expr.RuntimeContext(v, os.Environ())
	switch strings.ToLower(s.Type) {
	case "epl", "":
		return s.evalEPL(cxt)
	case "js", "javascript":
		return s.evalJS(cxt)
	default:
		return false, fmt.Errorf("Unsupported script type: %v", s.Type)
	}
}

func (s Script) evalEPL(cxt expr.Variables) (interface{}, error) {
	prg, err := epl.Compile(s.Source)
	if err != nil {
		return false, err
	}
	res, err := prg.Exec(cxt)
	if err != nil {
		return false, fmt.Errorf("Could not evaluate expression: {%v}: %v", s.Source, err)
	}
	return res, nil
}

func (s Script) evalJS(cxt expr.Variables) (interface{}, error) {
	vm := otto.New()
	for k, v := range cxt {
		vm.Set(k, v)
	}
	res, err := vm.Run(s.Source)
	if err != nil {
		return nil, err
	}
	return res, nil
}
