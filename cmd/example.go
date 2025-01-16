package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/max-weis/go-zigbee2mqtt"
)

func main() {
	// Define MQTT client options
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://raspberrypi:1883"). 
		SetClientID("zigbee2mqtt-client").
		SetCleanSession(true)

	// Create Paho MQTT client
	rawClient := mqtt.NewClient(opts)

	// Connect to the broker
	if token := rawClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}

	// Wrap the Paho client in RealMQTTClient
	client := zigbee2mqtt.NewPahoClient(rawClient)

	// Create Zigbee2MQTT instance with a 5-second timeout
	z2m, err := zigbee2mqtt.NewZigbee2MQTT(client, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to initialize Zigbee2MQTT: %v", err)
	}

	// Periodically check the state
	for {
		state, err := z2m.State()
		if err != nil {
			if err == zigbee2mqtt.ErrStateTimeout {
				fmt.Println("Timeout occurred. Using cached state:", state)
			} else {
				log.Fatalf("Error retrieving state: %v", err)
			}
		} else {
			fmt.Println("State received:", state)
			if state {
				os.Exit(0)
			}
		}

		time.Sleep(1 * time.Second) // Poll every second
	}
}
