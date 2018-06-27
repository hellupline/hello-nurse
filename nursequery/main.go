package main // import "github.com/hellupline/hello-nurse/nursequery"

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
)

var (
	amqpConnection *amqp.Connection
	amqpChannel    *amqp.Channel
	nurseQueue     amqp.Queue
)

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
			postsGroup.POST("", HttpHandlePostDelete)
		}

		datasetGroup := v1Group.Group("/dataset")
		{
			datasetGroup.GET("/download/gob", HttpHandleDatasetDownloadGOB)
			datasetGroup.POST("/upload/gob", HttpHandleDatasetUploadGOB)
		}

		fetchGroup := v1Group.Group("/tasks/nurse-fetch")
		{
			fetchGroup.POST("", HttpHandleBooruFetch)
		}
	}

	router.GET("/", HttpHandleHealthCheck)
	router.GET("/_ah/health", HttpHandleHealthCheck)
}

func main() {
	var err error

	amqpConnection, err = amqp.Dial("amqp://nurse:animaniacs@rabbitmq:5672/")
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

	appengine.Main()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
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
		"",              // exchange
		nurseQueue.Name, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf("[x] Sent %s", string(body))

	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}
