package zigbee2mqtt

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotConnected = errors.New("mqtt client not connected")
	ErrStateTimeout = errors.New("timeout waiting for state update")
)

// MQTTClient defines an interface for the MQTT client.
type MQTTClient interface {
	IsConnected() bool
	Subscribe(topic string, qos byte, callback func(client MQTTClient, msg MQTTMessage)) error
}

// MQTTMessage defines an interface for an MQTT message.
type MQTTMessage interface {
	Payload() []byte
}

// TopicHandler defines an interface for handling specific MQTT topics.
type TopicHandler interface {
	HandleMessage(client MQTTClient, msg MQTTMessage)
	GetTopic() string
}

// Zigbee2MQTT manages MQTT communication and topic handlers.
type Zigbee2MQTT struct {
	client   MQTTClient
	timeout  time.Duration
	handlers map[string]TopicHandler // Registered handlers
	mutex    sync.RWMutex            // Protects the handlers map
}

// NewZigbee2MQTT creates a new Zigbee2MQTT instance.
func NewZigbee2MQTT(client MQTTClient, timeout time.Duration) (*Zigbee2MQTT, error) {
	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	z2m := &Zigbee2MQTT{
		client:   client,
		timeout:  timeout,
		handlers: make(map[string]TopicHandler),
	}

	return z2m, nil
}

// RegisterHandler registers a handler for a specific topic.
func (z *Zigbee2MQTT) RegisterHandler(handler TopicHandler) error {
	topic := handler.GetTopic()

	// Subscribe to the topic
	if err := z.client.Subscribe(topic, 0, func(client MQTTClient, msg MQTTMessage) {
		handler.HandleMessage(client, msg)
	}); err != nil {
		return err
	}

	// Add the handler to the map
	z.mutex.Lock()
	defer z.mutex.Unlock()

	z.handlers[topic] = handler

	return nil
}
