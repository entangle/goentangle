package goentangle

import (
	"testing"
)

func TestExceptionDefinition(t *testing.T) {
	// Make a definition.
	def := NewExceptionDefinition("test", "NameException")

	// Make an exception.
	var err Exception = def.New("Description")
	if err.Definition() != "test" {
		t.Errorf("invalid definition: %s", err.Definition())
	} else if err.Name() != "NameException" {
		t.Errorf("invalid name: %s", err.Name())
	} else if err.Error() != "Description" {
		t.Errorf("invalid exception message: %s", err.Error())
	}

	// Make a formated exception.
	err = def.Newf("Description: %d, %s", 123, "something")
	if err.Definition() != "test" {
		t.Errorf("invalid definition: %s", err.Definition())
	} else if err.Name() != "NameException" {
		t.Errorf("invalid name: %s", err.Name())
	} else if err.Error() != "Description: 123, something" {
		t.Errorf("invalid exception message: %s", err.Error())
	}
}
