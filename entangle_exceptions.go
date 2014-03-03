package goentangle

var (
	// Bad message error definition.
	BadMessageError = NewExceptionDefinition("entangle", "BadMessage")

	// Internal server error.
	InternalServerError = NewExceptionDefinition("entangle", "InternalServerError")

	// Unknown method.
	UnknownMethodError = NewExceptionDefinition("entangle", "UnknownMethod")

	// Unknown exception.
	UnknownExceptionError = NewExceptionDefinition("entangle", "UnknownException")
)
