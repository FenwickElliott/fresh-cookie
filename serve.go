package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":7000", nil))
}

func root(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, I'm the root\n")
}
