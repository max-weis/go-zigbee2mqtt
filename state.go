package zigbee2mqtt

import (
	"sync"
	"time"
)

const TOPIC_STATE = "zigbee2mqtt/bridge/state"

// StateHandler handles the state topic.
type StateHandler struct {
	stateChannel chan bool
	stateMutex   sync.RWMutex
	currentState bool
	timeout      time.Duration
}

// NewStateHandler creates a new StateHandler.
func NewStateHandler(timeout time.Duration) *StateHandler {
	return &StateHandler{
		stateChannel: make(chan bool, 1),
		timeout:      timeout,
	}
}

// HandleMessage processes incoming state messages.
func (h *StateHandler) HandleMessage(_ MQTTClient, msg MQTTMessage) {
	state := string(msg.Payload()) == "online"

	h.stateMutex.Lock()
	h.currentState = state
	h.stateMutex.Unlock()

	// Send state to the channel (non-blocking)
	select {
	case h.stateChannel <- state:
	default:
	}
}

// GetTopic returns the topic for the handler.
func (h *StateHandler) GetTopic() string {
	return TOPIC_STATE
}

// GetState retrieves the current state, blocking until a new state is received or timeout occurs.
func (h *StateHandler) GetState() (bool, error) {
	select {
	case state := <-h.stateChannel:
		return state, nil
	case <-time.After(h.timeout):
		h.stateMutex.RLock()
		defer h.stateMutex.RUnlock()
		return h.currentState, ErrStateTimeout
	}
}
