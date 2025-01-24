package chat

import (
	"encoding/json"
	"fmt"
	"log"
)

func (chat *Chat) RegisterClient() {
	message, err := chat.Nats.Recive(CLIENT_REGISTER)
	if err != nil {
		log.Fatal(err)
	}
	defer close(message)

	for msg := range message {

		if msg != nil {
			Client := Client{}
			err = json.Unmarshal(msg, &Client)
			if err != nil {
				log.Fatal(err)
			}

			// add Client to db .
			chat.Db.mtx.Lock()
			chat.Db.storage[Client.UserId] = Client
			chat.Db.mtx.Unlock()

			fmt.Printf("Welcome dear %s to cli chat room !\n", Client.UserName)

			// add new client to chat room
			chat.NewClient(&Client)
		}
	}
}
