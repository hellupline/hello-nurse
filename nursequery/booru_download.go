package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
)

type (
	BooruDownloadPipelineCallback func(*NurseDownloadTask, chan<- *NurseDownloadTask)
)

func DownloadFile(task *NurseDownloadTask, out chan<- *NurseDownloadTask) {
	base := filepath.Join(baseDir, "files", task.Domain)
	os.MkdirAll(base, 0755)
	path := filepath.Join(base, filepath.Base(task.URL))

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return
	}

	resp, err := http.Get(task.URL)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close() // nolint: errcheck

	// XXX: use io.copy
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	err = ioutil.WriteFile(path, body, 0644)
	if err != nil {
		log.Error(err)
		return
	}

	out <- task
}

func BooruDownloadPipeline(in <-chan *NurseDownloadTask, cb BooruDownloadPipelineCallback) <-chan *NurseDownloadTask {
	out := make(chan *NurseDownloadTask, 1)

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
