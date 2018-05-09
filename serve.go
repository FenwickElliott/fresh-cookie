package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("port")
	if port == "" {
		port = "80"
	}
	http.HandleFunc("/", root)
	fmt.Println("Listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func root(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, I'm the root\n")
}
