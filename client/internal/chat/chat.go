package chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type ChatInterface interface {
	NewCLient() error
	SendMessage(topic string)
	ReciveMessage(topic string)
	GetUserId() uuId
}

type Chat struct {
	Scanner        *bufio.Scanner
	NatsConnection *nats.Conn
	Client         *Client
}

// NewChat creates and returns a new instance of Chat
func NewChat(natsURL string) (ChatInterface, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	scnr := bufio.NewScanner(os.Stdin)

	clnt := &Client{}

	return &Chat{
		Scanner:        scnr,
		NatsConnection: nc,
		Client:         clnt,
	}, nil
}

// Add new client
func (chat *Chat) NewCLient() error {
	fmt.Printf("Welcome to chat room..\nWhat is your name : ")
	chat.Scanner.Scan()
	Client := Client{
		UserName: chat.Scanner.Text(),
		UserId:   uuId(uuid.New().String()),
		Online:   true,
	}
	data, err := json.Marshal(Client)
	if err != nil {
		return err
	}

	chat.Client = &Client

	err = chat.NatsConnection.Publish("Client.Register", data)
	return err
}

// Send message to server
func (chat *Chat) SendMessage(topic string) {
	for chat.Scanner.Scan() {
		fmt.Printf("> ")
		text := chat.Scanner.Text()
		err := chat.NatsConnection.Publish(topic, []byte(text))
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	if err := chat.Scanner.Err(); err != nil {
		log.Fatal(err.Error())
	}
}

// Recive message from server
func (chat *Chat) ReciveMessage(topic string) {
	chat.NatsConnection.Subscribe(topic, func(msg *nats.Msg) {
		data := strings.Split(string(msg.Data), "-")
		if data[0] == chat.Client.UserName {
			fmt.Printf("New message : %s\n", string(data[1]))
		} else {
			fmt.Printf("New message from %s : %s\n", data[0], string(data[1]))
		}
		fmt.Printf("> ")
	})
}

// Get user id
func (chat *Chat) GetUserId() uuId {
	return chat.Client.UserId
}
