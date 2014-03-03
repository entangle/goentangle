package goentangle

import (
	"net"
)

// Server.
type Server interface {
	// Accept accepts connections on the listener and serves requests for each
	// incoming connection.
	//
	// Accept blocks; the caller typically invokes it in a go statement.
	Accept(l net.Listener) error

	// Serve connections until the listener is closed.
	Serve(l net.Listener) error

	// ServeConn runs the server on a single connection.
	//
	// ServeConn blocks, serving the connection until the client hangs up. The
	// caller typically invokes ServeConn in a go statement.
	ServeConn(conn *Conn)

	// Wait for all connections to finish.
	//
	// Only valid when using Serve.
	Wait()
}
