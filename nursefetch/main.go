package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type (
	PipelineCallback func(*TagPage, chan<- *TagPage)
	DomainURL        func(*TagPage) *url.URL

	NurseFetchTask struct {
		Domain string `json:"domain" binding:"required"`
		Tag    string `json:"tag" binding:"required"`
	}

	TagPage struct {
		Domain string `xml:"-"`
		Name   string `xml:"-"`
		Page   int    `xml:"-"`
		Limit  int    `xml:"-"`

		Count int         `xml:"count,attr"`
		Posts []*PostPage `xml:"post"`

		Logger *log.Entry `xml:"-"`
	}
	PostPage struct {
		URL     string `xml:"file_url,attr"`
		Key     string `xml:"id,attr"`
		RawTags string `xml:"tags,attr"`
	}
)

var (
	domainURLs = map[string]DomainURL{
		"konachan.net": KonachanUrlBuilder,
	}

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

// TagPage API

func NewTagPage(domain, name string) *TagPage {
	return &TagPage{
		Domain: domain,
		Name:   name,
		Page:   0,
		Limit:  100,

		Logger: log.WithFields(log.Fields{
			"domain": domain,
			"name":   name,
			"page":   0,
			"limit":  100,
		}),
	}

}

func (t *TagPage) Fetch() error {
	if t.Page < 1 {
		t.CreateCacheDir()
	}
	filename := t.Filename()
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		u := t.URL().String()
		resp, err := http.Get(u)
		if err != nil {
			return err
		}
		defer resp.Body.Close() // nolint: errcheck

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, body, 0660)
		if err != nil {
			return err
		}
	}

	return xml.Unmarshal(body, t)
}

func (t *TagPage) NextPage(page int) *TagPage {
	return &TagPage{
		Domain: t.Domain,
		Name:   t.Name,
		Page:   page,
		Limit:  t.Limit,

		Logger: log.WithFields(log.Fields{
			"domain": t.Domain,
			"name":   t.Name,
			"page":   page,
			"limit":  t.Limit,
		}),
	}

}

func (t *TagPage) Filename() string {
	basePath := "/data" // XXX: use envvar
	prefix := filepath.Join(basePath, "cache", t.Domain, t.Name)
	return filepath.Join(prefix, fmt.Sprintf("%04d.xml", t.Page))
}

func (t *TagPage) CreateCacheDir() {
	basePath := "/data" // XXX: use envvar
	prefix := filepath.Join(basePath, "cache", t.Domain, t.Name)
	os.MkdirAll(prefix, os.ModePerm)
}

func (t *TagPage) URL() *url.URL {
	return domainURLs[t.Domain](t)
}

func (t *TagPage) Pages() int {
	return (t.Count / t.Limit) + 1
}

func (p *PostPage) Tags() []string {
	return strings.Split(strings.TrimSpace(p.RawTags), " ")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func KonachanUrlBuilder(tag *TagPage) *url.URL {
	return &url.URL{
		RawQuery: url.Values{
			"limit": []string{strconv.Itoa(tag.Limit)},
			"page":  []string{strconv.Itoa(tag.Page)},
			"tags":  []string{tag.Name},
		}.Encode(),
		Scheme: "https",
		Host:   "konachan.net",
		Path:   "/post.xml",
	}
}
