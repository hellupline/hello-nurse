package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

type (
	PipelineCallback func(*TagPage, chan<- *TagPage)
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

func SaveToQueryServer(tag *TagPage, out chan<- *TagPage) {
	for _, post := range tag.Posts {
		log.Info("Saving Post", post.Key)
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(post); err != nil {
			log.Warning("Failed to encode json", post.Key)
		}
		if _, err := http.Post(nurseQueryPostURI, "application/json", &buf); err != nil {
			log.Warning("Failed to push json", post.Key)
		}
	}
	out <- tag
}

func Pipeline(in <-chan *TagPage, cb PipelineCallback) <-chan *TagPage {
	out := make(chan *TagPage, 1)

	wg := sync.WaitGroup{}
	worker := func() {
		defer wg.Done()
		for tag := range in {
			cb(tag, out)
		}
	}
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go worker()
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}
