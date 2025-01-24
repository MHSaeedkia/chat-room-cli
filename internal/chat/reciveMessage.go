package chat

import (
	"fmt"
	"strings"
)

func (chat *Chat) ReciveMessage(subj string, Client *Client) {
	// if client equal to nil , it mean , we are in server mode
	if Client != nil {
		message, _ := chat.Nats.Recive(subj)
		defer close(message)
		for msg := range message {
			fmt.Printf("Message from %s : %s\n", Client.UserName, string(msg))
			chat.SendMessage("", msg, Client)
		}
	} else {
		message, _ := chat.Nats.Recive(subj)
		defer close(message)
		for msg := range message {
			data := strings.Split(string(msg), "-")
			if data[0] == chat.Client.UserName {
				fmt.Printf("New message : %s\n", string(data[1]))
			} else {
				fmt.Printf("New message from %s : %s\n", data[0], string(data[1]))
			}
			fmt.Printf("> ")
		}
	}
}
