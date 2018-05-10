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
	"gopkg.in/mgo.v2/bson"
)

var (
	err      error
	db       *mgo.Database
	partners []partner
)

type association struct {
	NativeID  string
	ForeignID string
}

type partner struct {
	PartnerID  string
	AuthHeader string
	Scope      []string
}

func main() {
	port := flag.String("port", "80", "port")
	flag.Parse()

	http.HandleFunc("/", root)
	http.HandleFunc("/in", in)
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

func in(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nativeID, err := r.Cookie("nativeID")
	if nativeID != nil {
		check(err)
	}

	c := db.C(r.FormValue("partner"))
	err = c.Insert(association{nativeID.Value, r.FormValue("cookie")})
	check(err)

	err = c.Find(bson.M{"nativeid": nativeID.Value}).One(&association{})
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, "Added to database")
	}
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
