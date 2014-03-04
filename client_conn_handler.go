package goentangle

import (
	"errors"
	"io"
	"sync"
)

// Client connection is shut down.
var ErrShutdown = errors.New("connection is shut down")

// Client connection handler.
//
// Convenience type for handling requests and responses for clients.
type ClientConnHandler struct {
	// Underlying connection.
	conn *Conn

	// Pending messages.
	pending map[MessageId]chan Message

	// Lock for pending table.
	pendingLock sync.Mutex

	// Closing.
	closing  bool

	// Shut down.
    shutdown bool

	// Lock for state.
	stateLock sync.Mutex
}

// Receive messages and dispatch.
func (h *ClientConnHandler) receive() {
	// Receive messages for as long as possible.
	var err error

	for err == nil {
		var msg Message
		if msg, err = h.conn.Receive(); err != nil {
			break
		}

		// Ignore anything that is not a response or exception.
		ignore := true

		switch msg.(type) {
		case *ResponseMessage, *ExceptionMessage:
			ignore = false
		}

		if ignore {
			continue
		}

		// Determine the recipient of the message.
		var done chan Message

		h.pendingLock.Lock()
		done, _ = h.pending[msg.MessageId()]
		delete(h.pending, msg.MessageId())
		h.pendingLock.Unlock()

		// Send the message to the recipient if possible.
		if done != nil {
			done <- msg
		}
	}

	// Close all pending message channels.
	h.pendingLock.Lock()
	h.stateLock.Lock()

	h.shutdown = true

	for _, done := range h.pending {
		done <- nil
	}

	h.stateLock.Unlock()
	h.pendingLock.Unlock()
}

// Call a remote function.
func (h *ClientConnHandler) Call(method string, args []interface{}, notify bool, trace bool) (resp Message, err error) {
	// Make sure we're in normal operating state.
	h.stateLock.Lock()
	if h.shutdown || h.closing {
		err = ErrShutdown
		h.stateLock.Unlock()
		return
	}
	h.stateLock.Unlock()

	// Acquire the lock for the pending table.
	h.pendingLock.Lock()

	// Send the message.
	var msgId MessageId

	if notify {
		msgId, err = h.conn.SendNotification(method, args)
	} else {
		msgId, err = h.conn.SendRequest(method, args, trace)
	}

	if err != nil {
		if err == io.EOF {
			err = ErrShutdown
		}

		h.pendingLock.Unlock()
		return
	}

	if notify {
		h.pendingLock.Unlock()
		return
	}

	// Allocate a channel for awaiting the response, update the pending table
	// and unlock.
	done := make(chan Message, 1)
	h.pending[msgId] = done
	h.pendingLock.Unlock()

	// Wait for the response.
	defer close(done)
	if resp = <- done; resp == nil {
		err = ErrShutdown
		return
	}

	return
}

// Close the connection.
func (h *ClientConnHandler) Close() error {
	h.stateLock.Lock()
	defer h.stateLock.Unlock()
	if h.closing || h.shutdown {
		return ErrShutdown
	}
	h.closing = true
	h.conn.Close()
	return nil
}

// New client connection handler.
func NewClientConnHandler(conn *Conn) (h *ClientConnHandler) {
	h = &ClientConnHandler{
		conn:    conn,
		pending: make(map[MessageId]chan Message),
	}

	go h.receive()

	return
}
