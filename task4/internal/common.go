package internal

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
	QueueMes = "AMQP_QUEUE_MESSAGES"
	QueueNot = "AMQP_QUEUE_NOTIFICATIONS"
)

func SetupRMQ() (*amqp.Channel, amqp.Queue, amqp.Queue) {
	host, user, pass :=
		os.Getenv(Host),
		os.Getenv(Username),
		os.Getenv(Password)
	queueMes := os.Getenv(QueueMes)
	queueNot := os.Getenv(QueueNot)
	desc := fmt.Sprintf("amqp://%s:%s@%s/", user, pass, host)
	log.Printf("connecting to RabbitMQ at %s", desc)
	conn, err := amqp.Dial(desc)
	FailOnError(err, "failed to connect to RabbitMQ")
	//defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "failed to open a channel")
	//defer ch.Close()

	messages, err := ch.QueueDeclare(
		queueMes, // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	FailOnError(err, "failed to declare a queue")
	notifications, err := ch.QueueDeclare(
		queueNot, // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	FailOnError(err, "failed to declare a queue")

	return ch, messages, notifications
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
