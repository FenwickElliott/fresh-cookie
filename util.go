package main

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

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
