package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

//HttpChatServer exposes some parts of the chat via http
type HttpChatServer struct {
	logger   zerolog.Logger
	chatRoom *ChatRoom
}

func NewHttpChatServer(logger zerolog.Logger, chatRoom *ChatRoom) *HttpChatServer {
	hcs := &HttpChatServer{
		logger:   logger,
		chatRoom: chatRoom,
	}

	addMessage := hcs.HandlerAddMessage()
	listMessages := hcs.ListMessages()
	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			listMessages(w, r)
		} else {
			addMessage(w, r)
		}
	})

	return hcs
}

// HandlerAddMessage creates a function that will convert a POST into Message
func (hcs *HttpChatServer) HandlerAddMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			fmt.Printf("error %s", err)
			hcs.logger.Error().AnErr("error", err)
			http.Error(w, `{"error": "invalid"}`, 400)
			return
		}

		var msg Message
		err = json.Unmarshal(b, &msg)
		if err != nil {
			fmt.Printf("error %s", err)
			hcs.logger.Error().AnErr("error", err)
			http.Error(w, `{"error": "invalid"}`, 400)
			return
		}
		defer r.Body.Close()
		msg.Sent = time.Now()
		hcs.chatRoom.AddMessage(msg)
		w.Header().Add("Content-Type", "application/json")
		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("error %s", err)
			hcs.logger.Error().AnErr("error", err)
			http.Error(w, `{"error": "invalid"}`, 400)
			return
		}
		w.Write(data)
	}
}

// HandlerAddMessage creates a function that will convert a POST into Message
func (hcs *HttpChatServer) ListMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messages := hcs.chatRoom.ListMessages()
		w.Header().Add("Content-Type", "application/json")
		data, err := json.Marshal(messages)
		if err != nil {
			http.Error(w, `{"error": "invalid"}`, 400)
			return
		}
		w.Write(data)
	}
}
