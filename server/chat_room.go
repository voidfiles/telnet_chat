package server

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

//Message is a message bound for chatRoom
type Message struct {
	ClientID int
	Sent     time.Time
	Text     string
}

//ChatRoom manages a persistant chat
type ChatRoom struct {
	logger           *zerolog.Logger
	history          []Message
	inboundMessages  chan Message
	outboundMessages chan string
}

//NewChatRoom creates a ChatRoom
func NewChatRoom(logger *zerolog.Logger, inboundMessages chan Message, outboundMessages chan string) *ChatRoom {
	return &ChatRoom{
		logger:           logger,
		inboundMessages:  inboundMessages,
		outboundMessages: outboundMessages,
		history:          make([]Message, 0),
	}
}

func (cr *ChatRoom) formatMessage(msg Message) string {
	return fmt.Sprintf("%s - Client %d > %s", msg.Sent.Format(time.RFC3339), msg.ClientID, msg.Text)
}

//AddMessage appends a message to the chat room and then broacasts it
func (cr *ChatRoom) AddMessage(msg Message) {
	cr.history = append(cr.history, msg)
	go func() {
		cr.outboundMessages <- cr.formatMessage(msg)
	}()
}

func (cr *ChatRoom) Run() {
	for {
		select {
		case msg := <-cr.inboundMessages:
			cr.AddMessage(msg)
		}
	}
}
