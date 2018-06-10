package main

import (
	"github.com/deckarep/golang-set"
)

func parseQuery(queryRaw interface{}) mapset.Set {
	switch queryRaw.(type) {
	case map[string]interface{}:
		queryMap := queryRaw.(map[string]interface{})
		var items mapset.Set

		andKeys, andOk := queryMap["and"]
		orKeys, orOk := queryMap["or"]
		notKeys, notOk := queryMap["not"]

		if andOk && orOk {
			andItems := keysIntersect(andKeys.([]interface{})...)
			orItems := keysUnion(orKeys.([]interface{})...)

			items = andItems.Intersect(orItems)
		} else if andOk {
			items = keysIntersect(andKeys.([]interface{})...)
		} else if orOk {
			items = keysUnion(orKeys.([]interface{})...)
		}
		if notOk {
			for _, key := range notKeys.([]interface{}) {
				items = items.Difference(parseQuery(key))
			}
		}
		return items

	case string:
		if tag, ok := databaseTagRead(queryRaw.(string)); ok {
			return tag
		}
		// using "unsafe" because database already has a lock
		return mapset.NewThreadUnsafeSet()
	}
	// using "unsafe" because database already has a lock
	return mapset.NewThreadUnsafeSet()
}

func keysIntersect(keys ...interface{}) mapset.Set {
	first, others := keys[0], keys[1:]
	result := parseQuery(first)
	for _, key := range others {
		result = result.Intersect(parseQuery(key))
	}
	return result
}

func keysUnion(keys ...interface{}) mapset.Set {
	first, others := keys[0], keys[1:]
	result := parseQuery(first)
	for _, key := range others {
		result = result.Union(parseQuery(key))
	}
	return result
}
