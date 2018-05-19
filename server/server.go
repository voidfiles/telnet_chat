package server

import (
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog"
)

// ChatServer exposes a ChatRoom via telnet
type ChatServer struct {
	clientPool *ClientPool
	server     net.Listener
	logger     *zerolog.Logger
}

//NewChatServer constructs a ChatServer and returns it
func NewChatServer(logger *zerolog.Logger, config Config, clientPool *ClientPool) *ChatServer {
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.IP, config.Port))
	if err != nil {
		logger.Fatal().AnErr("error", err).Msg("Failed to launch server")
		os.Exit(1)
	}

	return &ChatServer{
		logger:     logger,
		server:     server,
		clientPool: clientPool,
	}
}

//Listen starts listening for inbound connections
func (cs *ChatServer) Listen() {
	for {
		conn, err := cs.server.Accept()
		if err != nil {
			cs.logger.Fatal().AnErr("error", err).Msg("Failed to launch server")
			os.Exit(1)
		}
		cs.clientPool.AddConnection(conn)
	}
}
