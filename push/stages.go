package main

import (
	"bufio"
	"bytes"
	"net/http"
	"os"

	"github.com/hellupline/hello-nurse/nursetags"

	log "github.com/sirupsen/logrus"
)

func FetchFirstPageStage(tag *TagPage, out chan<- *TagPage) {
	tag.Logger.Info("fetch started")
	if err := tag.Fetch(); err != nil {
		tag.Logger.WithField("url", tag.URL()).WithError(err).Error("Failed to fetch TagPage")
		return
	}

	pages := tag.Pages()
	tag.Logger.WithFields(log.Fields{
		"count": tag.Count,
		"pages": pages,
	}).Info("fetch done")

	for i := 0; i < pages; i++ {
		out <- tag.NextPage(i)
	}
}

func FetchAllPagesStage(tag *TagPage, out chan<- *TagPage) {
	tag.Logger.Info("fetch started")
	if err := tag.Fetch(); err != nil {
		tag.Logger.WithField("url", tag.URL()).WithError(err).Error("Failed to fetch TagPage")
		return
	}
	tag.Logger.WithFields(log.Fields{
		"count": tag.Count,
		"pages": tag.Pages(),
	}).Info("fetch done")
	out <- tag
}

func DatabaseInsertStage(tag *TagPage, out chan<- *TagPage) {
	for _, post := range tag.Posts {
		nursetags.DefaultDatabase.PostCreate(nursetags.PostData{
			PostKey: nursetags.PostKey{
				Namespace: tag.Domain,
				ID:        post.Key,
			},
			Tags:  post.Tags(),
			Value: "https:" + post.URL,
		})
	}
	out <- tag
}

// input
func TagGenerator() chan *TagPage {
	// tags := ReadTagsFile()
	tags := []*TagPage{
		NewTagPage("konachan.net", "landscape"),
		NewTagPage("konachan.net", "moon"),
		NewTagPage("konachan.net", "night"),
		NewTagPage("konachan.net", "scenic"),
		NewTagPage("konachan.net", "sky"),
		NewTagPage("konachan.net", "star"),
		NewTagPage("konachan.net", "sunset"),
		NewTagPage("konachan.net", "ruins"),
	}

	out := make(chan *TagPage, len(tags))
	go func() {
		defer close(out)
		for _, tag := range tags {
			out <- tag
		}
	}()
	return out
}

func ReadTagsFile() []*TagPage {
	tags := make([]*TagPage, 0)

	f, err := os.Open("./tags.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // nolint: errcheck

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tags = append(tags, NewTagPage("konachan.net", scanner.Text()))
	}
	return tags
}

// output
func DatabaseUpload() {
	log.Info("Database push started")
	buffer := &bytes.Buffer{}
	if err := nursetags.DefaultDatabase.Write(buffer); err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", uploadURL, buffer)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := http.DefaultClient.Do(req); err != nil {
		log.Fatal(err)
	}
	log.Info("Database push done")
}