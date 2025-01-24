package chat

import (
	"fmt"
	"time"
)

func (chat *Chat) CheckOnline(Client *Client) {
	//
	if Client != nil {
		// it check each five secund online users .
		for {
			time.Sleep(5 * time.Second)
			subj := fmt.Sprintf("%s.%s.%s", "Server", "Online", Client.UserId)
			err := chat.Nats.Request(subj, []byte{}, 1*time.Second)
			if err != nil {
				chat.Db.mtx.Lock()
				delete(chat.Db.storage, Client.UserId)
				chat.Db.mtx.Unlock()
			}
		}
	} else {
		subj := fmt.Sprintf("%s.%s.%s", "Server", "Online", chat.GetUserId())
		chat.Nats.Response(subj, []byte{})
	}
}
