package goentangle

import (
	"io"
	"sync"
)

type pipe struct {
	channel      chan []byte
	lock         sync.Mutex
	readerClosed bool
	writerClosed bool
}

type pipeSender struct {
	p *pipe
}

func (s *pipeSender) Write(p []byte) (n int, err error) {
	s.p.lock.Lock()
	closed := s.p.writerClosed || s.p.readerClosed
	s.p.lock.Unlock()

	if closed {
		return 0, io.EOF
	}

	data := make([]byte, len(p))
	copy(data, p)
	s.p.channel <- data
	return len(p), nil
}

func (s *pipeSender) Close() error {
	s.p.lock.Lock()
	defer s.p.lock.Unlock()

	s.p.writerClosed = true
	close(s.p.channel)

	return nil
}

type pipeReceiver struct {
	p *pipe
}

func (r *pipeReceiver) Read(p []byte) (n int, err error) {
	r.p.lock.Lock()
	closed := r.p.readerClosed
	r.p.lock.Unlock()

	if closed {
		return 0, io.EOF
	}

	data := <-r.p.channel
	if data == nil {
		return 0, io.EOF
	}

	copy(p, data)
	return len(data), nil
}

func (s *pipeReceiver) Close() error {
	s.p.lock.Lock()
	defer s.p.lock.Unlock()

	s.p.readerClosed = true

	return nil
}

func newPipe() (sender *pipeSender, receiver *pipeReceiver) {
	p := &pipe{
		channel: make(chan []byte, 16),
	}

	return &pipeSender{
			p,
		}, &pipeReceiver{
			p,
		}
}

type testingConn struct {
	reader io.ReadCloser
	writer io.WriteCloser
}

func (c *testingConn) Read(data []byte) (n int, err error) {
	return c.reader.Read(data)
}

func (c *testingConn) Write(p []byte) (n int, err error) {
	return c.writer.Write(p)
}

func (c *testingConn) Close() error {
	c.reader.Close()
	c.writer.Close()
	return nil
}

func newTestingPipe() (clientConn *testingConn, serverConn *testingConn) {
	clientToServerWriter, clientToServerReader := newPipe()
	serverToClientWriter, serverToClientReader := newPipe()

	clientConn = &testingConn{
		reader: serverToClientReader,
		writer: clientToServerWriter,
	}

	serverConn = &testingConn{
		reader: clientToServerReader,
		writer: serverToClientWriter,
	}

	return
}

func newTestingConnPipe() (clientConn *Conn, serverConn *Conn) {
	pipeClientConn, pipeServerConn := newTestingPipe()
	return NewConn(pipeClientConn, "test server"), NewConn(pipeServerConn, "test client")
}
