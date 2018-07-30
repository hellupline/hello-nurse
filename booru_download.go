package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/pool.v3"
)

func BooruGetFile(p pool.Pool, type_, key string) pool.WorkFunc { // nolint
	return func(wu pool.WorkUnit) (interface{}, error) {
		post, ok := database.PostRead(PostKey{type_, key})
		logger := log.WithFields(log.Fields{
			"type": type_,
			"key":  key,
		})

		if !ok {
			logger.Warning("does not exists")
			return nil, nil
		}

		if wu.IsCancelled() {
			return nil, nil
		}

		_ = BooruDownloadPost(post)
		return nil, nil
	}
}

func BooruDownloadPost(p Post) error { // nolint
	logger := log.WithFields(log.Fields{
		"type": p.Type,
		"key":  p.Key,
	})

	url, ok := p.Value["file_url"]
	if !ok {
		logger.Warning("does not contain key file_url")
		return nil
	}

	base := filepath.Join(baseDir, "files", p.Type)
	if err := os.MkdirAll(base, 0755); err != nil {
		logger.Error(err)
		return err
	}

	filename := filepath.Join(base, filepath.Base(url))
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		logger.Infof("file %s exists", filename)
		return nil // file exists, abort
	}

	logger.WithFields(log.Fields{"url": url}).Info("downloading")
	resp, err := http.Get("https:" + url)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer resp.Body.Close() // nolint

	// XXX: use io.copy
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.WithFields(log.Fields{"url": url}).Info("done")

	return nil
}
