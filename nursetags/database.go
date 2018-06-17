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
		ID        string `json:"id" binding:"required"`
	}
	PostData struct {
		PostKey

		Tags  []string `json:"tags" binding:"required"`
		Value string   `json:"value" binding:"required"`
	}
)

var (
	database = Database{
		DataTables: DataTables{
			Favorites: FavoritesDB{},
			Tags:      TagsDB{},
			Posts:     PostsDB{},
		},
	}
)

func (d *Database) Read(r io.Reader) error {
	database.Lock()
	defer database.Unlock()

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
	database.RLock()
	defer database.RUnlock()

	return gob.NewEncoder(w).Encode(d.DataTables)
}

func DatabaseFavoritesQuery() []Favorite {
	favorites := make([]Favorite, 0)
	database.RLock()
	defer database.RUnlock()
	for _, favorite := range database.Favorites {
		favorites = append(favorites, favorite)
	}
	return favorites
}

func DatabaseFavoriteCreate(favorite Favorite) {
	database.Lock()
	defer database.Unlock()
	database.Favorites[favorite.Name] = favorite
}

func DatabaseFavoriteRead(key string) (Favorite, bool) {
	database.RLock()
	defer database.RUnlock()
	favorite, ok := database.Favorites[key]
	return favorite, ok
}

func DatabaseFavoriteDelete(key string) {
	database.Lock()
	defer database.Unlock()
	delete(database.Favorites, key)
}

func DatabaseTagsQuery() []string {
	keys := make([]string, 0, len(database.Tags))
	database.RLock()
	defer database.RUnlock()
	for key := range database.Tags {
		keys = append(keys, key)
	}
	return keys
}

func DatabaseTagRead(key string) (Tag, bool) {
	database.RLock()
	defer database.RUnlock()
	tag, ok := database.Tags[key]
	return tag, ok
}

func DatabasePostsQuery(query string) []PostData {
	database.RLock()
	defer database.RUnlock()

	posts := make([]PostData, 0)

	if len(query) == 0 {
		for _, post := range database.Posts {
			posts = append(posts, post)
		}
		return posts
	}

	keys, err := parseExpr(query)
	if err != nil {
		return posts
	}

	for postKey := range keys.Iter() {
		if post, ok := database.Posts[postKey]; ok {
			posts = append(posts, post)
		}
	}
	return posts

}

func DatabasePostCreate(post PostData) {
	database.Lock()
	defer database.Unlock()
	for _, tagName := range post.Tags {
		tag, ok := database.Tags[tagName]
		if !ok {
			tag = Tag{}
			database.Tags[tagName] = tag
		}
		Set(tag).Add(post.PostKey)
	}
	database.Posts[post.PostKey] = post
}

func DatabasePostRead(key PostKey) (PostData, bool) {
	database.RLock()
	defer database.RUnlock()
	post, ok := database.Posts[key]
	return post, ok
}

func DatabasePostDelete(key PostKey) {
	database.Lock()
	defer database.Unlock()
	if post, ok := database.Posts[key]; ok {
		// remove post from tags
		for _, tagName := range post.Tags {
			tag := database.Tags[tagName]
			Set(tag).Remove(post.PostKey)

			// if tag is empty, delete it
			if len(tag) == 0 {
				delete(database.Tags, tagName)
			}
		}

		delete(database.Posts, key)
	}
}
