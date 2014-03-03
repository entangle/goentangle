package goentangle

import (
	"testing"
)

func TestErrorDefinition(t *testing.T) {
	// Make a definition.
	def := NewErrorDefinition("test", "NameError")

	// Make an error.
	var err Error = def.New("Description")
	if err.Namespace() != "test" {
		t.Errorf("invalid namespace: %s", err.Namespace())
	} else if err.Name() != "NameError" {
		t.Errorf("invalid name: %s", err.Name())
	} else if err.Error() != "Description" {
		t.Errorf("invalid error message: %s", err.Error())
	}

	// Make a formated error.
	err = def.Newf("Description: %d, %s", 123, "something")
	if err.Namespace() != "test" {
		t.Errorf("invalid namespace: %s", err.Namespace())
	} else if err.Name() != "NameError" {
		t.Errorf("invalid name: %s", err.Name())
	} else if err.Error() != "Description: 123, something" {
		t.Errorf("invalid error message: %s", err.Error())
	}
}
