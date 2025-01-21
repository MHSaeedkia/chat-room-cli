package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type uuId string

type Chat struct {
	scanner        *bufio.Scanner
	natsConnection *nats.Conn
	client         *Client
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

	scnr := bufio.NewScanner(os.Stdin)

	clnt := &Client{}

	return &Chat{
		scanner:        scnr,
		natsConnection: nc,
		client:         clnt,
	}, nil
}

func (chat *Chat) NewCLient() error {
	fmt.Printf("Welcome to chat room..\nWhat is your name : ")
	chat.scanner.Scan()
	client := Client{
		UserName: chat.scanner.Text(),
		UserId:   uuId(uuid.New().String()),
		Online:   true,
	}
	data, err := json.Marshal(client)
	if err != nil {
		return err
	}

	chat.client = &client

	err = chat.natsConnection.Publish("Client.Register", data)
	return err
}

func (chat *Chat) SendMessage(topic string) {
	for chat.scanner.Scan() {
		fmt.Printf("> ")
		text := chat.scanner.Text()
		err := chat.natsConnection.Publish(topic, []byte(text))
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	if err := chat.scanner.Err(); err != nil {
		log.Fatal(err.Error())
	}
}

func (chat *Chat) ReciveMessage(topic string) {
	chat.natsConnection.Subscribe(topic, func(msg *nats.Msg) {
		data := strings.Split(string(msg.Data), "-")
		if data[0] == chat.client.UserName {
			fmt.Printf("New message : %s\n", string(data[1]))
		} else {
			fmt.Printf("New message from %s : %s\n", data[0], string(data[1]))
		}
		fmt.Printf("> ")
	})
}

func main() {
	chat, err := NewChat(nats.DefaultURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connected to " + nats.DefaultURL)

	err = chat.NewCLient()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("> ")

	// send messages to server
	go chat.SendMessage(fmt.Sprintf("%s", chat.client.UserId))

	// get another users message
	go chat.ReciveMessage(fmt.Sprintf("%s.%s", "Server", chat.client.UserId))

	// Keep the connection alive
	runtime.Goexit()
}
