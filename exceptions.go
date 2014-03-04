package goentangle

import (
	"fmt"
)

// Entangle exception.
type Exception interface {
	error

	// Definition.
	Definition() string

	// Name.
	Name() string
}

// Entangle exception definition.
//
// Exception definition that can produce an exception of a specific type.
type ExceptionDefinition interface {
	// New error.
	New(description string) Exception

	// New formatted error.
	//
	// Behaves like fmt.Printf.
	Newf(format string, a ...interface{}) Exception
}

// Entangle exception implementation.
type entangleException struct {
	definition  string
	name        string
	description string
}

func (e *entangleException) Definition() string {
	return e.definition
}

func (e *entangleException) Name() string {
	return e.name
}

func (e *entangleException) Error() string {
	return e.description
}

// Entangle exception definition implementation.
type entangleExceptionDefinition struct {
	definition string
	name       string
}

func (d *entangleExceptionDefinition) New(description string) Exception {
	return &entangleException{
		d.definition,
		d.name,
		description,
	}
}

func (d *entangleExceptionDefinition) Newf(format string, a ...interface{}) Exception {
	return &entangleException{
		d.definition,
		d.name,
		fmt.Sprintf(format, a...),
	}
}

// New error definition.
func NewExceptionDefinition(definition, name string) ExceptionDefinition {
	return &entangleExceptionDefinition{
		definition,
		name,
	}
}
