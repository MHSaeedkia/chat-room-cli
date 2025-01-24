package nats

import (
	"time"

	natsPkg "github.com/nats-io/nats.go"
)

type NatsInterface interface {
	Recive(subj string) (chan []byte, error)
	Send(subj string, message []byte) error
	Request(subj string, message []byte, delay time.Duration) error
	Response(subj string, message []byte) error
}

type NatsConnection *natsPkg.Conn

type Nats struct {
	NatsConn *natsPkg.Conn
}

func NewNats(natsURL string) (NatsInterface, error) {
	nc, err := natsPkg.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	return &Nats{
		NatsConn: nc,
	}, err
}
