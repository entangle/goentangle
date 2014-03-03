package goentangle

import (
	"errors"
	"testing"
)

func testConnReceiveFails(t *testing.T, data []byte, expectedErr error) {
	clientPipe, serverPipe := newTestingPipe()
	serverConn := NewConn(serverPipe, "test")

	if _, err := clientPipe.Write(data); err != nil {
		t.Fatalf("writing to client pipe failed unexpectedly: %v", err)
	}

	_, err := serverConn.Receive()
	if err != expectedErr {
		t.Errorf("Expected '%v' from Receive sending %v, but got '%v'", expectedErr, data, err)
	}

	clientPipe.Close()
	serverPipe.Close()
}

// Test low level receive error handling of Conn.
func TestConnReceiveLowLevel(t *testing.T) {
	// Write a piece of data to the client pipe, that is not an array and
	// make sure that receiving fails.
	testConnReceiveFails(t, []byte{0x00}, ErrInvalidMessageData)

	// Write an array with 0 or 1 elements to the client pipe, and make sure
	// that receiving fails.
	testConnReceiveFails(t, []byte{0x90}, ErrInvalidMessageData)
	testConnReceiveFails(t, []byte{0x91, 0xc0}, ErrInvalidMessageData)

	// Write an invalid opcode to the client pipe, and make sure that receiving
	// fails.
	for _, data := range [][]byte{
		// String.
		{0x92, 0xa0, 0x01},

		// Boolean.
		{0x92, 0xc3, 0x01},
	} {
		testConnReceiveFails(t, data, ErrInvalidMessageOpcode)
	}

	// Write an invalid message ID to the client pipe, and make sure that
	// receiving fails.
	for _, data := range [][]byte{
		// String.
		{0x92, 0x00, 0xa1, 0x00},

		// Signed 32 bit integer < 0.
		{0x92, 0x00, 0xd2, 0xff, 0xff, 0xff, 0xff},

		// 64 bit unsigned integer.
		{0x92, 0x00, 0xcf, 0x55, 0x12, 0xc5, 0x16, 0x55, 0x12, 0xc5, 0x16},
	} {
		testConnReceiveFails(t, data, ErrInvalidMessageId)
	}

	// Write a valid opcode and message ID with no more elements, and make sure
	// that receiving fails.
	testConnReceiveFails(t, []byte{0x92, 0x00, 0x00}, ErrBadMessage)
}

func testConnSendRequestReceive(t *testing.T, method string, arguments []interface{}, trace bool) {
	clientConn, serverConn := newTestingConnPipe()
	defer clientConn.Close()
	defer serverConn.Close()

	sentMessageId, err := clientConn.SendRequest(method, arguments, trace)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
		return
	}

	msg, err := serverConn.Receive()
	if err != nil {
		t.Errorf("Error receiving message: %v", err)
		return
	}

	if msg.MessageId() != sentMessageId {
		t.Errorf("Sent and received message IDs do not match: %d sent, %d received", sentMessageId, msg.MessageId())
		return
	}

	req, ok := msg.(*RequestMessage)
	if !ok {
		t.Errorf("Received message is not a request")
		return
	}

	if req.Method != method {
		t.Errorf("Expected request for method '%s', but request is for '%s'", method, req.Method)
	}

	if req.Arguments == nil {
		t.Errorf("Received request has nil arguments")
		return
	} else if len(req.Arguments) != len(arguments) {
		t.Errorf("Unexpected number of request arguments, expected %d, got %d", len(arguments), len(req.Arguments))
	} else {
		for i, expected := range arguments {
			actual := req.Arguments[i]

			if actual != expected {
				t.Errorf("Expected argument %d to be %v, but it is %v", i, expected, actual)
			}
		}
	}

	if req.Trace != trace {
		t.Errorf("Expected trace to be %v but it is %v", trace, req.Trace)
	}
}

