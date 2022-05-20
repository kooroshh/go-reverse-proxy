package main

import (
	"github.com/patrickmn/go-cache"
)

var userDB *cache.Cache

func CacheInit() {
	userDB = cache.New(cache.NoExpiration, cache.NoExpiration)
}
