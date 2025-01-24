package nats

import (
	natsPkg "github.com/nats-io/nats.go"
)

func (nats *Nats) Response(subj string, message []byte) error {
	var err error
	nats.NatsConn.Subscribe(subj, func(msg *natsPkg.Msg) {
		err = msg.Respond([]byte{})
	})
	return err
}
