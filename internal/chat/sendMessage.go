package chat

import (
	"fmt"
	"log"
	"strings"
)

func (chat *Chat) SendMessage(subj string, message []byte, Client *Client) {
	// if client equal to nil , it mean , we are in server mode
	if Client != nil {
		chat.Db.mtx.RLock()
		var payLoad string
		switch string(message) {
		case "#users":
			payLoad += fmt.Sprintf("%s-\n**%-12s**\n", Client.UserName, "Online users")
			for _, clnt := range chat.Db.storage {
				payLoad += fmt.Sprintf("**%-12s**\n", clnt.UserName)
			}
			subj := fmt.Sprintf("%s.%s", SERVER, Client.UserId)
			chat.Nats.Send(subj, []byte(payLoad))
		default:
			if strings.Contains(string(message), "#") {
				break
			}
			for uuid := range chat.Db.storage {
				payLoad := fmt.Sprintf("%s-", Client.UserName) + string(message)
				if uuid != Client.UserId {
					subj := fmt.Sprintf("%s.%s", SERVER, uuid)
					chat.Nats.Send(subj, []byte(payLoad))
				}
			}
		}
		chat.Db.mtx.RUnlock()
	} else {
		for chat.Scanner.Scan() {
			fmt.Printf("> ")
			text := chat.Scanner.Text()
			chat.Nats.Send(subj, []byte(text))
		}
		if err := chat.Scanner.Err(); err != nil {
			log.Fatal(err.Error())
		}
	}
}
