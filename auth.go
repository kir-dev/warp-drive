package main

import (
	"code.google.com/p/goauth2/oauth"
	"crypto"
	_ "crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type group struct {
	Id     int
	Name   string
	Status string
}

type apiResponse struct {
	EduPersonEntitlement []group
}

func oauthWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := getSession(w, r)
		// when logged in execute handler
		if s.isLoggedIn() {
			if isAuthorized(s) {
				h.ServeHTTP(w, r)
			} else {
				http.Error(w, "Nem vagy Kir-Dev tag!", http.StatusForbidden)
			}
			return
		}

		hash := crypto.SHA1.New()
		_, err := hash.Write([]byte(r.RemoteAddr + r.UserAgent()))
		if err != nil {
			log.Printf("Something really went wrong: %s", err)
			http.Error(w, "Ups...", http.StatusInternalServerError)
			return
		}

		ustate := fmt.Sprintf("%x", hash.Sum(nil))
		if s.setUserState(ustate) {
			http.Redirect(w, r, config.oauth().AuthCodeURL(ustate), http.StatusFound)
		}
	})
}

func isAuthorized(sess *session) bool {
	if config.GroupId <= 0 {
		// every logged in user can upload
		return true
	}

	res, err := http.Get("https://auth.sch.bme.hu/api/profile/?access_token=" + sess.accessToken())
	if err != nil {
		log.Printf("Error while getting profile information: %s", err)
	}

	defer res.Body.Close()
	apiRes := apiResponse{}
	if err = json.NewDecoder(res.Body).Decode(&apiRes); err != nil {
		log.Printf("Error decoding auth.sch api response: %s ", err)
		return false
	}

	for _, g := range apiRes.EduPersonEntitlement {
		if g.Id == config.GroupId {
			return true
		}
	}

	return false
}

func oauthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	t := &oauth.Transport{Config: config.oauth()}
	tok, err := t.Exchange(r.FormValue("code"))
	if err != nil {
		log.Printf("Error getting access token: %s", err)
		return
	}

	sess := getSession(w, r)

	ustate := sess.userState()
	if r.FormValue("state") != ustate {
		http.Error(w, "Wrong user state", http.StatusBadRequest)
		return
	}

	sess.setAccessToken(tok.AccessToken)
	http.Redirect(w, r, "/upload", http.StatusFound)
}
