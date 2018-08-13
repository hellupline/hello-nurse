package nurseworkers

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/pool.v3"

	"github.com/hellupline/hello-nurse/nursedatabase"

	"github.com/hellupline/hello-nurse/booru"
)

type (
	TaskManager struct { // nolint: golint
		database *nursedatabase.Database
		pool     pool.Pool
	}
)

func NewTaskManager(db *nursedatabase.Database, p pool.Pool) *TaskManager { // nolint: golint
	return &TaskManager{database: db, pool: p}

}

func (tm *TaskManager) BooruGetTagPage(c booru.Client, name string, page int) pool.WorkUnit { // nolint: golint
	return tm.pool.Queue(func(wu pool.WorkUnit) (interface{}, error) {
		t := c.NewTag(name, page)

		logger := log.WithFields(log.Fields{"name": name, "page": page})
		logger.Info("fetch started")
		if err := t.Fetch(); err != nil {
			logger.WithError(err).Error("Failed to fetch TagPage")
			return nil, nil
		}

		pages := t.Pages()
		logger.WithFields(log.Fields{
			// "count": t.Count,
			"pages": pages,
		}).Info("fetch done")

		// if I am the first one, I manage it
		if page == 0 {
			for i := 1; i < pages; i++ {
				tm.BooruGetTagPage(c, name, i)
			}
		}

		SaveBooruPosts(tm.database, c.Name(), t)
		return nil, nil
	})
}

func (tm *TaskManager) BooruGetFile(baseDir, typeKey, key string) pool.WorkUnit { // nolint: golint
	return tm.pool.Queue(func(wu pool.WorkUnit) (interface{}, error) {
		post, ok := tm.database.PostRead(nursedatabase.PostKey{Type: typeKey, Key: key})
		if !ok {
			log.WithFields(log.Fields{"type": typeKey, "key": key}).Warning("does not exists")
			return nil, nil
		}

		_ = DownloadPostFile(post, baseDir)
		return nil, nil
	})
}
