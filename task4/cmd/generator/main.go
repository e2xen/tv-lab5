package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"lab5/internal"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	ch, qMes, qNot := internal.SetupRMQ()

	for {
		publish(ch, qMes, qNot)
		time.Sleep(time.Second)
	}
}

func publish(ch *amqp.Channel, qMes amqp.Queue, qNot amqp.Queue) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := int64(rand.Intn(10000))
	if r%2 == 0 {
		body := "message " + strconv.FormatInt(r, 10)
		err := ch.PublishWithContext(ctx,
			"",        // exchange
			qMes.Name, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		internal.FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent message %s\n", body)
	}

	body := "notification " + strconv.FormatInt(r, 10)
	err := ch.PublishWithContext(ctx,
		"",        // exchange
		qNot.Name, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	internal.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent notification %s\n", body)
}
