package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	err      error
	c        *mgo.Collection
	partners []string
)

type doc struct {
	NativeID  string
	ForeignID string
	PartnerID string
}

func main() {
	partners = append(partners, "inception")
	port := flag.String("port", "80", "port")
	mongoServer := flag.String("mongoServer", "127.0.0.1", "Mongo Server Address")
	flag.Parse()

	session, err := mgo.Dial(*mongoServer)
	check(err)
	defer session.Close()
	c = session.DB("db").C("collection")

	http.HandleFunc("/", root)
	fmt.Println("Listening on port:", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func root(w http.ResponseWriter, r *http.Request) {
	nativeID, err := r.Cookie("nativeID")
	if nativeID == nil {
		fmt.Println("Seting new native cookie")
		http.SetCookie(w, &http.Cookie{Name: "nativeID", Value: newID(r), Expires: time.Now().Add(365 * 24 * time.Hour)})
	} else {
		check(err)
	}
	io.WriteString(w, "Hello, I'm the root\n")
}

func newID(r *http.Request) string {
	h := sha1.New()
	rand := time.Now().String()
	for _, c := range r.Cookies() {
		rand += c.Value
	}
	h.Write([]byte(rand))
	return hex.EncodeToString(h.Sum(nil))
}

func check(err error) {
	if err != nil {
		log.Print(err)
	}
}
