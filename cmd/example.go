package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/max-weis/go-zigbee2mqtt"
)

func main() {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://raspberrypi:1883").
		SetClientID("zigbee2mqtt-client").
		SetCleanSession(true)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect: %v", token.Error())
	}

	z2m, err := zigbee2mqtt.NewZigbee2MQTT(zigbee2mqtt.NewPahoClient(client), 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to initialize Zigbee2MQTT: %v", err)
	}

	// Register StateHandler
	stateHandler := zigbee2mqtt.NewStateHandler(5 * time.Second)
	if err := z2m.RegisterHandler(stateHandler); err != nil {
		log.Fatalf("Failed to register state handler: %v", err)
	}

	// Register LogHandler
	logHandler := zigbee2mqtt.NewLogHandler()
	if err := z2m.RegisterHandler(logHandler); err != nil {
		log.Fatalf("Failed to register log handler: %v", err)
	}

	// Poll for state updates
	go func() {
		for {
			state, err := stateHandler.GetState()
			if err != nil {
				fmt.Println("State timeout. Cached state:", state)
			} else {
				fmt.Println("State update:", state)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	// Poll for logs
	go func() {
		for {
			logs := logHandler.GetLogs()
			for _, logMsg := range logs {
				fmt.Println("Log message:", logMsg)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	select {} // Keep the program running
}
