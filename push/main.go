package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

type (
	PipelineCallback func(*TagPage, chan<- *TagPage)
)

var (
	uploadURL      = "http://localhost:8080/v1/dataset/upload/gob"
	healthCheckURL = "http://localhost:8080/_ah/health"

	waitPeriod time.Duration = 10
	basePath   string
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stderr)

	if URL, ok := os.LookupEnv("NURSETAGS_URL"); ok {
		uploadURL = fmt.Sprintf("%s/v1/dataset/upload/gob", URL)
		healthCheckURL = fmt.Sprintf("%s/_ah/health", URL)
	}

	home, err := homedir.Dir()
	if err != nil {
		log.WithError(err).Fatal()
	}
	basePath = filepath.Join(home, ".booru")
}

func main() {
	HealthCheck()

	stage1 := Pipeline(TagGenerator(), FetchFirstPageStage)
	stage2 := Pipeline(stage1, FetchAllPagesStage)
	stage3 := Pipeline(stage2, DatabaseInsertStage)

	// wait all jobs
	for range stage3 {
	}

	DatabaseUpload()
}

func HealthCheck() {
	checkURL := func() bool {
		_, err := http.Get(healthCheckURL)
		return err == nil
	}
	for ok := checkURL(); !ok; ok = checkURL() {
		log.Errorf("Server not online, waiting %d seconds\n", waitPeriod)
		time.Sleep(waitPeriod * time.Second)
	}
}
