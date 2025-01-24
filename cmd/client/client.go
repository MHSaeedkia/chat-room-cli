package client

import (
	"fmt"
	"runtime"

	chatRoom "chat-room-cli/internal/chat"
)

const (
	DefaultURL = "nats://127.0.0.1:4222"
)

func Run() error {
	chat, err := chatRoom.NewChat(DefaultURL)
	if err != nil {
		return err
	}

	var client *chatRoom.Client

	chat.NewClient(client)

	fmt.Printf("> ")

	// send messages to server
	subj := fmt.Sprintf("%s", chat.GetUserId())
	go chat.SendMessage(subj, []byte(""), client)

	// get another users message
	subj = fmt.Sprintf("%s.%s", "Server", chat.GetUserId())
	go chat.ReciveMessage(subj, client)

	// check of being online
	go chat.CheckOnline(client)

	// Keep the connection alive
	runtime.Goexit()

	return nil
}
