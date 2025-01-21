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

type Client struct {
	UserName string `json:"userName"`
	UserId   uuId   `json:"userId"`
	Message  string `json:"Message"`
	Online   bool   `json:"online"`
}

func NewCLient(natsConnection *nats.Conn, scanner *bufio.Scanner) (Client, error) {
	fmt.Printf("Welcome to chat room..\nWhat is your name : ")
	scanner.Scan()
	name := scanner.Text()
	newUUID := uuid.New()
	clnt := Client{
		UserName: name,
		UserId:   uuId(newUUID.String()),
		Online:   true,
	}
	data, err := json.Marshal(clnt)
	if err != nil {
		return Client{}, err
	}
	err = natsConnection.Publish("Client.Register", data)
	return clnt, err
}

func main() {
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	fmt.Println("Connected to " + nats.DefaultURL)
	scanner := bufio.NewScanner(os.Stdin)

	client, err := NewCLient(natsConnection, scanner)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	fmt.Printf("> ")

	go func() {
		for scanner.Scan() {
			fmt.Printf("> ")
			text := scanner.Text()
			err := natsConnection.Publish(fmt.Sprintf("%s", client.UserId), []byte(text))
			if err != nil {
				log.Fatal(err.Error())
				break
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err)
		}
	}()

	go func() {
		natsConnection.Subscribe(fmt.Sprintf("%s.%s", "Server", client.UserId), func(msg *nats.Msg) {
			data := strings.Split(string(msg.Data), "-")
			fmt.Printf("New message from %s : %s\n", data[0], string(data[1]))
			fmt.Printf("> ")
		})
	}()

	// Keep the connection alive
	runtime.Goexit()
}
