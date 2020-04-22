package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connect(clientID string, uri *url.URL) (mqtt.Client, error) {
	var opts = mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientID)

	var client = mqtt.NewClient(opts)
	var token = client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	return client, token.Error()
}

func main() {
	var brokerURI = flag.String("brokerURI", "mqtt://127.0.0.1:1883", "URI of the MQTT broker")
	var clientID = flag.String("clientID", "mqtt-publish", "client ID for MQTT")
	var topic = flag.String("topic", "some-topic", "MQTT topic to subscribe to")
	var qos = flag.Int("qos", 0, "QOS value of the message")
	var payload = flag.String("payload", time.Now().Format(time.RFC1123Z), "payload to publish. defaults to current time")
	var retained = flag.Bool("retained", false, "whether to publish a retained value")

	flag.Parse()

	mqttURI, err := url.Parse(*brokerURI)
	if err != nil {
		log.Fatal(err)
	}

	client, err := connect(*clientID, mqttURI)
	if err != nil {
		log.Fatal(err)
	}

	token := client.Publish(*topic, byte(*qos), *retained, *payload)
	token.WaitTimeout(time.Millisecond * 2000)
	if token.Error() != nil {
		log.Fatal(token.Error())
	}
}