func testConnSendRequestCompressedReceive(t *testing.T, method string, arguments []interface{}, trace bool) {
	clientConn, serverConn := newTestingConnPipe()
	defer clientConn.Close()
	defer serverConn.Close()

	sentMessageId := clientConn.nextMessageId()
	sentMsg := &RequestMessage{
		messageId: sentMessageId,
		Method:    method,
		Arguments: arguments,
		Trace:     trace,
	}

	err := clientConn.sendCompressed(sentMsg, nil, SnappyCompression)
	if err != nil {
		t.Errorf("Error sending compressed message: %v", err)
		return
	}

	msg, err := serverConn.Receive()
	if err != nil {
		t.Errorf("Error receiving message: %v", err)
		return
	}

	if msg.MessageId() != sentMessageId {
		t.Errorf("Sent and received message IDs do not match: %d sent, %d received", sentMessageId, msg.MessageId())
		return
	}

	req, ok := msg.(*RequestMessage)
	if !ok {
		t.Errorf("Received message is not a request")
		return
	}

	if req.Method != method {
		t.Errorf("Expected request for method '%s', but request is for '%s'", method, req.Method)
	}

	if req.Arguments == nil {
		t.Errorf("Received request has nil arguments")
		return
	} else if len(req.Arguments) != len(arguments) {
		t.Errorf("Unexpected number of request arguments, expected %d, got %d", len(arguments), len(req.Arguments))
	} else {
		for i, expected := range arguments {
			actual := req.Arguments[i]

			if actual != expected {
				t.Errorf("Expected argument %d to be %v, but it is %v", i, expected, actual)
			}
		}
	}

	if req.Trace != trace {
		t.Errorf("Expected trace to be %v but it is %v", trace, req.Trace)
	}
}

func testConnSendNotificationReceive(t *testing.T, method string, arguments []interface{}) {
	clientConn, serverConn := newTestingConnPipe()
	defer clientConn.Close()
	defer serverConn.Close()

	sentMessageId, err := clientConn.SendNotification(method, arguments)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
		return
	}

	msg, err := serverConn.Receive()
	if err != nil {
		t.Errorf("Error receiving message: %v", err)
		return
	}

	if msg.MessageId() != sentMessageId {
		t.Errorf("Sent and received message IDs do not match: %d sent, %d received", sentMessageId, msg.MessageId())
		return
	}

	req, ok := msg.(*NotificationMessage)
	if !ok {
		t.Errorf("Received message is not a notification")
		return
	}

	if req.Method != method {
		t.Errorf("Expected request for method '%s', but request is for '%s'", method, req.Method)
	}

	if req.Arguments == nil {
		t.Errorf("Received request has nil arguments")
		return
	} else if len(req.Arguments) != len(arguments) {
		t.Errorf("Unexpected number of request arguments, expected %d, got %d", len(arguments), len(req.Arguments))
	} else {
		for i, expected := range arguments {
			actual := req.Arguments[i]

			if actual != expected {
				t.Errorf("Expected argument %d to be %v, but it is %v", i, expected, actual)
			}
		}
	}
}

func testConnRespondExceptionReceive(t *testing.T, exception error, trace interface{}) {
	clientConn, serverConn := newTestingConnPipe()
	defer clientConn.Close()
	defer serverConn.Close()

	sentMessageId := MessageId(123)

	err := serverConn.RespondException(exception, &RequestMessage{
		messageId: sentMessageId,
	}, trace)
	if err != nil {
		t.Errorf("Error sending exception response: %v", err)
		return
	}

	msg, err := clientConn.Receive()
	if err != nil {
		t.Errorf("Error receiving message: %v", err)
		return
	}

	if msg.MessageId() != sentMessageId {
		t.Errorf("Sent and received message IDs do not match: %d sent, %d received", sentMessageId, msg.MessageId())
		return
	}

	exc, ok := msg.(*ExceptionMessage)
	if !ok {
		t.Errorf("Received message is not an exception")
		return
	}

	var expectedNamespace, expectedName, expectedDescription string

	if eErr, ok := exception.(Error); ok {
		expectedNamespace = eErr.Namespace()
		expectedName = eErr.Name()
		expectedDescription = eErr.Error()
	} else {
		expectedNamespace = "entangle"
		expectedName = "InternalServerError"
		expectedDescription = exception.Error()
	}

	if exc.Namespace != expectedNamespace {
		t.Errorf("Expected exception namespace to be '%s' but it is '%s'", expectedNamespace, exc.Namespace)
	}

	if exc.Name != expectedName {
		t.Errorf("Expected exception name to be '%s' but it is '%s'", expectedName, exc.Name)
	}

	if exc.Description != expectedDescription {
		t.Errorf("Expected exception description to be '%s' but it is '%s'", expectedDescription, exc.Description)
	}
}

