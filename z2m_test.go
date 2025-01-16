package zigbee2mqtt_test

import (
	"testing"
	"time"

	"github.com/max-weis/go-zigbee2mqtt"
)

func TestCachedStateTimout(t *testing.T) {
	mockClient := zigbee2mqtt.NewMockMQTTClient()
	mockClient.Connected = true

	z2m, err := zigbee2mqtt.NewZigbee2MQTT(mockClient, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Error creating Zigbee2MQTT: %v", err)
	}

	state, err := z2m.State()
	if err != zigbee2mqtt.ErrStateTimeout {
		t.Fatalf("Expected ErrStateTimeout, got %v", err)
	}
	if state {
		t.Fatalf("Expected state to be false, got %v", state)
	}
}

func TestStateUpdate(t *testing.T) {
	mockClient := zigbee2mqtt.NewMockMQTTClient()
	mockClient.Connected = true

	z2m, err := zigbee2mqtt.NewZigbee2MQTT(mockClient, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Error creating Zigbee2MQTT: %v", err)
	}

	mockClient.SendMessage("online")
	
	state, err := z2m.State()
	if err != nil {
		t.Fatalf("Error retrieving state: %v", err)
	}
	if !state {
		t.Fatalf("Expected state to be false, got %v", state)
	}
}

func TestStateFallback(t *testing.T) {
	mockClient := zigbee2mqtt.NewMockMQTTClient()
	mockClient.Connected = true

	z2m, err := zigbee2mqtt.NewZigbee2MQTT(mockClient, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Error creating Zigbee2MQTT: %v", err)
	}

	mockClient.SendMessage("offline")
	
	state, err := z2m.State()
	if err != nil {
		t.Fatalf("Error retrieving state: %v", err)
	}
	if state {
		t.Fatalf("Expected state to be false, got %v", state)
	}
}
