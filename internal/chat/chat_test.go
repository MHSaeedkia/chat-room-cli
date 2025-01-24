package chat

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockNats struct {
	messages map[string][]byte
	mtx      sync.Mutex
}

func NewMockNats() *MockNats {
	return &MockNats{
		messages: make(map[string][]byte),
	}
}

func (m *MockNats) Recive(subj string) (chan []byte, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	ch := make(chan []byte, 1)
	if msg, exists := m.messages[subj]; exists {
		ch <- msg
	}
	return ch, nil
}

func (m *MockNats) Send(subj string, message []byte) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.messages[subj] = message
	return nil
}

func (m *MockNats) Request(subj string, message []byte, delay time.Duration) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, exists := m.messages[subj]; !exists {
		return fmt.Errorf("No response for subject: %s", subj)
	}
	return nil
}

func (m *MockNats) Response(subj string, message []byte) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.messages[subj] = message
	return nil
}

func TestRegisterClient(t *testing.T) {
	mockNats := NewMockNats()

	chat := &Chat{
		Nats: mockNats,
		Db: &inMemoryDb{
			storage: make(map[uuId]Client),
			mtx:     sync.RWMutex{},
		},
	}

	client := Client{
		UserName: "test-user",
		UserId:   uuId(uuid.New().String()),
	}

	data, _ := json.Marshal(client)
	err := mockNats.Send(CLIENT_REGISTER, data)
	assert.NoError(t, err, "Failed to send registration message")

	go chat.RegisterClient()

	time.Sleep(100 * time.Millisecond)

	chat.Db.mtx.RLock()
	defer chat.Db.mtx.RUnlock()

	storedClient, exists := chat.Db.storage[client.UserId]
	assert.True(t, exists, "Client should be stored in the database")
	assert.Equal(t, client.UserName, storedClient.UserName, "Client username should match")
}

func TestCheckOnline(t *testing.T) {
	mockNats := NewMockNats()

	client := &Client{
		UserName: "test-client",
		UserId:   uuId(uuid.New().String()),
	}

	chat := &Chat{
		Nats: mockNats,
		Db: &inMemoryDb{
			storage: make(map[uuId]Client),
			mtx:     sync.RWMutex{},
		},
	}

	chat.Db.mtx.Lock()
	chat.Db.storage[client.UserId] = *client
	chat.Db.mtx.Unlock()

	subj := fmt.Sprintf("%s.%s.%s", SERVER, ONLINE, client.UserId)
	mockNats.Response(subj, []byte{})

	go chat.CheckOnline(client)

	time.Sleep(100 * time.Millisecond)

	chat.Db.mtx.RLock()
	defer chat.Db.mtx.RUnlock()

	_, exists := chat.Db.storage[client.UserId]
	assert.True(t, exists, "Client should remain online")
}

func TestSendMessage(t *testing.T) {
	mockNats := NewMockNats()

	client := &Client{
		UserName: "sender",
		UserId:   uuId(uuid.New().String()),
	}

	chat := &Chat{
		Nats: mockNats,
		Db: &inMemoryDb{
			storage: map[uuId]Client{
				client.UserId: *client,                                       // Add the sender to the DB
				"receiver-id": {UserName: "receiver", UserId: "receiver-id"}, // Add a receiver
			},
			mtx: sync.RWMutex{},
		},
	}

	message := "hello world"

	// Trigger SendMessage
	go chat.SendMessage("", []byte(message), client)

	// Wait for the asynchronous operation
	time.Sleep(100 * time.Millisecond)

	// Verify if the message was sent to the expected receivers
	for _, receiver := range chat.Db.storage {
		if receiver.UserId != client.UserId {
			subj := fmt.Sprintf("%s.%s", SERVER, receiver.UserId)
			msg, exists := mockNats.messages[subj]
			assert.True(t, exists, "Message should be sent to the NATS")
			assert.Contains(t, string(msg), client.UserName, "Message should include sender's username")
			assert.Contains(t, string(msg), message, "Message should include the original message")
		}
	}
}

func TestReciveMessage(t *testing.T) {
	mockNats := NewMockNats()

	client := &Client{
		UserName: "receiver",
		UserId:   uuId(uuid.New().String()),
	}

	chat := &Chat{
		Nats: mockNats,
		Db: &inMemoryDb{
			storage: map[uuId]Client{
				client.UserId: *client,
			},
			mtx: sync.RWMutex{},
		},
	}

	subj := client.UserId
	mockNats.Send(string(subj), []byte("test message"))

	go chat.ReciveMessage(string(subj), client)

	time.Sleep(100 * time.Millisecond)

	// Since we're just verifying output, there is no return value to assert.
	// This test mainly ensures the absence of errors and that the function runs.
}

func TestGetUserId(t *testing.T) {
	client := &Client{
		UserId: uuId(uuid.New().String()),
	}

	chat := &Chat{
		Client: client,
	}

	assert.Equal(t, client.UserId, chat.GetUserId(), "GetUserId should return the client's UserId")
}
