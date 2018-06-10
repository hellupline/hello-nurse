package main

import (
	"fmt"
	"sync"

	"github.com/deckarep/golang-set"
	"github.com/go-playground/validator"
)

type (
	FavoritesDB map[string]Favorite
	Favorite    struct {
		Name string `json:"name" binding:"required"`
	}

	TagsDB map[string]Tag
	Tag    mapset.Set

	PostsDB map[string]Post
	Post    struct {
		Tags      []string `json:"tags" binding:"required"`
		Namespace string   `json:"namespace" binding:"required"`
		External  bool     `json:"external"`
		ID        string   `json:"id" binding:"required"`
		Value     string   `json:"value" binding:"required"`
	}

	Database struct {
		Favorites FavoritesDB
		Tags      TagsDB
		Posts     PostsDB

		sync.RWMutex
	}
)

var (
	database = Database{
		Favorites: FavoritesDB{},
		Tags:      TagsDB{},
		Posts:     PostsDB{},
	}
)

func bindErrorResponse(err error) map[string][]string {
	errors := make(map[string][]string)
	switch msg := err.(type) {
	case validator.ValidationErrors:
		for _, v := range msg {
			errors[v.Field()] = append(errors[v.Field()], v.Tag())
		}
	default:
		errors["unknown"] = []string{fmt.Sprintln(err)}
	}
	return errors

}

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
	keys := make([]string, len(database.Tags))
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

func databasePostsQuery(query interface{}) []Post {
	posts := make([]Post, 0)
	database.RLock()
	defer database.RUnlock()
	if query != nil {
		keys := parseQuery(query)
		for postKey := range keys.Iter() {
			if post, ok := database.Posts[postKey.(string)]; ok {
				posts = append(posts, post)
			}
		}
	} else {
		for _, post := range database.Posts {
			posts = append(posts, post)
		}
	}

	return posts
}

func databasePostCreate(post Post) {
	database.Lock()
	defer database.Unlock()
	for _, tagName := range post.Tags {
		tag, ok := database.Tags[tagName]
		if !ok {
			// using "unsafe" because database already has a lock
			tag = mapset.NewThreadUnsafeSet()
			database.Tags[tagName] = tag
		}
		tag.Add(post.ID)
	}
	database.Posts[post.ID] = post
}

func databasePostRead(key string) (Post, bool) {
	database.RLock()
	defer database.RUnlock()
	post, ok := database.Posts[key]
	return post, ok
}

func databasePostDelete(key string) {
	database.Lock()
	defer database.Unlock()
	if post, ok := database.Posts[key]; ok {
		// remove post from tags
		for _, tagName := range post.Tags {
			tag := database.Tags[tagName]
			tag.Remove(post.ID)

			// if tag is empty, delete it
			if tag.Cardinality() == 0 {
				delete(database.Tags, tagName)
			}
		}

		delete(database.Posts, key)
	}
}
