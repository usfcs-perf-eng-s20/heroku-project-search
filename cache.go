package main

import (
	"github.com/golang/groupcache/lru"
	"log"
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
	MaxEntries, _ := strconv.Atoi(os.Getenv("CACHE_MAX_ENTRIES"))
	log.Println("Using a cache with ", MaxEntries, " entries")
	return lru.New(MaxEntries)
}
