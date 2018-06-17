package nursetags

import (
	"sync"
)

type (
	Database struct {
		Favorites FavoritesDB
		Tags      TagsDB
		Posts     PostsDB

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
		Favorites: FavoritesDB{},
		Tags:      TagsDB{},
		Posts:     PostsDB{},
	}
)

func databaseFavoritesQuery() []Favorite {
	favorites := make([]Favorite, 0)
	database.RLock()
	defer database.RUnlock()
	for _, favorite := range database.Favorites {
		favorites = append(favorites, favorite)
	}
	return favorites
}

func databaseFavoriteCreate(favorite Favorite) {
	database.Lock()
	defer database.Unlock()
	database.Favorites[favorite.Name] = favorite
}

func databaseFavoriteRead(key string) (Favorite, bool) {
	database.RLock()
	defer database.RUnlock()
	favorite, ok := database.Favorites[key]
	return favorite, ok
}

func databaseFavoriteDelete(key string) {
	database.Lock()
	defer database.Unlock()
	delete(database.Favorites, key)
}

func databaseTagsQuery() []string {
	keys := make([]string, 0, len(database.Tags))
	database.RLock()
	defer database.RUnlock()
	for key := range database.Tags {
		keys = append(keys, key)
	}
	return keys
}

func databaseTagRead(key string) (Tag, bool) {
	database.RLock()
	defer database.RUnlock()
	tag, ok := database.Tags[key]
	return tag, ok
}

func databasePostsQuery(query string) []PostData {
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

func databasePostCreate(post PostData) {
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

func databasePostRead(key PostKey) (PostData, bool) {
	database.RLock()
	defer database.RUnlock()
	post, ok := database.Posts[key]
	return post, ok
}

func databasePostDelete(key PostKey) {
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
