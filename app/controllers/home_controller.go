package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/krishh1at/test/config"
)

func Home(w http.ResponseWriter, r *http.Request) {
	sid, err := r.Cookie("sessionId")

	if err != nil {
		http.Redirect(w, r, "/signup", http.StatusPermanentRedirect)
		return
	}

	xsid := strings.SplitN(sid.Value, "|", 2)
	if len(xsid) != 2 {
		http.Redirect(w, r, "/signup", http.StatusPermanentRedirect)
		return
	}

	mac1, err := base64.StdEncoding.DecodeString(xsid[0])
	if err != nil {
		log.Fatalln("unnable to decode sessionid")
	}

	mac := hmac.New(sha256.New, Key)
	mac.Write([]byte(xsid[1]))
	mac2 := mac.Sum(nil)

	if !hmac.Equal(mac1, mac2) {
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}

	config.Template.ExecuteTemplate(w, "home.html", nil)
}
