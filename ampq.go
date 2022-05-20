package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

var ampqStopChannel chan struct{}
var ampqTag = "[AMPQ]"

func StartAmpq(connectionString string, exchange string, queue string) {
	ampqStopChannel = make(chan struct{})
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, ampqTag, "Unable to connect to Endpoint", err)
		return
	}
	Log(LOG_INFO, ampqTag, "Connected TO Endpoint")
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		Log(LOG_ERROR, ampqTag, "Unable to create AMPQ channel", err)
		return
	}
	Log(LOG_INFO, ampqTag, "Channel Created")
	defer ch.Close()
	_, err = ch.QueueDeclare(queue, false, true, false, true, nil)
	if err != nil {
		Log(LOG_ERROR, ampqTag, "Unable to create queue", err)
		return
	}
	Log(LOG_INFO, ampqTag, "Creating Queue", queue)

	err = ch.QueueBind(queue, "", exchange, false, nil)
	if err != nil {
		Log(LOG_ERROR, ampqTag, "Unable to bind to Exchange", err)
		return
	}
	Log(LOG_INFO, ampqTag, "Queue Binding Done")

	go startConsumer(ch, queue)
	<-ampqStopChannel
	_, err = ch.QueueDelete(queue, false, false, true)
	if err != nil {
		Log(LOG_WARNING, ampqTag, "Unable to delete Queue")
	} else {
		Log(LOG_INFO, ampqTag, "Queue Deleted")
	}

}
func startConsumer(ch *amqp.Channel, queue string) {
	msgs, err := ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		Log(LOG_ERROR, ampqTag, "Unable to read messages from queue", err)
		ampqStopChannel <- struct{}{}
		return
	}
	for message := range msgs {
		go processMessage(message.Body)
	}
}
func processMessage(message []byte) {
	var cmd AMPQ_Message
	err := json.Unmarshal(message, &cmd)
	if err != nil {
		Log(LOG_ERROR, ampqTag, "Invalid Message Queue", err)
		return
	}
	Log(LOG_DEBUG, ampqTag, cmd)
	if cmd.Action == "delete" {
		userDB.Delete(cmd.Address)
	} else if cmd.Action == "add" {
		userDB.Add(cmd.Address, true, -1)
	}
}
func StopAmpq() {
	ampqStopChannel <- struct{}{}
}
