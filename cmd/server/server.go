package server

import (
	"chat-room-cli/internal/chat"
	"chat-room-cli/internal/http"
	"log"
	"runtime"
)

const (
	DefaultURL = "nats://127.0.0.1:4222"
	HttpPort   = "12345"
)

func Run() error {
	chat, err := chat.NewChat(DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	// register new Client in Server
	go chat.RegisterClient()

	// run http server
	log.Fatal(http.Run(HttpPort))

	// Keep the connection alive
	runtime.Goexit()
	return nil
}
