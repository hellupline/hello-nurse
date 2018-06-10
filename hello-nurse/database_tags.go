package main

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
