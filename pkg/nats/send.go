package nats

func (nats *Nats) Send(subj string, message []byte) error {
	return nats.NatsConn.Publish(subj, message)
}
