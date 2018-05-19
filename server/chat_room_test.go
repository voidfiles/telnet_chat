package server

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewChatRoom(t *testing.T) {
	logger := zerolog.New(os.Stdout)
	inboundMessages := make(chan Message)
	outboundMessages := make(chan string)
	chatRoom := NewChatRoom(&logger, inboundMessages, outboundMessages)
	assert.IsType(t, &ChatRoom{}, chatRoom)
}

func TestFormatMessage(t *testing.T) {
	msg := Message{
		Sent:     time.Date(2018, 11, 11, 11, 11, 1, 1, time.UTC),
		Text:     "Yolo",
		ClientID: 0,
	}
	logger := zerolog.New(os.Stdout)
	inboundMessages := make(chan Message)
	outboundMessages := make(chan string)
	chatRoom := NewChatRoom(&logger, inboundMessages, outboundMessages)

	text := chatRoom.formatMessage(msg)
	assert.Equal(t, "2018-11-11T11:11:01Z - Client 0 > Yolo", text)
}

func TestAddMessage(t *testing.T) {
	msg := Message{
		Sent:     time.Date(2018, 11, 11, 11, 11, 1, 1, time.UTC),
		Text:     "Yolo",
		ClientID: 0,
	}
	logger := zerolog.New(os.Stdout)
	inboundMessages := make(chan Message)
	outboundMessages := make(chan string)
	chatRoom := NewChatRoom(&logger, inboundMessages, outboundMessages)
	go func() {
		inboundMessages <- msg
	}()
	go chatRoom.Run()
	out_msg := <-outboundMessages
	assert.Equal(t, "2018-11-11T11:11:01Z - Client 0 > Yolo", out_msg)
}
