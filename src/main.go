package main

import (
	"embed"
	"log"
	"net/http"
)

//go:embed web/*
var content embed.FS

func main() {
	http.Handle("/", http.FileServer(http.FS(content)))
	log.Println("Serving on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
