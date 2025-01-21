package main

import (
	"chat-room-cli/client/internal/chat"
	"fmt"
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
)

func main() {
	chat, err := chat.NewChat(nats.DefaultURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connected to " + nats.DefaultURL)

	err = chat.NewCLient()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("> ")

	// send messages to server
	go chat.SendMessage(fmt.Sprintf("%s", chat.GetUserId()))

	// get another users message
	go chat.ReciveMessage(fmt.Sprintf("%s.%s", "Server", chat.GetUserId()))

	// Keep the connection alive
	runtime.Goexit()
}
