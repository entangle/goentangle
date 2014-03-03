package goentangle

// Message.
type Message interface {
	// Message ID.
	MessageId() MessageId

	// Serialize the message.
	Serialize() []interface{}
}

// Request message.
type RequestMessage struct {
	// Message ID.
	messageId MessageId

	// Method.
	Method string

	// Arguments.
	Arguments []interface{}

	// Trace.
	Trace bool
}

func (m *RequestMessage) MessageId() MessageId {
	return m.messageId
}

func (m *RequestMessage) Serialize() []interface{} {
	return []interface{}{
		RequestOpcode,
		m.messageId,
		m.Method,
		m.Arguments,
		m.Trace,
	}
}

// Notification message.
type NotificationMessage struct {
	// Message ID.
	messageId MessageId

	// Method.
	Method string

	// Arguments.
	Arguments []interface{}
}

func (m *NotificationMessage) MessageId() MessageId {
	return m.messageId
}

func (m *NotificationMessage) Serialize() []interface{} {
	return []interface{}{
		NotificationOpcode,
		m.messageId,
		m.Method,
		m.Arguments,
	}
}

// Response message.
type ResponseMessage struct {
	// Message ID.
	messageId MessageId

	// Result.
	Result interface{}

	// Trace.
	Trace interface{}
}

func (m *ResponseMessage) MessageId() MessageId {
	return m.messageId
}

func (m *ResponseMessage) Serialize() []interface{} {
	return []interface{}{
		ResponseOpcode,
		m.messageId,
		m.Result,
		m.Trace,
	}
}

// Exception message.
type ExceptionMessage struct {
	// Message ID.
	messageId MessageId

	// Definition.
	Definition string

	// Name.
	Name string

	// Description.
	Description string

	// Trace.
	Trace interface{}
}

func (m *ExceptionMessage) MessageId() MessageId {
	return m.messageId
}

func (m *ExceptionMessage) Serialize() []interface{} {
	return []interface{}{
		ExceptionOpcode,
		m.messageId,
		m.Definition,
		m.Name,
		m.Description,
		m.Trace,
	}
}
