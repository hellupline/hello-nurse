package main

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type (
	PipelineCallback func(*NurseDownloadTask, chan<- *NurseDownloadTask)
)

func DownloadFile(task *NurseDownloadTask, out chan<- *NurseDownloadTask) {
	log.Info(task)
	out <- task
}

func Pipeline(in <-chan *NurseDownloadTask, cb PipelineCallback) <-chan *NurseDownloadTask {
	out := make(chan *NurseDownloadTask, 1)

	wg := sync.WaitGroup{}
	worker := func() {
		defer wg.Done()
		for task := range in {
			cb(task, out)
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
