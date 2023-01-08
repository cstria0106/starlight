package util

import (
	"io"

	"github.com/gorilla/websocket"
)

type wsConnReadWriteCloser struct {
	conn *websocket.Conn
}

func (c *wsConnReadWriteCloser) Read(p []byte) (n int, err error) {
	_, bytes, err := c.conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	return copy(p, bytes), nil
}

func (c *wsConnReadWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), c.conn.WriteMessage(websocket.BinaryMessage, p)
}

func (c *wsConnReadWriteCloser) Close() error {
	return c.conn.Close()
}

func WebsocketConnToReadWriteCloser(conn *websocket.Conn) io.ReadWriteCloser {
	return &wsConnReadWriteCloser{conn: conn}
}