func testConnRespondResponseReceive(t *testing.T, result interface{}, trace interface{}) {
	clientConn, serverConn := newTestingConnPipe()
	defer clientConn.Close()
	defer serverConn.Close()

	sentMessageId := MessageId(123)

	err := serverConn.RespondResponse(result, &RequestMessage{
		messageId: sentMessageId,
	}, trace)
	if err != nil {
		t.Errorf("Error sending exception response: %v", err)
		return
	}

	msg, err := clientConn.Receive()
	if err != nil {
		t.Errorf("Error receiving message: %v", err)
		return
	}

	if msg.MessageId() != sentMessageId {
		t.Errorf("Sent and received message IDs do not match: %d sent, %d received", sentMessageId, msg.MessageId())
		return
	}

	resp, ok := msg.(*ResponseMessage)
	if !ok {
		t.Errorf("Received message is not a response")
		return
	}

	switch resp.Result.(type) {
	case []interface{}:
		actualArr := resp.Result.([]interface{})
		expectedArr := result.([]interface{})

		if len(actualArr) != len(expectedArr) {
			t.Errorf("Expected response result to have an array length of %d, but is has an array length of %d", len(expectedArr), len(actualArr))
		}

		for i, expected := range expectedArr {
			actual := actualArr[i]

			if actual != expected {
				t.Errorf("Expected response result array item %d to be %v but it is %v", i, expected, actual)
			}
		}

	default:
		if resp.Result != result {
			t.Errorf("Expected response result to be %v, but it is %v", result, resp.Result)
		}
	}
}

// Test SendRequest and subsequent Receive.
func TestConnSendRequestReceive(t *testing.T) {
	testConnSendRequestReceive(t, "MethodName", []interface{}{}, false)
	testConnSendRequestReceive(t, "MethodName", []interface{}{}, true)
	testConnSendRequestReceive(t, "MethodName", []interface{}{
		"Foo",
		int64(123),
	}, false)
	testConnSendRequestReceive(t, "MethodName", []interface{}{
		"Foo",
		int64(123),
	}, true)

	testConnSendRequestCompressedReceive(t, "MethodName", []interface{}{}, false)
	testConnSendRequestCompressedReceive(t, "MethodName", []interface{}{}, true)
	testConnSendRequestCompressedReceive(t, "MethodName", []interface{}{
		"Foo",
		int64(123),
	}, false)
	testConnSendRequestCompressedReceive(t, "MethodName", []interface{}{
		"Foo",
		int64(123),
	}, true)
}

// Test SendNotification and subsequent Receive.
func TestConnSendNotificationReceive(t *testing.T) {
	testConnSendNotificationReceive(t, "MethodName", []interface{}{})
	testConnSendNotificationReceive(t, "MethodName", []interface{}{
		"Foo",
		int64(123),
	})
}

// Test RespondException and subsequent Receive.
func TestConnRespondExceptionReceive(t *testing.T) {
	// Test with non-Entangle error.
	testConnRespondExceptionReceive(t, errors.New("non-entangle error"), nil)

	// Test with Entangle error.
	def := NewErrorDefinition("testing", "TestError")

	testConnRespondExceptionReceive(t, def.New("Something went awry"), nil)
}

// Test RespondResponse and subsequent Receive.
func TestConnRespondResponseReceive(t *testing.T) {
	testConnRespondResponseReceive(t, nil, nil)
	testConnRespondResponseReceive(t, []interface{}{}, nil)
	testConnRespondResponseReceive(t, "Test", nil)
	testConnRespondResponseReceive(t, uint64(12346), nil)
	testConnRespondResponseReceive(t, []interface{}{"Hello", int64(123)}, nil)
}
