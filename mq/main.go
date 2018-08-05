package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func listen(msgs <-chan amqp.Delivery) {
	log.Println("listening for new messages...")
	for msg := range msgs {
		log.Println(string(msg.Body))
	}
}

func main() {
	mqAddr := os.Getenv("MQADDR")
	if len(mqAddr) == 0 {
		mqAddr = "localhost:5672"
	}
	mqURL := fmt.Sprintf("amqp://%s", mqAddr)
	conn, err := amqp.Dial(mqURL)
	if err != nil {
		log.Fatalf("error connecting to RabbitMQ: %v", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("error creating channel: %v", err)
	}
	q, err := channel.QueueDeclare("testQ", false, false, false, false, nil)

	msgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	go listen(msgs)

	neverEnd := make(chan bool)
	<-neverEnd
}
