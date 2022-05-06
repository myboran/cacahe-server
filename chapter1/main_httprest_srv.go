package main

import (
	"mycache/httprest/cache"
	"mycache/httprest/http"
)

func main() {
	c := cache.New("inmemory")
	http.New(c).Listen()
}
