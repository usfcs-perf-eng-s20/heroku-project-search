package main

import (
	"github.com/golang/groupcache/lru"
	"os"
	"strconv"
)


var cache = createCache()

type cacheKey struct {
	path string
	argument string
}

func getCacheKey(path string, argument string) cacheKey {
	return cacheKey{
		path:     path,
		argument: argument,
	}
}

func createCache() *lru.Cache {
	MaxEntries, _ := strconv.Atoi(os.Getenv("ANALYTICS_URL"))
	return lru.New(MaxEntries)
}
