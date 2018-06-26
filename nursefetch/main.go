package main

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type (
	NurseFetchTask struct {
		Domain string `json:"domain" binding:"required"`
		Tag    string `json:"tag" binding:"required"`
	}
)

var (
	amqpConnection *amqp.Connection
	amqpChannel    *amqp.Channel
	nurseQueue     amqp.Queue

	uploadURL      = "http://nursequery:8080/v1/dataset/upload/gob"
	healthCheckURL = "http://nursequery:8080/_ah/health"

	waitPeriod time.Duration = 10
)

func main() {
	var err error

	amqpConnection, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer amqpConnection.Close()

	amqpChannel, err = amqpConnection.Channel()
	failOnError(err, "Failed to open a channel")
	defer amqpChannel.Close()

	nurseQueue, err = amqpChannel.QueueDeclare(
		"nurse-fetch", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := amqpChannel.Consume(
		nurseQueue.Name, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			data := NurseFetchTask{}
			json.Unmarshal(d.Body, &data)

			tagPage := NewTagPage(data.Domain, data.Tag)
			log.Printf("Received a message: %#v", tagPage)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	forever := make(chan bool)
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
