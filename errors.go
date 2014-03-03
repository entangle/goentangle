package goentangle

import (
	"fmt"
)

// Entangle error.
type Error interface {
	error

	// Namespace.
	Namespace() string

	// Name.
	Name() string
}

// Entangle error definition.
//
// Error definition that can produce Entangle of a specific type.
type ErrorDefinition interface {
	// New error.
	New(description string) Error

	// New formatted error.
	//
	// Behaves like fmt.Printf.
	Newf(format string, a ...interface{}) Error
}

// Entangle error implementation.
type entangleError struct {
	namespace   string
	name        string
	description string
}

func (e *entangleError) Namespace() string {
	return e.namespace
}

func (e *entangleError) Name() string {
	return e.name
}

func (e *entangleError) Error() string {
	return e.description
}

// Entangle error definition implementation.
type entangleErrorDefinition struct {
	namespace string
	name      string
}

func (d *entangleErrorDefinition) New(description string) Error {
	return &entangleError{
		d.namespace,
		d.name,
		description,
	}
}

func (d *entangleErrorDefinition) Newf(format string, a ...interface{}) Error {
	return &entangleError{
		d.namespace,
		d.name,
		fmt.Sprintf(format, a...),
	}
}

// New error definition.
func NewErrorDefinition(namespace, name string) ErrorDefinition {
	return &entangleErrorDefinition{
		namespace,
		name,
	}
}
