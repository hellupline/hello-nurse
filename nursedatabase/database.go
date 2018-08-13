package nursedatabase

import (
	"encoding/gob"
	"io"
	"sync"
)

type (
	Database struct { // nolint: golint
		sync.RWMutex
		Dataset
	}

	Dataset struct { // nolint: golint
		Posts PostStorage
		Tags  TagStorage
	}

	PostStorage map[PostKey]Post // nolint: golint

	TagStorage map[TagKey]Tag // nolint: golint

	Post struct { // nolint: golint
		PostKey

		Value map[string]string `json:"value"`

		Tags []string `json:"tags"`
	}

	PostKey struct { // nolint: golint
		Type string `json:"type"`
		Key  string `json:"key"`
	}

	Tag PostKeySet // nolint: golint

	TagKey string // nolint: golint

	TagCount struct { // nolint: golint
		Name  TagKey `json:"name"`
		Count int    `json:"count"`
	}
)

func NewDatabase() *Database { // nolint: golint
	return &Database{
		Dataset: Dataset{
			Posts: PostStorage{},
			Tags:  TagStorage{},
		},
	}
}

func (d *Database) Read(r io.Reader) error {
	var t Dataset

	d.Lock()
	defer d.Unlock()

	err := gob.NewDecoder(r).Decode(&t)
	if err != nil {
		return err
	}

	d.Tags = t.Tags
	d.Posts = t.Posts
	return nil
}

func (d *Database) Write(w io.Writer) error {
	d.RLock()
	defer d.RUnlock()

	return gob.NewEncoder(w).Encode(d.Dataset)
}

func (d *Database) PostIndex(query string) []Post { // nolint: golint
	posts := make([]Post, 0)

	d.RLock()
	defer d.RUnlock()

	if len(query) == 0 {
		for _, post := range d.Posts {
			posts = append(posts, post)
		}
		return posts
	}

	keys, err := d.ParseExpr(query)
	if err != nil {
		return posts
	}

	for postKey := range keys.Iter() {
		if post, ok := d.Posts[postKey]; ok {
			posts = append(posts, post)
		}
	}
	return posts

}

func (d *Database) PostCreate(post Post) { // nolint: golint
	d.Lock()
	defer d.Unlock()

	for _, tagName := range post.Tags {
		tag, ok := d.Tags[TagKey(tagName)]
		if !ok {
			tag = Tag{}
			d.Tags[TagKey(tagName)] = tag
		}
		PostKeySet(tag).Add(post.PostKey)
	}
	d.Posts[post.PostKey] = post
}

func (d *Database) PostRead(key PostKey) (Post, bool) { // nolint: golint
	d.RLock()
	defer d.RUnlock()

	post, ok := d.Posts[key]
	return post, ok
}

func (d *Database) PostDelete(key PostKey) { // nolint: golint
	d.Lock()
	defer d.Unlock()

	if post, ok := d.Posts[key]; ok {
		// remove post from tags
		for _, tagName := range post.Tags {
			tag := d.Tags[TagKey(tagName)]
			PostKeySet(tag).Remove(post.PostKey)

			// if tag is empty, delete it
			if len(tag) == 0 {
				delete(d.Tags, TagKey(tagName))
			}
		}

		delete(d.Posts, key)
	}
}

func (d *Database) TagIndex() []TagCount { // nolint: golint
	tagCounts := make([]TagCount, 0, len(d.Tags))

	d.RLock()
	defer d.RUnlock()

	for key, value := range d.Tags {
		tagCounts = append(tagCounts, TagCount{key, len(value)})
	}
	return tagCounts
}

func (d *Database) TagRead(key string) (Tag, bool) { // nolint: golint
	d.RLock()
	defer d.RUnlock()

	tag, ok := d.Tags[TagKey(key)]
	return tag, ok
}
