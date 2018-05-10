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
	db       *mgo.Database
	c        *mgo.Collection
	partners []string
)

type doc struct {
	NativeID  string
	ForeignID string
}

func main() {
	partners = append(partners, "inception")
	port := flag.String("port", "80", "port")
	mongoServer := flag.String("mongoServer", "127.0.0.1", "Mongo Server Address")
	flag.Parse()

	session, err := mgo.Dial(*mongoServer)
	check(err)
	defer session.Close()
	db = session.DB("db")

	http.HandleFunc("/", root)
	fmt.Println("Listening on port:", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func root(w http.ResponseWriter, r *http.Request) {
	nativeID, err := r.Cookie("nativeID")
	if nativeID == nil {
		fmt.Println("Seting new native cookie")
		setNativeCookie(&w, r)
	} else {
		check(err)
	}
	io.WriteString(w, "Hello, I'm the root\n")
}

func setNativeCookie(w *http.ResponseWriter, r *http.Request) {
	h := sha1.New()
	h.Write([]byte(time.Now().String() + r.RemoteAddr))
	http.SetCookie(*w, &http.Cookie{Name: "nativeID", Value: hex.EncodeToString(h.Sum(nil)), Expires: time.Now().Add(365 * 24 * time.Hour)})
}

func check(err error) {
	if err != nil {
		log.Print(err)
	}
}
