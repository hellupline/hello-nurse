package nursetags

import (
	"encoding/gob"
	"io"
	"sync"
)

type (
	DataTables struct {
		Favorites FavoritesDB
		Tags      TagsDB
		Posts     PostsDB
	}
	Database struct {
		DataTables
		sync.RWMutex
	}

	FavoritesDB map[string]Favorite
	Favorite    struct {
		Name string `json:"name" binding:"required"`
	}

	TagsDB map[string]Tag
	Tag    Set

	PostsDB map[PostKey]PostData
	PostKey struct {
		Namespace string `json:"namespace" binding:"required"`
		Key       string `json:"key" binding:"required"`
	}
	PostData struct {
		PostKey

		Tags []string `json:"tags" binding:"required"`

		Type  string `json:"type" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	TagCount struct {
		Count int    `json:"count"`
		Name  string `json:"name"`
	}
)

var (
	DefaultDatabase = Database{
		DataTables: DataTables{
			Favorites: FavoritesDB{},
			Tags:      TagsDB{},
			Posts:     PostsDB{},
		},
	}
)

func (d *Database) Read(r io.Reader) error {
	d.Lock()
	defer d.Unlock()

	var data DataTables
	err := gob.NewDecoder(r).Decode(&data)
	if err != nil {
		return err
	}

	d.Favorites = data.Favorites
	d.Tags = data.Tags
	d.Posts = data.Posts
	return nil
}

func (d *Database) Write(w io.Writer) error {
	d.RLock()
	defer d.RUnlock()

	return gob.NewEncoder(w).Encode(d.DataTables)
}

func (d *Database) FavoriteQuery() []Favorite {
	favorites := make([]Favorite, 0)

	d.RLock()
	defer d.RUnlock()

	for _, favorite := range d.Favorites {
		favorites = append(favorites, favorite)
	}
	return favorites
}

func (d *Database) FavoriteCreate(favorite Favorite) {
	d.Lock()
	defer d.Unlock()

	d.Favorites[favorite.Name] = favorite
}

func (d *Database) FavoriteRead(key string) (Favorite, bool) {
	d.RLock()
	defer d.RUnlock()

	favorite, ok := d.Favorites[key]
	return favorite, ok
}

func (d *Database) FavoriteDelete(key string) {
	d.Lock()
	defer d.Unlock()

	delete(d.Favorites, key)
}

func (d *Database) TagQuery() []TagCount {
	tagCounts := make([]TagCount, 0, len(d.Tags))

	d.RLock()
	defer d.RUnlock()

	for key, value := range d.Tags {
		tagCounts = append(tagCounts, TagCount{len(value), key})
	}
	return tagCounts
}

func (d *Database) TagRead(key string) (Tag, bool) {
	d.RLock()
	defer d.RUnlock()

	tag, ok := d.Tags[key]
	return tag, ok
}

func (d *Database) PostQuery(query string) []PostData {
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

func (d *Database) PostCreate(post PostData) {
	d.Lock()
	defer d.Unlock()

	for _, tagName := range post.Tags {
		tag, ok := d.Tags[tagName]
		if !ok {
			tag = Tag{}
			d.Tags[tagName] = tag
		}
		Set(tag).Add(post.PostKey)
	}
	d.Posts[post.PostKey] = post
}

func (d *Database) PostRead(key PostKey) (PostData, bool) {
	d.RLock()
	defer d.RUnlock()

	post, ok := d.Posts[key]
	return post, ok
}

func (d *Database) PostDelete(key PostKey) {
	d.Lock()
	defer d.Unlock()

	if post, ok := d.Posts[key]; ok {
		// remove post from tags
		for _, tagName := range post.Tags {
			tag := d.Tags[tagName]
			Set(tag).Remove(post.PostKey)

			// if tag is empty, delete it
			if len(tag) == 0 {
				delete(d.Tags, tagName)
			}
		}

		delete(d.Posts, key)
	}
}
