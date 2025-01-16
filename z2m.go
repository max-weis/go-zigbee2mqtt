package zigbee2mqtt

import (
	"errors"
	"sync"
	"time"
)

const (
	TOPIC_STATE = "zigbee2mqtt/bridge/state"
	TOPIC_LOG   = "zigbee2mqtt/bridge/logging"
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

// Zigbee2MQTT handles MQTT state communication.
type Zigbee2MQTT struct {
	client  MQTTClient
	timeout time.Duration

	stateChannel chan bool
	stateMutex   sync.RWMutex
	currentState bool
}

// NewZigbee2MQTT creates a new Zigbee2MQTT instance with a specified timeout.
func NewZigbee2MQTT(client MQTTClient, timeout time.Duration) (*Zigbee2MQTT, error) {
	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	z2m := &Zigbee2MQTT{
		client:       client,
		stateChannel: make(chan bool, 1),
		timeout:      timeout,
	}

	if err := client.Subscribe(TOPIC_STATE, 0, z2m.stateCallback); err != nil {
		return nil, err
	}

	if err := client.Subscribe(TOPIC_LOG, 0, z2m.logCallback); err != nil {
		return nil, err
	}

	return z2m, nil
}

// State blocks until a new state is sent via the channel or a timeout occurs.
// If no new state is received, it returns the cached state and a sentinel error for timeout.
func (z *Zigbee2MQTT) State() (bool, error) {
	select {
	case state := <-z.stateChannel:
		z.updateCachedState(state)

		return state, nil
	case <-time.After(z.timeout):
		z.stateMutex.RLock()
		defer z.stateMutex.RUnlock()

		return z.currentState, ErrStateTimeout
	}
}

// stateCallback processes incoming messages and updates the state.
func (z *Zigbee2MQTT) stateCallback(client MQTTClient, msg MQTTMessage) {
	state := string(msg.Payload()) == "online"
	z.updateCachedState(state)

	select {
	case z.stateChannel <- state:
	default:
	}
}

// updateCachedState updates the cached state safely.
func (z *Zigbee2MQTT) updateCachedState(state bool) {
	z.stateMutex.Lock()
	defer z.stateMutex.Unlock()

	z.currentState = state
}
