package errors

import (
	"fmt"
)

/**
 * An error
 */
type Error struct {
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`
	Cause   error       `json:"-"`
}

/**
 * Create a status error
 */
func Errorf(f string, a ...interface{}) *Error {
	return &Error{fmt.Sprintf(f, a...), nil, nil}
}

/**
 * Set detail
 */
func (e *Error) SetDetail(d interface{}) *Error {
	e.Detail = d
	return e
}

/**
 * Set the underlying cause
 */
func (e *Error) SetCause(c error) *Error {
	e.Cause = c
	return e
}

/**
 * Obtain the error message
 */
func (e Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%v: %v", e.Message, e.Cause.Error())
	} else {
		return e.Message
	}
}

/**
 * A set of errors
 */
type Set []error

/**
 * Create a set of errors. Only non-nil parameters are included. If only
 * one non-nil parameter is provided it is simply returned and a set is
 * not actually created.
 */
func NewSet(e ...error) error {
	s := make(Set, 0)
	for _, v := range e {
		if v != nil {
			s = append(s, v)
		}
	}
	if len(s) == 1 {
		return s[0]
	} else {
		return s
	}
}

/**
 * Obtain the error message
 */
func (e Set) Error() string {
	var s string
	for i, v := range e {
		if i > 0 {
			s += "; "
		}
		s += v.Error()
	}
	return s
}
