// docker run -d --hostname my-rabbit --name some-rabbit rabbitmq:3
// docker run -d -p 5672:5672 --hostname devrabbit --name devrabbit rabbitmq
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func listen(msgs <-chan amqp.Delivery) {
	log.Println("listening for new messages....")
	// can only read from. can neve rput anything in the channel
	// hence the arrow
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

	q, err := channel.QueueDeclare("messagingQ", false, false, false, false, nil)
	msgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	go listen(msgs) // if you put go in front of the func, it creates a new routine!

	neverEnd := make(chan bool) // make a new channel
	<-neverEnd                  // means i want to read a bool from this channel
	// if it cannot read anything, then that means it will stop until you put something in there
}
