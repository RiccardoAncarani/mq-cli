package queue

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Queue struct {
	connection *amqp.Connection
	channel    string
}

func (q *Queue) OpenQueue(connectionString string, channel string) {
	var err error
	q.channel = channel
	q.connection, err = amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")

}

func (q *Queue) PublishStringMessage(message string) {
	ch, err := q.connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	t, err := ch.QueueDeclare(
		q.channel, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		t.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	failOnError(err, "Failed to publish a message")

}

func (q *Queue) ReadStringMessage() {
	ch, err := q.connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	t, err := ch.QueueDeclare(
		q.channel, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		t.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Println(string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (q *Queue) CloseQueue() {
	q.connection.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
