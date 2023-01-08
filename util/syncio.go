package util

import (
	"io"
	"sync"
)

type syncConn struct {
	io.ReadWriteCloser
	lock *sync.Mutex
}

func (s *syncConn) Read(p []byte) (n int, err error) {
	return s.ReadWriteCloser.Read(p)
}

func (s *syncConn) Write(p []byte) (n int, err error) {
	s.lock.Lock()
	n, err = s.ReadWriteCloser.Write(p)
	s.lock.Unlock()
	return
}

func (s *syncConn) Close() error {
	return s.Close()
}

func WrapConn(conn io.ReadWriteCloser) io.ReadWriteCloser {
	return &syncConn{conn, &sync.Mutex{}}
}
