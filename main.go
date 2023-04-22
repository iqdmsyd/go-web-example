package main

import (
	"fmt"
	"net/http"
)

func main() {
	// process dynamic requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Dek")
	})

	// server static assets
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// accept connections
	http.ListenAndServe(":8080", nil)
}
