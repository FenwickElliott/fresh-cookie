package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

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
	mongoServer := flag.String("mongoServer", "127.0.0.1", "Mongo Server Address")
	flag.Parse()

	session, err := mgo.Dial(*mongoServer)
	check(err)
	defer session.Close()
	db = session.DB("db")

	http.HandleFunc("/", root)
	http.HandleFunc("/in", in)
	http.HandleFunc("/find", find)
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
	if nativeID == nil {
		fmt.Println("Seting new native cookie")
		setNativeCookie(&w, r)
	} else {
		check(err)
	}

	err = checkFormValues(r, "cookie", "partner")
	if err != nil {
		fmt.Println(err)
		io.WriteString(w, err.Error()+"\n")
	} else {
		c := db.C(r.FormValue("partner"))
		err = c.Insert(association{nativeID.Value, r.FormValue("cookie")})
		check(err)
		res := association{}
		err = c.Find(bson.M{"nativeid": nativeID.Value}).One(&res)
		if err != nil {
			fmt.Println(err)
			io.WriteString(w, err.Error()+"\n")
		} else {
			fmt.Println("Added:", res)
			io.WriteString(w, "Added to database"+"\n")
		}
	}
}

func find(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	err = checkFormValues(r, "cookie", "partner")
	fmt.Println(r.FormValue("cookie"))
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		res := association{}
		c := db.C(r.FormValue("partner"))
		err = c.Find(bson.M{"foreignid": r.FormValue("cookie")}).One(&res)
		fmt.Println(res)
		if err != nil {
			io.WriteString(w, err.Error()+"\n")
		} else {
			io.WriteString(w, res.NativeID+"\n")
		}
	}
}
