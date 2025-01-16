package zigbee2mqtt

import (
	"errors"
)

// MockMQTTClient is a mock implementation of the MQTTClient interface.
type MockMQTTClient struct {
	Connected bool
	Callback  func(client MQTTClient, msg MQTTMessage)
}

func (m *MockMQTTClient) IsConnected() bool {
	return m.Connected
}

// NewMockMQTTClient creates a new MockMQTTClient.
func NewMockMQTTClient() *MockMQTTClient {
	return &MockMQTTClient{}
}

func (m *MockMQTTClient) Subscribe(topic string, qos byte, callback func(client MQTTClient, msg MQTTMessage)) error {
	if !m.IsConnected() {
		return errors.New("not connected")
	}
	m.Callback = callback
	return nil
}

// SendMessage simulates receiving a message on the mock client.
func (m *MockMQTTClient) SendMessage(payload string) {
	if m.Callback != nil {
		m.Callback(m, &MockMQTTMessage{payload: []byte(payload)})
	}
}

// MockMQTTMessage is a mock implementation of the MQTTMessage interface.
type MockMQTTMessage struct {
	payload []byte
}

func NewMockMQTTMessage(payload string) *MockMQTTMessage {
	return &MockMQTTMessage{payload: []byte(payload)}
}

func (m *MockMQTTMessage) Payload() []byte {
	return m.payload
}
