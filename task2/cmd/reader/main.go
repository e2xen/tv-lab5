package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

const (
	Host     = "AMQP_HOST"
	Username = "AMQP_USER"
	Password = "AMQP_PASSWORD"
	Queue    = "AMQP_QUEUE"
)

func main() {
	host, user, pass :=
		os.Getenv(Host),
		os.Getenv(Username),
		os.Getenv(Password)
	queue := os.Getenv(Queue)
	desc := fmt.Sprintf("amqp://%s:%s@%s/", user, pass, host)
	log.Printf("connecting to RabbitMQ at %s", desc)
	conn, err := amqp.Dial(desc)
	failOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	var forever chan struct{}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
