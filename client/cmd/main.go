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

	return &Chat{
		scanner:        scnr,
		natsConnection: nc,
	}, nil
}

func (chat *Chat) NewCLient() (Client, error) {
	fmt.Printf("Welcome to chat room..\nWhat is your name : ")
	chat.scanner.Scan()
	client := Client{
		UserName: chat.scanner.Text(),
		UserId:   uuId(uuid.New().String()),
		Online:   true,
	}
	data, err := json.Marshal(client)
	if err != nil {
		return Client{}, err
	}
	err = chat.natsConnection.Publish("Client.Register", data)
	return client, err
}

func (chat *Chat) Publisher(client *Client) {
	for chat.scanner.Scan() {
		fmt.Printf("> ")
		text := chat.scanner.Text()
		err := chat.natsConnection.Publish(fmt.Sprintf("%s", client.UserId), []byte(text))
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	if err := chat.scanner.Err(); err != nil {
		log.Fatal(err.Error())
	}
}

func (chat *Chat) Subscriber(topic string) {
	chat.natsConnection.Subscribe(topic, func(msg *nats.Msg) {
		data := strings.Split(string(msg.Data), "-")
		fmt.Printf("New message from %s : %s\n", data[0], string(data[1]))
		fmt.Printf("> ")
	})
}

func main() {
	chat, err := NewChat(nats.DefaultURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connected to " + nats.DefaultURL)

	client, err := chat.NewCLient()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("> ")

	// send messages to server
	go chat.Publisher(&client)

	// get another users message
	go chat.Subscriber(fmt.Sprintf("%s.%s", "Server", client.UserId))

	// Keep the connection alive
	runtime.Goexit()
}
