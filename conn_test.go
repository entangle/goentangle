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

func testConnRaiseExceptionReceive(t *testing.T, exception error, trace Trace) {
	clientConn, serverConn := newTestingConnPipe()
	defer clientConn.Close()
	defer serverConn.Close()

	sentMessageId := MessageId(123)

	err := serverConn.RaiseException(exception, &RequestMessage{
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

	var expectedDefinition, expectedName, expectedDescription string

	if eErr, ok := exception.(Exception); ok {
		expectedDefinition = eErr.Definition()
		expectedName = eErr.Name()
		expectedDescription = eErr.Error()
	} else {
		expectedDefinition = "entangle"
		expectedName = "InternalServerError"
		expectedDescription = exception.Error()
	}

	if exc.Definition != expectedDefinition {
		t.Errorf("Expected exception definition to be '%s' but it is '%s'", expectedDefinition, exc.Definition)
	}

	if exc.Name != expectedName {
		t.Errorf("Expected exception name to be '%s' but it is '%s'", expectedName, exc.Name)
	}

	if exc.Description != expectedDescription {
		t.Errorf("Expected exception description to be '%s' but it is '%s'", expectedDescription, exc.Description)
	}

	if trace != nil && exc.Trace == nil {
		t.Errorf("Expected exception to contain trace")
	}
}

func testConnRespondReceive(t *testing.T, result interface{}, trace Trace) {
	clientConn, serverConn := newTestingConnPipe()
	defer clientConn.Close()
	defer serverConn.Close()

	sentMessageId := MessageId(123)

	err := serverConn.Respond(result, &RequestMessage{
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

	if trace != nil && resp.Trace == nil {
		t.Errorf("Expected response to contain trace")
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

// Test RaiseException and subsequent Receive.
func TestConnRaiseExceptionReceive(t *testing.T) {
	for _, trace := range []Trace{
		nil,
		NewTrace("Test"),
	} {
		// Test with non-Entangle error.
		testConnRaiseExceptionReceive(t, errors.New("non-entangle error"), trace)

		// Test with Entangle error.
		def := NewExceptionDefinition("testing", "TestError")

		testConnRaiseExceptionReceive(t, def.New("Something went awry"), trace)
	}
}

// Test Respond and subsequent Receive.
func TestConnRespondReceive(t *testing.T) {
	for _, trace := range []Trace{
		nil,
		NewTrace("Test"),
	} {
		testConnRespondReceive(t, nil, trace)
		testConnRespondReceive(t, []interface{}{}, trace)
		testConnRespondReceive(t, "Test", trace)
		testConnRespondReceive(t, uint64(12346), trace)
		testConnRespondReceive(t, []interface{}{"Hello", int64(123)}, trace)
	}
}

// Test AcknowledgeNotification and subsequent Receive.
func TestConnAcknowledgeNotificationReceive(t *testing.T) {
	clientConn, serverConn := newTestingConnPipe()
	defer clientConn.Close()
	defer serverConn.Close()

	sentMessageId := MessageId(123)

	err := serverConn.AcknowledgeNotification(&RequestMessage{
		messageId: sentMessageId,
	})
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

	if _, ok := msg.(*NotificationAcknowledgementMessage); !ok {
		t.Errorf("Received message is not a notification acknowledgement")
		return
	}
}
