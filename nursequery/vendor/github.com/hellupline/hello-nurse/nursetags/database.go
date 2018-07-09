package nursetags

import (
	"encoding/gob"
	"io"
	"sync"
)

type (
	Database struct { // nolint
		sync.RWMutex
		Tables
	}

	Tables struct { // nolint
		Posts PostsDB
		Tags  TagsDB
	}

	PostsDB  map[PostKey]PostData // nolint
	PostData struct {             // nolint
		PostKey

		Value map[string]string `json:"value"`

		Tags []string `json:"tags"`
		Type string   `json:"type"`
	}
	PostKey struct { // nolint
		Namespace string `json:"namespace"`
		Key       string `json:"key"`
	}

	TagsDB map[TagKey]Tag // nolint
	Tag    Set            // nolint
	TagKey string         // nolint

	TagCount struct { // nolint
		Name  TagKey `json:"name"`
		Count int    `json:"count"`
	}
)

func NewDatabase() *Database { // nolint
	return &Database{
		Tables: Tables{
			Posts: PostsDB{},
			Tags:  TagsDB{},
		},
	}
}

func (d *Database) Read(r io.Reader) error {
	var t Tables

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

	return gob.NewEncoder(w).Encode(d.Tables)
}

func (d *Database) PostIndex(query string) []PostData { // nolint
	posts := make([]PostData, 0)

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

func (d *Database) PostCreate(post PostData) { // nolint
	d.Lock()
	defer d.Unlock()

	for _, tagName := range post.Tags {
		tag, ok := d.Tags[TagKey(tagName)]
		if !ok {
			tag = Tag{}
			d.Tags[TagKey(tagName)] = tag
		}
		Set(tag).Add(post.PostKey)
	}
	d.Posts[post.PostKey] = post
}

func (d *Database) PostRead(key PostKey) (PostData, bool) { // nolint
	d.RLock()
	defer d.RUnlock()

	post, ok := d.Posts[key]
	return post, ok
}

func (d *Database) PostDelete(key PostKey) { // nolint
	d.Lock()
	defer d.Unlock()

	if post, ok := d.Posts[key]; ok {
		// remove post from tags
		for _, tagName := range post.Tags {
			tag := d.Tags[TagKey(tagName)]
			Set(tag).Remove(post.PostKey)

			// if tag is empty, delete it
			if len(tag) == 0 {
				delete(d.Tags, TagKey(tagName))
			}
		}

		delete(d.Posts, key)
	}
}

func (d *Database) TagIndex() []TagCount { // nolint
	tagCounts := make([]TagCount, 0, len(d.Tags))

	d.RLock()
	defer d.RUnlock()

	for key, value := range d.Tags {
		tagCounts = append(tagCounts, TagCount{key, len(value)})
	}
	return tagCounts
}

func (d *Database) TagRead(key string) (Tag, bool) { // nolint
	d.RLock()
	defer d.RUnlock()

	tag, ok := d.Tags[TagKey(key)]
	return tag, ok
}
