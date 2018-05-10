package main

import (
	"encoding/csv"
	"flag"
	"os"

	mgo "gopkg.in/mgo.v2"
)

func init() {
	go getPartners()
	go setDB()
}

func getPartners() {
	partnerFile := "./partnerFile.csv"

	f, err := os.Open(partnerFile)
	check(err)
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	check(err)

	for _, line := range lines {
		scope := []string{}
		for i := 2; i < len(line); i++ {
			scope = append(scope, line[i])
		}
		p := partner{line[0], line[1], scope}
		partners = append(partners, p)
	}
}

func setDB() {
	mongoServer := flag.String("mongoServer", "127.0.0.1", "Mongo Server Address")
	flag.Parse()

	session, err := mgo.Dial(*mongoServer)
	check(err)
	defer session.Close()
	db = session.DB("db")
}
