package main

func databaseFavoritesQuery() []Favorite {
	favorites := make([]Favorite, 0)
	for _, favorite := range favoritesDB {
		favorites = append(favorites, favorite)
	}
	return favorites
}

func databaseFavoriteCreate(favorite Favorite) {
	mux.Lock()
	defer mux.Unlock()
	favoritesDB[favorite.Name] = favorite
}

func databaseFavoriteRead(key string) (Favorite, bool) {
	favorite, ok := favoritesDB[key]
	return favorite, ok
}

func databaseFavoriteDelete(key string) {
	mux.Lock()
	defer mux.Unlock()
	delete(favoritesDB, key)
}
