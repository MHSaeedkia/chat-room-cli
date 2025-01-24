package chat

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func (chat *Chat) NewClient(client *Client) {
	// if client equal to nil , it mean , we are in server mode
	if client != nil {
		// send join message to all Client exept registerd Client .
		go chat.SendMessage("", []byte("joind to chat"), client)

		// listen for incomming message from registerd clinet .
		subj := fmt.Sprintf("%s", client.UserId)
		go chat.ReciveMessage(subj, client)

		// check about online or not .
		go chat.CheckOnline(client)

	} else {
		fmt.Printf("Welcome to chat room ..\nWhat is your name : ")
		chat.Scanner.Scan()
		client := Client{
			UserName: chat.Scanner.Text(),
			UserId:   uuId(uuid.New().String()),
		}
		data, _ := json.Marshal(client)

		chat.Client = &client

		chat.Nats.Send(CLIENT_REGISTER, data)
	}

}
