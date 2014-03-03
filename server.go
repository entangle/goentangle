package goentangle

import (
	"net"
)

// Server.
type Server interface {
	// Accept.
	Accept(l net.Listener) error

	// Serve connections until the listener is closed.
	Serve(l net.Listener) error

	// Serve connection.
	ServeConn(conn *Conn)

	// Wait for all connections to finish.
	//
	// Only valid when using Serve.
	Wait()
}
