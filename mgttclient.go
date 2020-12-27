package main

import (
	"fmt"
	//import the Paho Go MQTT library
	"log"
	"net"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//watch here: http://www.hivemq.com/demos/websocket-client/
//documentation here: https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("mqtt://localhost:1883")
	opts.SetClientID("gotest")
	opts.SetUsername("testclient")
	opts.SetPassword("1qaz2wsx")
	opts.SetDefaultPublishHandler(f)
	opts.SetWill("/Device/james/lwt", "HCF - Elvis has left the building", 2, false)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe("/device/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	//dt := time.Now()
	//curr := fmt.Sprintf("Current date and time is: $d", dt.Format("01-02-2006 15:04:05"))
	//token2 := c.Publish("/testtopic/james/curtime", 1, false, curr)
	//token2.Wait()

	ipaddr := GetOutboundIP()
	//ipreport := fmt.Sprintf(ip)
	token := c.Publish("device/imac", 0, false, ipaddr)
	token.Wait()

	//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message

	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("device/james/3", 0, false, text)
		token.Wait()
	}

	time.Sleep(2500 * time.Second)

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("device/james/"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}
