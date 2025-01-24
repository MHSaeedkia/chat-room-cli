package nats

import (
	natsPkg "github.com/nats-io/nats.go"
)

func (nats *Nats) Recive(subj string) (chan []byte, error) {
	message := make(chan []byte)
	_, err := nats.NatsConn.Subscribe(subj, func(msg *natsPkg.Msg) {
		message <- msg.Data
	})
	return message, err
}
