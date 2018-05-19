package server

import (
	"io"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type MyReadWriter struct {
	reader io.Reader
	writer io.Writer
}

func NewMyReadWriter() *MyReadWriter {
	r, w := io.Pipe()
	return &MyReadWriter{
		reader: r,
		writer: w,
	}
}

func (rw *MyReadWriter) Read(p []byte) (n int, err error) {
	return rw.reader.Read(p)
}

func (rw *MyReadWriter) Write(p []byte) (n int, err error) {
	return rw.writer.Write(p)
}

func TestNewTelnetChatConnector(t *testing.T) {
	logger := zerolog.New(os.Stdout)
	clientPool := NewClientPool(
		&logger,
		make(chan Message),
		make(chan string),
	)

	assert.IsType(t, &ClientPool{}, clientPool)
}

func TestTelnetChatConnector(t *testing.T) {
	logger := zerolog.New(os.Stdout)

	clientPool := NewClientPool(
		&logger,
		make(chan Message),
		make(chan string),
	)

	conn1 := NewMyReadWriter()
	clientPool.AddConnection(conn1)
	assert.Equal(t, 1, clientPool.ClientCount)
	_, found1 := clientPool.allClients[conn1]
	assert.True(t, found1)
	clientPool.RemoveConnection(conn1)
	assert.Equal(t, 1, clientPool.ClientCount)
	_, found2 := clientPool.allClients[conn1]
	assert.False(t, found2)
}
