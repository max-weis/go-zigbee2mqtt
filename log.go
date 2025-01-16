package zigbee2mqtt

const TOPIC_LOG = "zigbee2mqtt/bridge/logging"

// LogHandler handles the log topic.
type LogHandler struct {
	logChannel chan string
}

// NewLogHandler creates a new LogHandler.
func NewLogHandler() *LogHandler {
	return &LogHandler{
		logChannel: make(chan string, 10),
	}
}

// HandleMessage processes incoming log messages.
func (h *LogHandler) HandleMessage(_ MQTTClient, msg MQTTMessage) {
	log := string(msg.Payload())

	// Send log to the channel (non-blocking)
	select {
	case h.logChannel <- log:
	default:
	}
}

// GetTopic returns the topic for the handler.
func (h *LogHandler) GetTopic() string {
	return TOPIC_LOG
}

// GetLogs retrieves the logs from the channel.
func (h *LogHandler) GetLogs() []string {
	logs := []string{}
	for {
		select {
		case log := <-h.logChannel:
			logs = append(logs, log)
		default:
			return logs
		}
	}
}
