package handler

import (
	"chat-room-cli/client/internal/chat"
	"fmt"
	"runtime"

	"github.com/nats-io/nats.go"
)

func Run() error {
	chat, err := chat.NewChat(nats.DefaultURL)
	if err != nil {
		return err
	}
	fmt.Println("Connected to " + nats.DefaultURL)

	err = chat.NewCLient()
	if err != nil {
		return err
	}

	fmt.Printf("> ")

	// send messages to server
	go chat.SendMessage(fmt.Sprintf("%s", chat.GetUserId()))

	// get another users message
	go chat.ReciveMessage(fmt.Sprintf("%s.%s", "Server", chat.GetUserId()))

	// check of being online
	go chat.CheckOnline(fmt.Sprintf("%s.%s.%s", "Server", "Online", chat.GetUserId()))

	// Keep the connection alive
	runtime.Goexit()

	return nil
}
