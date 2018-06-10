package main

import (
	"github.com/deckarep/golang-set"
)

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
