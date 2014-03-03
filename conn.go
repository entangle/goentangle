package goentangle

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/vmihailenco/msgpack"
	"io"
	"sync"
	"sync/atomic"
)

var (
	ErrInvalidMessageData   = errors.New("invalid message data received")
	ErrInvalidMessageOpcode = errors.New("invalid message opcode received")
	ErrInvalidMessageId     = errors.New("invalid message ID received")
	ErrBadMessage           = errors.New("bad message received")

	compressionThreshold = 1460 * 5
)

func encodeSlice(slice []interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := msgpack.NewEncoder(buffer)
	err := encoder.Encode(slice)
	return buffer.Bytes(), err
}

// Connection.
//
// Reading is only safe from one goroutine while writing is safe in a blocking
// manner form any number of goroutines.
type Conn struct {
	// Message ID counter.
	messageIdCounter uint32

	// Description.
	description string

	// Closer.
	closer io.Closer

	// Writer.
	writer *bufio.Writer

	// Reader.
	reader *bufio.Reader

	// Decoder.
	decoder *msgpack.Decoder

	// Write lock.
	writeLock sync.Mutex
}

// New connection.
func NewConn(conn io.ReadWriteCloser, description string) *Conn {
	reader := bufio.NewReader(conn)

	return &Conn{
		description: description,
		closer:      conn,
		writer:      bufio.NewWriter(conn),
		reader:      reader,
		decoder:     msgpack.NewDecoder(reader),
	}
}

// Lock for writing and get writer.
//
// Make absolutely sure to call UnlockWriter after writing is done, preferably
// using a defer statement immediately after retrieving the writer.
func (c *Conn) lockAndWriter() *bufio.Writer {
	c.writeLock.Lock()
	return c.writer
}

// Unlock writer.
func (c *Conn) unlockWriter() {
	c.writeLock.Unlock()
}

// Close connection.
func (c *Conn) Close() {
	c.closer.Close()
}

// Description.
func (c *Conn) Description() string {
	return c.description
}

// Deserialize a message.
func (c *Conn) deserializeMessage(opcode Opcode, messageId MessageId, messageData []interface{}) (msg Message, err error) {
	// Parse the incoming message based on its opcode.
	switch opcode {
	case RequestOpcode:
		if len(messageData) != 3 {
			err = ErrBadMessage
			return
		}

		method, methodOk := messageData[0].(string)
		arguments, argumentsOk := messageData[1].([]interface{})
		trace, traceOk := messageData[2].(bool)

		if !methodOk || method == "" || !argumentsOk || arguments == nil || !traceOk {
			err = ErrBadMessage
			return
		}

		msg = &RequestMessage{
			messageId: messageId,
			Method:    method,
			Arguments: arguments,
			Trace:     trace,
		}

	case NotificationOpcode:
		if len(messageData) != 2 {
			err = ErrBadMessage
			return
		}

		method, methodOk := messageData[0].(string)
		arguments, argumentsOk := messageData[1].([]interface{})

		if !methodOk || method == "" || !argumentsOk || arguments == nil {
			err = ErrBadMessage
			return
		}

		msg = &NotificationMessage{
			messageId: messageId,
			Method:    method,
			Arguments: arguments,
		}

	case ResponseOpcode:
		if len(messageData) != 2 {
			err = ErrBadMessage
			return
		}

		var trace Trace
		var traceErr error
		result, resultOk := messageData[0].(interface{})
		rawTrace, rawTraceOk := messageData[1].(interface{})
		if rawTraceOk && rawTrace != nil {
			trace, traceErr = DeserializeTrace(rawTrace)
		}

		if (!resultOk && messageData[0] != nil) || (!rawTraceOk && messageData[1] != nil) || (traceErr != nil) {
			err = ErrBadMessage
			return
		}

		msg = &ResponseMessage{
			messageId: messageId,
			Result:    result,
			Trace:     trace,
		}

	case ExceptionOpcode:
		if len(messageData) != 4 {
			err = ErrBadMessage
			return
		}

		var trace Trace
		var traceErr error
		definition, definitionOk := messageData[0].(string)
		name, nameOk := messageData[1].(string)
		description, descriptionOk := messageData[2].(string)
		rawTrace, rawTraceOk := messageData[3].(interface{})
		if rawTraceOk {
			trace, traceErr = DeserializeTrace(rawTrace)
		}

		if !definitionOk || !nameOk || !descriptionOk || (!rawTraceOk && messageData[3] != nil) || (traceErr != nil) {
			err = ErrBadMessage
			return
		}

		msg = &ExceptionMessage{
			messageId:   messageId,
			Definition:  definition,
			Name:        name,
			Description: description,
			Trace:       trace,
		}

	case CompressedMessageOpcode:
		if len(messageData) != 2 {
			err = ErrBadMessage
			return
		}

		method, methodOk := DeserializeCompressionMethod(messageData[0])
		compressed, compressedErr := DeserializeBinary(messageData[1])
		if !methodOk || compressedErr != nil {
			err = ErrBadMessage
			return
		}

		decompressed, decompressionErr := method.Decompress(compressed)
		if decompressionErr != nil {
			err = ErrBadMessage
			return
		}

		reader := bytes.NewReader(decompressed)
		msg, err = c.readMessage(msgpack.NewDecoder(reader))

	default:
		panic("implementation error in message parsing")
	}

	return
}

