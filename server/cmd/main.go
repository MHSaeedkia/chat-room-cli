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
	NatsConnection *nats.Conn
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
		NatsConnection: nc,
	}, nil
}

func (chat *Chat) RegisterClient() {
	chat.NatsConnection.Subscribe("Client.Register", func(msg *nats.Msg) {
		Client := Client{}
		err := json.Unmarshal(msg.Data, &Client)
		if err != nil {
			log.Fatal(err)
		}

		// add Client to db .
		chat.db.mtx.Lock()
		chat.db.storage[Client.UserId] = Client
		chat.db.mtx.Unlock()

		fmt.Printf("Welcome dear %s to cli chat room !\n", Client.UserName)
		go chat.NewClient(&Client)
	})
}

func (chat *Chat) NewClient(Client *Client) {
	// send join message to all Client exept registerd Client .
	go chat.SendMessage([]byte("joind to chat"), Client)

	// listen for incomming message from registerd clinet .
	go chat.ReciveMessage(fmt.Sprintf("%s", Client.UserId), Client)
}

func (chat *Chat) SendMessage(msg []byte, Client *Client) {
	chat.db.mtx.RLock()
	defer chat.db.mtx.RUnlock()
	var payLoad string
	switch string(msg) {
	case "#users":
		payLoad += fmt.Sprintf("%s-\n**%-12s**\n", Client.UserName, "Online users")
		for _, clnt := range chat.db.storage {
			if clnt.Online {
				payLoad += fmt.Sprintf("**%-12s**\n", clnt.UserName)
			}
		}
		chat.NatsConnection.Publish(fmt.Sprintf("%s.%s", "Server", Client.UserId), []byte(payLoad))
	default:
		for uuid := range chat.db.storage {
			payLoad := fmt.Sprintf("%s-", Client.UserName) + string(msg)
			if uuid != Client.UserId {
				chat.NatsConnection.Publish(fmt.Sprintf("%s.%s", "Server", uuid), []byte(payLoad))
			}
		}

	}
}

func (chat *Chat) ReciveMessage(topic string, Client *Client) {
	chat.NatsConnection.Subscribe(topic, func(msg *nats.Msg) {
		fmt.Printf("Message from %s : %s\n", Client.UserName, string(msg.Data))

		// send incomming message from registerd clinet to all others Client .
		go chat.SendMessage(msg.Data, Client)
	})
}

func main() {
	chat, err := NewChat(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to " + nats.DefaultURL)

	// register new Client
	go chat.RegisterClient()

	// Keep the connection alive
	runtime.Goexit()
}
