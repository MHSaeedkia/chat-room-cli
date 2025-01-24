package server

import (
	"chat-room-cli/internal/chat"
	"log"
	"runtime"
)

const (
	DefaultURL = "nats://127.0.0.1:4222"
)

func Run() error {
	chat, err := chat.NewChat(DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	// register new Client in Server
	chat.RegisterClient()

	// Keep the connection alive
	runtime.Goexit()
	return nil
}
