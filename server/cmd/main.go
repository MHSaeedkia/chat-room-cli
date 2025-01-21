package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/nats-io/nats.go"
)

type uuId string

type inRamMemory struct {
	storage map[uuId]Client
	mtx     sync.RWMutex
}

type Client struct {
	UserName string `json:"userName"`
	UserId   uuId   `json:"userId"`
	Online   bool   `json:"online"`
}

func listener(client Client, natsConnection *nats.Conn) {
	natsConnection.Subscribe(fmt.Sprintf("%s", client.UserId), func(msg *nats.Msg) {
		fmt.Printf("Message from %s : %s\n", client.UserName, string(msg.Data))
	})
}

func main() {
	// Create server connection
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	fmt.Println("Connected to " + nats.DefaultURL)

	Storage := inRamMemory{
		storage: map[uuId]Client{},
	}

	go func() {
		natsConnection.Subscribe("Client.Register", func(msg *nats.Msg) {
			clnt := Client{}
			err := json.Unmarshal(msg.Data, &clnt)
			if err != nil {
				log.Fatal(err)
			}
			Storage.mtx.Lock()
			Storage.storage[clnt.UserId] = clnt
			Storage.mtx.Unlock()

			fmt.Printf("Welcome dear %s to cli chat room !\n", clnt.UserName)
			go func(client Client, natsConnection *nats.Conn) {
				welcome := "joind to chat"
				go func(msg []byte, natsConnection *nats.Conn) {
					for uuid, _ := range Storage.storage {
						pay_load := fmt.Sprintf("%s-", client.UserName) + string(msg)
						if uuid != client.UserId {
							natsConnection.Publish(fmt.Sprintf("%s.%s", "Server", uuid), []byte(pay_load))
						}
					}
				}([]byte(welcome), natsConnection)

				natsConnection.Subscribe(fmt.Sprintf("%s", client.UserId), func(msg *nats.Msg) {
					fmt.Printf("Message from %s : %s\n", client.UserName, string(msg.Data))
					go func(msg []byte, natsConnection *nats.Conn) {
						for uuid, _ := range Storage.storage {
							pay_load := fmt.Sprintf("%s-", client.UserName) + string(msg)
							if uuid != client.UserId {
								natsConnection.Publish(fmt.Sprintf("%s.%s", "Server", uuid), []byte(pay_load))
							}
						}
					}(msg.Data, natsConnection)
				})
			}(clnt, natsConnection)
		})

	}()

	// Keep the connection alive
	runtime.Goexit()
}
