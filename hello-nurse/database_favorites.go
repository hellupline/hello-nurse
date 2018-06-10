package main

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
