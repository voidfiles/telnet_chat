package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
	server.NewHTTPChatServer(*logger, chatRoom)

	go chatRoom.Run()
	go clientPool.Run()
	go func() {
		http.ListenAndServe(fmt.Sprintf("%s:%s", config.HTTPIP, config.HTTPPort), nil)
	}()
	fmt.Printf("Chat Server is running\n")
	fmt.Printf("Telnet running at %s:%s\n", config.TelnetIP, config.TelnetPort)
	fmt.Printf("HTTP Server running at %s:%s\n", config.HTTPIP, config.HTTPPort)
	chatServer.Listen()
}
