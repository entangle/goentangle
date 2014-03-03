package goentangle

var (
	// Bad message error definition.
	BadMessageError = NewErrorDefinition("entangle", "BadMessage")

	// Internal server error.
	InternalServerError = NewErrorDefinition("entangle", "InternalServerError")

	// Unknown method.
	UnknownMethodError = NewErrorDefinition("entangle", "UnknownMethod")

	// Unknown exception.
	UnknownExceptionError = NewErrorDefinition("entangle", "UnknownException")
)
