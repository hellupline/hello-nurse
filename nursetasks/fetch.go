package nursetasks

import (
	"github.com/hellupline/hello-nurse/nursedatabase"

	"github.com/hellupline/hello-nurse/booruapi"
)

func SaveBooruPosts(database *nursedatabase.Database, domain string, t booruapi.Tag) { // nolint: golint
	for _, p := range t.Posts() {
		database.PostCreate(nursedatabase.Post{
			PostKey: nursedatabase.PostKey{Type: domain, Key: p.Key()},
			Value:   p.Body(),
			Tags:    p.Tags(),
		})
	}
}
