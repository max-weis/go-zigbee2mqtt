package zigbee2mqtt_test

import (
	"errors"
	"testing"
	"time"

	"github.com/max-weis/go-zigbee2mqtt"
)

// MockTopicHandler is a mock implementation of the TopicHandler interface.
type MockTopicHandler struct {
	topic    string
	messages []string
}

// NewMockTopicHandler creates a new MockTopicHandler.
func NewMockTopicHandler(topic string) *MockTopicHandler {
	return &MockTopicHandler{topic: topic, messages: []string{}}
}

// HandleMessage appends the received message payload to the messages slice.
func (h *MockTopicHandler) HandleMessage(_ zigbee2mqtt.MQTTClient, msg zigbee2mqtt.MQTTMessage) {
	h.messages = append(h.messages, string(msg.Payload()))
}

// GetTopic returns the topic for this handler.
func (h *MockTopicHandler) GetTopic() string {
	return h.topic
}

// GetMessages retrieves all messages received by the handler.
func (h *MockTopicHandler) GetMessages() []string {
	return h.messages
}

func TestZigbee2MQTTInitialization(t *testing.T) {
	mockClient := zigbee2mqtt.NewMockMQTTClient()
	mockClient.Connected = true

	z2m, err := zigbee2mqtt.NewZigbee2MQTT(mockClient, time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize Zigbee2MQTT: %v", err)
	}

	if z2m == nil {
		t.Fatal("Expected Zigbee2MQTT instance, got nil")
	}
}

func TestZigbee2MQTTInitializationFailsWhenNotConnected(t *testing.T) {
	mockClient := zigbee2mqtt.NewMockMQTTClient()
	mockClient.Connected = false

	_, err := zigbee2mqtt.NewZigbee2MQTT(mockClient, time.Second)
	if !errors.Is(err, zigbee2mqtt.ErrNotConnected) {
		t.Fatalf("Expected ErrNotConnected, got: %v", err)
	}
}

func TestRegisterHandler(t *testing.T) {
	mockClient := zigbee2mqtt.NewMockMQTTClient()
	mockClient.Connected = true

	z2m, err := zigbee2mqtt.NewZigbee2MQTT(mockClient, time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize Zigbee2MQTT: %v", err)
	}

	handler := NewMockTopicHandler("test/topic")

	err = z2m.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Failed to register handler: %v", err)
	}

	if mockClient.Callback == nil {
		t.Fatal("Expected callback to be set in the mock client after handler registration")
	}
}

func TestHandleMessage(t *testing.T) {
	mockClient := zigbee2mqtt.NewMockMQTTClient()
	mockClient.Connected = true

	z2m, err := zigbee2mqtt.NewZigbee2MQTT(mockClient, time.Second)
	if err != nil {
		t.Fatalf("Failed to initialize Zigbee2MQTT: %v", err)
	}

	handler := NewMockTopicHandler("test/topic")
	err = z2m.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Failed to register handler: %v", err)
	}

	// Simulate a message
	mockClient.SendMessage("Hello, World!")

	if len(handler.GetMessages()) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(handler.GetMessages()))
	}

	if handler.GetMessages()[0] != "Hello, World!" {
		t.Fatalf("Expected message 'Hello, World!', got '%s'", handler.GetMessages()[0])
	}
}
