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

type Chat struct {
	db             *inMemoryDb
	natsConnection *nats.Conn
}

type inMemoryDb struct {
	storage map[uuId]Client
	mtx     sync.RWMutex
}

type Client struct {
	UserName string `json:"userName"`
	UserId   uuId   `json:"userId"`
	Message  string `json:"Message"`
	Online   bool   `json:"online"`
}

// NewChat creates and returns a new instance of Chat
func NewChat(natsURL string) (*Chat, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	return &Chat{
		db: &inMemoryDb{
			storage: make(map[uuId]Client),
			mtx:     sync.RWMutex{},
		},
		natsConnection: nc,
	}, nil
}

func (chat *Chat) RegisterClient() {
	chat.natsConnection.Subscribe("Client.Register", func(msg *nats.Msg) {
		client := Client{}
		err := json.Unmarshal(msg.Data, &client)
		if err != nil {
			log.Fatal(err)
		}

		// add client to db .
		chat.db.mtx.Lock()
		chat.db.storage[client.UserId] = client
		chat.db.mtx.Unlock()

		fmt.Printf("Welcome dear %s to cli chat room !\n", client.UserName)
		go chat.NewClient(&client)
	})
}

func (chat *Chat) NewClient(client *Client) {
	// send join message to all client exept registerd client .
	go chat.SendMessage([]byte("joind to chat"), client)

	// listen for incomming message from registerd clinet .
	go chat.ReciveMessage(fmt.Sprintf("%s", client.UserId), client)
}

func (chat *Chat) SendMessage(msg []byte, client *Client) {
	chat.db.mtx.RLock()
	defer chat.db.mtx.RUnlock()
	var payLoad string
	switch string(msg) {
	case "#users":
		payLoad += fmt.Sprintf("%s-\n**%-12s**\n", client.UserName, "Online users")
		for _, clnt := range chat.db.storage {
			if clnt.Online {
				payLoad += fmt.Sprintf("**%-12s**\n", clnt.UserName)
			}
		}
		chat.natsConnection.Publish(fmt.Sprintf("%s.%s", "Server", client.UserId), []byte(payLoad))
	default:
		for uuid := range chat.db.storage {
			payLoad := fmt.Sprintf("%s-", client.UserName) + string(msg)
			if uuid != client.UserId {
				chat.natsConnection.Publish(fmt.Sprintf("%s.%s", "Server", uuid), []byte(payLoad))
			}
		}

	}
}

func (chat *Chat) ReciveMessage(topic string, client *Client) {
	chat.natsConnection.Subscribe(topic, func(msg *nats.Msg) {
		fmt.Printf("Message from %s : %s\n", client.UserName, string(msg.Data))

		// send incomming message from registerd clinet to all others client .
		go chat.SendMessage(msg.Data, client)
	})
}

func main() {
	chat, err := NewChat(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to " + nats.DefaultURL)

	// register new client
	go chat.RegisterClient()

	// Keep the connection alive
	runtime.Goexit()
}
