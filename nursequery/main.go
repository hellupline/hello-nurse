package main // import "github.com/hellupline/hello-nurse/nursequery"

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"

	"google.golang.org/appengine"
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

		fetchGroup := v1Group.Group("/worker/booru/fetch")
		{
			fetchGroup.POST("", HttpHandleBooruFetch)
		}
	}

	router.GET("/", HttpHandleHealthCheck)
	router.GET("/_ah/health", HttpHandleHealthCheck)
}

func main() {
	appengine.Main()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func HttpHandleBooruFetch(c *gin.Context) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "hello"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")

	c.String(200, "hello")
}
