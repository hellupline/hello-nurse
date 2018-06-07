package main

import (
	"github.com/deckarep/golang-set"
)

func databasePostsQuery(query interface{}) []Post {
	posts := make([]Post, 0)
	if query != nil {
		keys := parseQuery(query)
		for postKey := range keys.Iter() {
			if post, ok := postsDB[postKey.(string)]; ok {
				posts = append(posts, post)
			}
		}
	} else {
		for _, post := range postsDB {
			posts = append(posts, post)
		}
	}

	return posts
}

func databasePostCreate(post Post) {
	mux.Lock()
	defer mux.Unlock()
	for _, tagName := range post.Tags {
		tag, ok := tagsDB[tagName]
		if !ok {
			tag = mapset.NewSet()
			tagsDB[tagName] = tag
		}
		tag.Add(post.ID)
	}
	postsDB[post.ID] = post
}

func databasePostRead(key string) (Post, bool) {
	post, ok := postsDB[key]
	return post, ok
}

func databasePostDelete(key string) {
	mux.Lock()
	defer mux.Unlock()
	if post, ok := databasePostRead(key); ok {
		// remove post from tags
		for _, tagName := range post.Tags {
			tag := tagsDB[tagName]
			tag.Remove(post.ID)

			// if tag is empty, delete it
			if tag.Cardinality() == 0 {
				delete(tagsDB, tagName)
			}
		}

		delete(postsDB, key)
	}
}
