package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/minio/minio-go"
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
	rabbitmqUser, rabbitmqPass, rabbitmqURI                 string
	minioAccessKey, minioSecretKey, minioURI                string
	minioUseSSL                                             bool
	nurseQueryDatasetURI, nurseQueryHealth, nurseFetchQueue string

	bucketName = "response-data"

	amqpConnection *amqp.Connection
	amqpChannel    *amqp.Channel
	nurseQueue     amqp.Queue

	minioClient *minio.Client

	retryDelay time.Duration = 10
)

func init() {
	rabbitmqUser = ensureEnv("RABBITMQ_DEFAULT_USER")
	rabbitmqPass = ensureEnv("RABBITMQ_DEFAULT_PASS")
	rabbitmqURI = ensureEnv("RABBITMQ_URI")

	minioAccessKey = ensureEnv("MINIO_ACCESS_KEY")
	minioSecretKey = ensureEnv("MINIO_SECRET_KEY")
	minioURI = ensureEnv("MINIO_URI")
	minioUseSSLRaw := ensureEnv("MINIO_USE_SSL")
	minioUseSSL = minioUseSSLRaw == "true"

	nurseQueryURI := ensureEnv("NURSEQUERY_URI")
	nurseFetchQueue = ensureEnv("NURSEFETCH_QUEUE")

	nurseQueryDatasetURI = fmt.Sprintf("http://%s/v1/dataset/upload/gob", nurseQueryURI)
	nurseQueryHealth = fmt.Sprintf("http://%s/_ah/health", nurseQueryURI)
}

func init() {
	var err error

	amqpConnection, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", rabbitmqUser, rabbitmqPass, rabbitmqURI))
	failOnError(err, "Failed to connect to RabbitMQ")

	amqpChannel, err = amqpConnection.Channel()
	failOnError(err, "Failed to open a channel")

	nurseQueue, err = amqpChannel.QueueDeclare(
		nurseFetchQueue, // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	minioClient, err = minio.New(minioURI, minioAccessKey, minioSecretKey, minioUseSSL)
	failOnError(err, "Failed to coonect to minio")
}

func main() {
	defer amqpConnection.Close()
	defer amqpChannel.Close()

	ensureBucket()

	stage0 := amqpConsume()

	go func() {
		for msg := range stage0 {
			log.Println(msg)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	wait := make(chan bool)
	<-wait
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ensureEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Environment variable %s missing.", key)
	}
	return value
}

func ensureBucket() {
	exists, err := minioClient.BucketExists(bucketName)
	failOnError(err, "Failed to check bucket")
	if !exists {
		err = minioClient.MakeBucket(bucketName, "hell")
		failOnError(err, "Failed to create bucket")
	}
}

func amqpConsume() <-chan NurseFetchTask {
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

	out := make(chan NurseFetchTask)

	go func() {
		for d := range msgs {
			data := NurseFetchTask{}
			json.Unmarshal(d.Body, &data)

			tagPage := NewTagPage(data.Domain, data.Tag)
			log.Printf("Received a message: %#v", tagPage)

			out <- data
		}
	}()

	return out
}

func writeToMinio(filename string, data []byte) error {
	_, err := minioClient.PutObject(
		bucketName,
		filename,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: "application/json"},
	)
	return err
}
