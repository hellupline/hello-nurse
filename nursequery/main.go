package main // import "github.com/hellupline/hello-nurse/nursequery"

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/streadway/amqp"

	"google.golang.org/appengine"
)

type (
	NurseFetchTask struct {
		Domain string `json:"domain" binding:"required"`
		Tag    string `json:"tag" binding:"required"`
	}
	NurseDownloadTask struct {
		Domain string `json:"domain" binding:"required"`
		URL    string `json:"url" binding:"required"`
	}
)

var (
	rabbitmqUser, rabbitmqPass, rabbitmqURI string

	nurseFetchQueueName    = "nurse-fetch"
	nurseDownloadQueueName = "nurse-download"

	amqpConnection     *amqp.Connection
	amqpChannel        *amqp.Channel
	nurseFetchQueue    amqp.Queue
	nurseDownloadQueue amqp.Queue
)

func init() {
	rabbitmqUser = ensureEnv("RABBITMQ_DEFAULT_USER")
	rabbitmqPass = ensureEnv("RABBITMQ_DEFAULT_PASS")
	rabbitmqURI = ensureEnv("RABBITMQ_URI")
}

func init() {
	var err error

	amqpConnection, err = amqp.Dial("amqp://nurse:animaniacs@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	amqpChannel, err = amqpConnection.Channel()
	failOnError(err, "Failed to open a channel")

	nurseFetchQueue, err = amqpChannel.QueueDeclare(
		nurseFetchQueueName, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare a queue")

	nurseDownloadQueue, err = amqpChannel.QueueDeclare(
		nurseDownloadQueueName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
}

func init() {
	router := gin.Default()

	http.Handle("/", router)

	router.Use(cors.Default())

	v1Group := router.Group("/v1")
	{
		favoritesGroup := v1Group.Group("/favorites")
		{
			favoritesGroup.GET("/:key", HttpHandleFavoriteRead)
			favoritesGroup.DELETE("/:key", HttpHandleFavoriteDelete)
			favoritesGroup.GET("", HttpHandleFavoriteIndex)
			favoritesGroup.POST("", HttpHandleFavoriteCreate)
		}

		tagsGroup := v1Group.Group("/tags")
		{
			tagsGroup.GET("", HttpHandleTagsIndex)
		}

		postsGroup := v1Group.Group("/posts")
		{
			postsGroup.GET("/:namespace/:key", HttpHandlePostRead)
			postsGroup.DELETE("/:namespace/:key", HttpHandlePostDelete)
			postsGroup.GET("", HttpHandlePostIndex)
			postsGroup.POST("", HttpHandlePostCreate)
		}

		datasetGroup := v1Group.Group("/dataset")
		{
			datasetGroup.GET("/download/gob", HttpHandleDatasetDownloadGOB)
			datasetGroup.POST("/upload/gob", HttpHandleDatasetUploadGOB)
		}

		fetchGroup := v1Group.Group("/tasks")
		{
			fetchGroup.POST("/nurse-fetch", HttpHandleBooruFetch)
			fetchGroup.POST("/nurse-download", HttpHandleBooruDownload)
		}
	}

	router.GET("/", HttpHandleHealthCheck)
	router.GET("/_ah/health", HttpHandleHealthCheck)
}

func main() {
	defer amqpConnection.Close()
	defer amqpChannel.Close()

	appengine.Main()
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

func bindErrorResponse(err error) map[string][]string {
	errors := make(map[string][]string)
	switch msg := err.(type) {
	case validator.ValidationErrors:
		for _, v := range msg {
			errors[v.Field()] = append(errors[v.Field()], v.Tag())
		}
	default:
		errors["unknown"] = []string{fmt.Sprintln(err)}
	}
	return errors
}

func HttpHandleBooruFetch(c *gin.Context) {
	payload := NurseFetchTask{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": bindErrorResponse(err),
			"result": "error",
		})
		return
	}

	body, err := json.Marshal(payload)
	failOnError(err, "Failed to marshal a message")

	err = amqpChannel.Publish(
		"",                   // exchange
		nurseFetchQueue.Name, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: 2,
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf("[x] Sent %s", string(body))

	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

func HttpHandleBooruDownload(c *gin.Context) {
	payload := NurseDownloadTask{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": bindErrorResponse(err),
			"result": "error",
		})
		return
	}

	body, err := json.Marshal(payload)
	failOnError(err, "Failed to marshal a message")

	err = amqpChannel.Publish(
		"", // exchange
		nurseDownloadQueue.Name, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: 2,
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf("[x] Sent %s", string(body))

	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}
