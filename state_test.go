package zigbee2mqtt_test

import (
	"errors"
	"testing"
	"time"

	"github.com/max-weis/go-zigbee2mqtt"
)

func TestStateHandlerInitialization(t *testing.T) {
	timeout := 100 * time.Millisecond
	handler := zigbee2mqtt.NewStateHandler(timeout)

	if handler == nil {
		t.Fatal("StateHandler should not be nil after initialization")
	}

	if handler.GetTopic() != zigbee2mqtt.TOPIC_STATE {
		t.Fatalf("Expected topic %s, got %s", zigbee2mqtt.TOPIC_STATE, handler.GetTopic())
	}

	state, err := handler.GetState()
	if state != false || err != zigbee2mqtt.ErrStateTimeout {
		t.Fatalf("Expected initial state to be false with ErrStateTimeout, got state: %v, error: %v", state, err)
	}
}

func TestStateHandler_HandleMessage(t *testing.T) {
	timeout := 100 * time.Millisecond
	handler := zigbee2mqtt.NewStateHandler(timeout)

	// Simulate a message indicating "online"
	message := zigbee2mqtt.NewMockMQTTMessage("online")
	handler.HandleMessage(nil, message)

	// Check the state
	state, err := handler.GetState()
	if state != true || err != nil {
		t.Fatalf("Expected state to be true with no error, got state: %v, error: %v", state, err)
	}

	// Simulate a message indicating "offline"
	message = zigbee2mqtt.NewMockMQTTMessage("offline")
	handler.HandleMessage(nil, message)

	// Check the state
	state, err = handler.GetState()
	if state != false || err != nil {
		t.Fatalf("Expected state to be false with no error, got state: %v, error: %v", state, err)
	}
}

func TestStateHandler_Timeout(t *testing.T) {
	timeout := 100 * time.Millisecond
	handler := zigbee2mqtt.NewStateHandler(timeout)

	// Simulate no incoming message and wait for timeout
	state, err := handler.GetState()
	if state != false || !errors.Is(err, zigbee2mqtt.ErrStateTimeout) {
		t.Fatalf("Expected state to be false with ErrStateTimeout, got state: %v, error: %v", state, err)
	}
}
