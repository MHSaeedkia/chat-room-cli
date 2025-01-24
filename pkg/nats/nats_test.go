package nats

import (
	"testing"
	"time"

	natsPkg "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestNewNats(t *testing.T) {
	natsURL := "nats://localhost:4222"
	natsInterface, err := NewNats(natsURL)
	assert.NoError(t, err, "Error while creating a new NATS connection")
	assert.NotNil(t, natsInterface, "NATS interface should not be nil")
}

func TestRecive(t *testing.T) {
	natsURL := "nats://localhost:4222"
	natsInterface, err := NewNats(natsURL)
	assert.NoError(t, err)
	assert.NotNil(t, natsInterface)

	nats := natsInterface.(*Nats)

	subject := "test.receive"
	messageChan, err := nats.Recive(subject)
	assert.NoError(t, err)

	// Send a test message
	go func() {
		time.Sleep(100 * time.Millisecond)
		err := nats.Send(subject, []byte("test message"))
		assert.NoError(t, err)
	}()

	// Receive the message
	select {
	case msg := <-messageChan:
		assert.Equal(t, "test message", string(msg))
	case <-time.After(1 * time.Second):
		t.Error("Did not receive the message in time")
	}
}

func TestSend(t *testing.T) {
	natsURL := "nats://localhost:4222"
	natsInterface, err := NewNats(natsURL)
	assert.NoError(t, err)
	assert.NotNil(t, natsInterface)

	nats := natsInterface.(*Nats)
	subject := "test.send"

	// Subscribe to the subject
	msgChan := make(chan *natsPkg.Msg, 1)
	_, err = nats.NatsConn.Subscribe(subject, func(msg *natsPkg.Msg) {
		msgChan <- msg
	})
	assert.NoError(t, err)

	// Send a message
	err = nats.Send(subject, []byte("test message"))
	assert.NoError(t, err)

	// Verify the received message
	select {
	case msg := <-msgChan:
		assert.Equal(t, "test message", string(msg.Data))
	case <-time.After(1 * time.Second):
		t.Error("Did not receive the message in time")
	}
}

func TestRequest(t *testing.T) {
	natsURL := "nats://localhost:4222"
	natsInterface, err := NewNats(natsURL)
	assert.NoError(t, err)
	assert.NotNil(t, natsInterface)

	nats := natsInterface.(*Nats)
	subject := "test.request"

	// Subscribe and respond to requests
	_, err = nats.NatsConn.Subscribe(subject, func(msg *natsPkg.Msg) {
		err := msg.Respond([]byte("response message"))
		assert.NoError(t, err)
	})
	assert.NoError(t, err)

	// Make a request
	msg, err := nats.NatsConn.Request(subject, []byte("request message"), 1*time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, "response message", string(msg.Data))
}

func TestResponse(t *testing.T) {
	natsURL := "nats://localhost:4222"
	natsInterface, err := NewNats(natsURL)
	assert.NoError(t, err)
	assert.NotNil(t, natsInterface)

	nats := natsInterface.(*Nats)
	subject := "test.response"

	// Send a request and listen for a response
	go func() {
		_, err := nats.NatsConn.Request(subject, []byte{}, 1*time.Second)
		assert.NoError(t, err)
	}()

	// Handle the response
	err = nats.Response(subject, []byte("response message"))
	assert.NoError(t, err)
}
