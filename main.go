package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/voidfiles/telnet_chat/server"
)

func main() {
	var configPath = flag.String("configpath", "config", "Path to config file")
	flag.Parse()

	config, err := server.ReadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", config)
	logfile, err := os.Create(config.LogPath)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to open (%s): %s", config.LogPath, err))
	}
	defer logfile.Close()

	// Create shared communication channels
	logger := server.NewLogger("chat-app", false, logfile)
	inboundMessages := make(chan server.Message)
	outboudMessages := make(chan string)

	clientPool := server.NewClientPool(logger, inboundMessages, outboudMessages)
	chatRoom := server.NewChatRoom(logger, inboundMessages, outboudMessages)
	chatServer := server.NewChatServer(logger, config, clientPool)

	go chatRoom.Run()
	go clientPool.Run()
	chatServer.Listen()
}
