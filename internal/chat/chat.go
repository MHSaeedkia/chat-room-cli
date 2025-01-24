package chat

import (
	"bufio"
	"chat-room-cli/pkg/nats"
	"os"
	"sync"
)

const (
	CLIENT_REGISTER = "Client.Register"
	SERVER          = "Server"
)

type ChatInterface interface {
	RegisterClient()
	NewClient(Client *Client)
	SendMessage(subj string, message []byte, Client *Client)
	ReciveMessage(subj string, Client *Client)
	GetUserId() uuId
	CheckOnline(Client *Client)
}

type uuId string

type Client struct {
	UserName string `json:"userName"`
	UserId   uuId   `json:"userId"`
	Message  string `json:"Message"`
	Online   bool   `json:"online"`
}

type inMemoryDb struct {
	storage map[uuId]Client
	mtx     sync.RWMutex
}

type Chat struct {
	Client  *Client
	Scanner *bufio.Scanner
	Nats    nats.NatsInterface
	Db      *inMemoryDb
}

// NewChat creates and returns a new instance of Chat
func NewChat(natsURL string) (ChatInterface, error) {
	nats, err := nats.NewNats(natsURL)
	if err != nil {
		return nil, err
	}

	// decaler new scaner
	scnr := bufio.NewScanner(os.Stdin)

	clnt := &Client{}

	return &Chat{
		Client:  clnt,
		Scanner: scnr,
		Db: &inMemoryDb{
			storage: make(map[uuId]Client),
			mtx:     sync.RWMutex{},
		},
		Nats: nats,
	}, nil
}
