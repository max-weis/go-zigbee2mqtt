package zigbee2mqtt

import (
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// RealMQTTClient is a production implementation of the MQTTClient interface.
type RealMQTTClient struct {
	client mqtt.Client
	mu     sync.Mutex // Protects callback map
	// Map to store topic-specific callbacks
	callbacks map[string]func(client MQTTClient, msg MQTTMessage)
}

// NewRealMQTTClient wraps a paho.mqtt.golang.Client and creates a RealMQTTClient.
func NewPahoClient(client mqtt.Client) *RealMQTTClient {
	return &RealMQTTClient{
		client:    client,
		callbacks: make(map[string]func(client MQTTClient, msg MQTTMessage)),
	}
}

// IsConnected checks if the MQTT client is connected.
func (r *RealMQTTClient) IsConnected() bool {
	return r.client.IsConnected()
}

// Subscribe subscribes to a topic and registers a callback for incoming messages.
func (r *RealMQTTClient) Subscribe(topic string, qos byte, callback func(client MQTTClient, msg MQTTMessage)) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Register the callback
	r.callbacks[topic] = callback

	// Use the internal paho client to subscribe
	token := r.client.Subscribe(topic, qos, func(c mqtt.Client, m mqtt.Message) {
		// Call the registered callback
		r.mu.Lock()
		cb := r.callbacks[topic]
		r.mu.Unlock()

		if cb != nil {
			cb(r, &PahoMQTTMessage{message: m})
		}
	})
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// PahoMQTTMessage is a production implementation of the MQTTMessage interface.
type PahoMQTTMessage struct {
	message mqtt.Message
}

// Payload returns the message payload.
func (m *PahoMQTTMessage) Payload() []byte {
	return m.message.Payload()
}
