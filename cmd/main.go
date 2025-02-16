package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/atadzan/drm-go/pkg/handlers"
)

func main() {
	port := flag.String("port", "7777", "Server running port")
	flag.Parse()
	parsedPort := *port
	handler := handlers.New()
	http.HandleFunc("/cas-server", handler.CasServerHandler)
	http.HandleFunc("/user-key", handler.UserKeyHandler)

	log.Printf("DRM Server is running on port %s...", parsedPort)
	log.Fatal(http.ListenAndServe(":"+parsedPort, nil))
}
