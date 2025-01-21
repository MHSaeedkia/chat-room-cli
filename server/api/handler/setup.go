package handler

import (
	"chat-room-cli/server/internal/chat"
	"fmt"
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
)

func Run() error {
	chat, err := chat.NewChat(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to " + nats.DefaultURL)

	// register new Client
	go chat.RegisterClient()

	// Keep the connection alive
	runtime.Goexit()
	return nil
}
