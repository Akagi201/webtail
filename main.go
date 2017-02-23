package main

import (
	"net/http"

	"github.com/Akagi201/light"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func main() {
	root := light.New()

	root.Get("/tail", handleTail)
	root.Get("/follow", websocket.Handler(handleFollow).ServeHTTP)
	root.Get("/", handleHome)

	log.Printf("HTTP listening at: %v", opts.ListenAddr)
	http.ListenAndServe(opts.ListenAddr, root)
}