// Read a message from a decoder.
func (c *Conn) readMessage(decoder *msgpack.Decoder) (msg Message, err error) {
	// Read the message.
	var messageData []interface{}
	if messageData, err = decoder.DecodeSlice(); err != nil {
		if err != io.EOF {
			err = ErrInvalidMessageData
		}

		return
	}

	// Make sure that we can parse an opcode and message ID from the message
	// data.
	if len(messageData) < 2 {
		err = ErrInvalidMessageData
		return
	}

	opcode, ok := ParseOpcode(messageData[0])
	if !ok || !opcode.Valid() {
		err = ErrInvalidMessageOpcode
		return
	}

	messageId, ok := ParseMessageId(messageData[1])
	if !ok {
		err = ErrInvalidMessageId
		return
	}

	// Deserialize the message.
	return c.deserializeMessage(opcode, messageId, messageData[2:])
}

// Receive a message.
//
// The returned error can be either io.EOF, ErrInvalidMessageData,
// ErrInvalidMessageOpcode or ErrInvalidMessageId, all of which are
// unrecoverable, or ErrBadMessage which doesn't prohibit the connection from
// continuing.
func (c *Conn) Receive() (msg Message, err error) {
	return c.readMessage(c.decoder)
}

// Write message data to the connection.
func (c *Conn) writeMessageData(data []byte) (err error) {
	writer := c.lockAndWriter()
	defer c.unlockWriter()

	var n int

	for {
		if n, err = writer.Write(data); err != nil {
			break
		}

		if n == len(data) {
			break
		}

		data = data[n:]
	}

	if err == nil {
		writer.Flush()
	}

	return
}

// Send a message.
func (c *Conn) send(msg Message) (err error) {
	// Serialize the message.
	serialized := msg.Serialize()
	var data []byte
	if data, err = encodeSlice(serialized); err != nil {
		return
	}

	// If the data size is over the compression threshold, let's compress it.
	if len(data) >= compressionThreshold {
		return c.sendCompressed(msg, data, SnappyCompression)
	}

	// Write the message.
	return c.writeMessageData(data)
}

// Send compressed message.
//
// If the message has previously been serialized, provide the serialized data,
// otherwise supply nil.
func (c *Conn) sendCompressed(msg Message, data []byte, compressionMethod CompressionMethod) (err error) {
	// Serialize the message if it has not already been serialized.
	if data == nil {
		serialized := msg.Serialize()
		if data, err = encodeSlice(serialized); err != nil {
			return
		}
	}

	// Compress the data.
	var compressedData []byte
	if compressedData, err = compressionMethod.Compress(data); err != nil {
		return
	}

	// Serialize the compressed message.
	var msgData []byte
	if msgData, err = encodeSlice([]interface{}{
		CompressedMessageOpcode,
		msg.MessageId(),
		compressionMethod,
		compressedData,
	}); err != nil {
		return
	}

	// Write the message.
	return c.writeMessageData(msgData)
}

// Get the next message ID.
func (c *Conn) nextMessageId() MessageId {
	messageId := atomic.AddUint32(&c.messageIdCounter, 1)
	return MessageId(messageId)
}

// Send a request.
func (c *Conn) SendRequest(method string, arguments []interface{}, trace bool) (MessageId, error) {
	messageId := c.nextMessageId()
	return messageId, c.send(&RequestMessage{
		messageId: messageId,
		Method:    method,
		Arguments: arguments,
		Trace:     trace,
	})
}

// Send a notification.
func (c *Conn) SendNotification(method string, arguments []interface{}) (MessageId, error) {
	messageId := c.nextMessageId()
	return messageId, c.send(&NotificationMessage{
		messageId: messageId,
		Method:    method,
		Arguments: arguments,
	})
}

// Respond with an exception.
func (c *Conn) RespondException(exception error, responseTo Message, trace Trace) error {
	// Transform the error into an Entangle error if it is not already.
	eErr, ok := exception.(Exception)

	if !ok {
		eErr = InternalServerError.New(exception.Error())
	}

	// Create and send the response.
	return c.send(&ExceptionMessage{
		messageId:   responseTo.MessageId(),
		Definition:  eErr.Definition(),
		Name:        eErr.Name(),
		Description: eErr.Error(),
		Trace:       trace,
	})
}

// Respond with a response.
func (c *Conn) RespondResponse(result interface{}, responseTo Message, trace Trace) error {
	// Create and send the response.
	return c.send(&ResponseMessage{
		messageId: responseTo.MessageId(),
		Result:    result,
		Trace:     trace,
	})
}
