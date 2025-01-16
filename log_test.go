package zigbee2mqtt_test

import (
	"testing"

	"github.com/max-weis/go-zigbee2mqtt"
)

func TestLogHandlerInitialization(t *testing.T) {
	handler := zigbee2mqtt.NewLogHandler()

	if handler == nil {
		t.Fatal("LogHandler should not be nil after initialization")
	}

	if handler.GetTopic() != zigbee2mqtt.TOPIC_LOG {
		t.Fatalf("Expected topic %s, got %s", zigbee2mqtt.TOPIC_LOG, handler.GetTopic())
	}

	if len(handler.GetLogs()) != 0 {
		t.Fatal("LogHandler should start with an empty log channel")
	}
}

func TestLogHandler_HandleMessage(t *testing.T) {
	handler := zigbee2mqtt.NewLogHandler()
	message := zigbee2mqtt.NewMockMQTTMessage("Test log message")

	handler.HandleMessage(nil, message)
	logs := handler.GetLogs()

	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0] != "Test log message" {
		t.Fatalf("Expected log message 'Test log message', got '%s'", logs[0])
	}
}

func TestLogHandler_GetLogs(t *testing.T) {
	handler := zigbee2mqtt.NewLogHandler()
	messages := []string{"Log 1", "Log 2", "Log 3"}

	// Send messages
	for _, msg := range messages {
		handler.HandleMessage(nil, zigbee2mqtt.NewMockMQTTMessage(msg))
	}

	// Retrieve logs
	logs := handler.GetLogs()

	if len(logs) != len(messages) {
		t.Fatalf("Expected %d logs, got %d", len(messages), len(logs))
	}

	// Verify messages
	for i, log := range logs {
		if log != messages[i] {
			t.Fatalf("Expected log '%s', got '%s'", messages[i], log)
		}
	}

	// Verify that channel is now empty
	if len(handler.GetLogs()) != 0 {
		t.Fatal("LogHandler.GetLogs() should return an empty slice after reading all logs")
	}
}
