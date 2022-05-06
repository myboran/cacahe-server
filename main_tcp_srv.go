package main

import (
	"mycache/httprest/cache"
	"mycache/httprest/http"
	"mycache/tcp/tcp"
)

func main() {
	ca := cache.New("inmemory")
	go tcp.New(ca).Listen()
	http.New(ca).Listen()
}
