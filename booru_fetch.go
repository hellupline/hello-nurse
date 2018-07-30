package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/pool.v3"
)

func BooruGetTagPage(p pool.Pool, domain, name string, page int) pool.WorkFunc { // nolint
	return func(wu pool.WorkUnit) (interface{}, error) {
		t := NewBooruTag(domain, name, page)

		t.Logger.Info("fetch started")
		if err := t.Fetch(); err != nil {
			t.Logger.WithError(err).Error("Failed to fetch TagPage")
			return nil, nil
		}

		if wu.IsCancelled() {
			return nil, nil
		}

		pages := t.Pages()
		t.Logger.WithFields(log.Fields{
			"count": t.Count,
			"pages": pages,
		}).Info("fetch done")

		if page == 0 {
			for i := 1; i < pages; i++ {
				p.Queue(GetBooruTagPage(p, domain, name, i))
			}
		}

		t.SavePosts()
		return nil, nil
	}
}

func (t *BooruTag) SavePosts() { // nolint
	for _, p := range t.Posts {
		database.PostCreate(Post{
			PostKey: PostKey{t.Domain, p.Key},
			Value: map[string]string{
				"preview_url": p.PreviewURL,
				"sample_url":  p.SampleURL,
				"file_url":    p.FileURL,
			},
			Tags: p.Tags(),
		})
	}
}
