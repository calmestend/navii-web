package main

import (
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/naviiapp/web/internal/router"
)

func main() {
	mux := router.Init()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found")
	}

	addr := ":" + port
	log.Infof("Starting server on port %s", port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

