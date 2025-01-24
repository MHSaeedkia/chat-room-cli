package chat

import (
	"fmt"
	"log"
	"time"
)

const (
	CHECKTIME = 2 * time.Second
	TIMEOUT   = 3 * time.Second
)

func (chat *Chat) CheckOnline(Client *Client) {
	//
	if Client != nil {
		// it check each five secund online users .
		subj := fmt.Sprintf("%s.%s.%s.%s", SERVER, ONLINE, Client.UserId, CLIENT)

		chat.Nats.Response(subj, []byte{})
		for {
			time.Sleep(CHECKTIME)
			subj := fmt.Sprintf("%s.%s.%s", SERVER, ONLINE, Client.UserId)
			err := chat.Nats.Request(subj, []byte{}, 1*time.Second)
			if err != nil {
				chat.Db.mtx.Lock()
				delete(chat.Db.storage, Client.UserId)
				chat.Db.mtx.Unlock()
			}
		}
	} else {
		subj := fmt.Sprintf("%s.%s.%s", SERVER, ONLINE, chat.GetUserId())
		chat.Nats.Response(subj, []byte{})
		for {
			time.Sleep(TIMEOUT)
			subj := fmt.Sprintf("%s.%s.%s.%s", SERVER, ONLINE, chat.GetUserId(), CLIENT)
			err := chat.Nats.Request(subj, []byte{}, 1*time.Second)
			if err != nil {
				log.Fatal("Server is Down .")
			}
		}
	}
}
