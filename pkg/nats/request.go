package nats

import (
	"time"
)

func (nats *Nats) Request(subj string, message []byte, delay time.Duration) error {
	_, err := nats.NatsConn.Request(subj, []byte{}, 1*time.Second)
	return err
}
