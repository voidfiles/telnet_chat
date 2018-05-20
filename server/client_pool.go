package server

import (
	"bufio"
	"io"
	"time"

	"github.com/rs/zerolog"
)

//ClientPool manages all the clients and sends messages to a chat room
type ClientPool struct {
	logger *zerolog.Logger
	// ClientCount is the total number of people who have ever connected.
	ClientCount int

	// allClients keeps track of current connections
	allClients map[io.ReadWriter]int

	// InboundMessage is a message from a client
	InboundMessage chan Message
	// BroadcastMessage is a channel you can send messages to all clients on
	BroadcastMessage chan string
}

//NewClientPool creates and returns a TelnetChatConnector
func NewClientPool(logger *zerolog.Logger, inbound chan Message, broadcast chan string) *ClientPool {
	return &ClientPool{
		logger:           logger,
		ClientCount:      0,
		allClients:       make(map[io.ReadWriter]int),
		InboundMessage:   inbound,
		BroadcastMessage: broadcast,
	}
}

func (cp *ClientPool) listenToConnection(conn io.ReadWriter, clientID int) {
	reader := bufio.NewReader(conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg := Message{
			ClientID: clientID,
			Sent:     time.Now(),
			Text:     incoming,
		}
		cp.logger.Info().
			Int("client_id", clientID).
			Str("text", incoming).
			Msg("incomming message from client")
		cp.InboundMessage <- msg
	}

	cp.RemoveConnection(conn)
}

//AddConnection adds a connected client to the pool
func (cp *ClientPool) AddConnection(conn io.ReadWriter) {
	cp.logger.Info().Msgf("Accepted new client, #%d", cp.ClientCount)

	cp.allClients[conn] = cp.ClientCount
	go cp.listenToConnection(conn, cp.ClientCount)
	cp.ClientCount++

}

//RemoveConnection takes a connected client out of the pool
func (cp *ClientPool) RemoveConnection(conn io.ReadWriter) {
	cp.logger.Info().
		Int("client_id", cp.allClients[conn]).
		Msg("Client Disconnected")
	delete(cp.allClients, conn)
}

//Run starts the client pool on a loop waiting to distribute messages to clients
func (cp *ClientPool) Run() {
	for msg := range cp.BroadcastMessage {
		for conn := range cp.allClients {

			go func(conn io.ReadWriter, msg string) {
				_, err := conn.Write([]byte(msg))
				if err != nil {
					cp.RemoveConnection(conn)
				}
			}(conn, msg)
		}
		cp.logger.Info().
			Int("num_clients", len(cp.allClients)).
			Str("text", msg).
			Msg("message broadcast")
	}

}
