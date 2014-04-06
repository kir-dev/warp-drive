package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func basicAuth(realm string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failed := false
		auth := r.Header.Get("Authorization")
		if auth != "" {
			digest := strings.TrimPrefix(auth, "Basic ")
			decoded, err := base64.StdEncoding.DecodeString(digest)

			if err != nil {
				log.Printf("Error decoding basic auth header: %s", err)

			} else {
				parts := strings.SplitN(string(decoded), ":", 2)

				if parts[0] == config.User && parts[1] == config.Password {
					log.Printf("Successful authentication from %s", r.RemoteAddr)
					h.ServeHTTP(w, r)
					return
				}
			}
			failed = true
		}

		if failed {
			log.Printf("Failed authentication from %s", r.RemoteAddr)
		}

		w.Header().Add("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", realm))
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
	})
}
