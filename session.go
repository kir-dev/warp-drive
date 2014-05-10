package main

import (
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

const (
	WarpSession        = "warp-session"
	UserAccessTokenKey = "access_token"
	UserOAuthStateKey  = "ustate"
)

var (
	sessionStore *sessions.CookieStore
)

func createSessionStore() {
	sessionStore = sessions.NewCookieStore([]byte(config.SessionSecret))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   config.Secure,
	}
}

type session struct {
	store sessions.Store
	req   *http.Request
	rw    http.ResponseWriter
}

func getSession(w http.ResponseWriter, r *http.Request) *session {
	return &session{sessionStore, r, w}
}

func (s *session) get() (*sessions.Session, bool) {
	session, err := sessionStore.Get(s.req, WarpSession)
	if err != nil {
		log.Printf("Could not get session: %s", err)
		return nil, false
	}

	return session, true
}

func (s *session) set(key, value string) bool {
	sess, ok := s.get()
	if !ok {
		return false
	}

	sess.Values[key] = value
	sess.Save(s.req, s.rw)
	return true
}

func (s *session) isLoggedIn() bool {
	session, ok := s.get()
	if !ok {
		return false
	}

	_, hasKey := session.Values[UserAccessTokenKey]
	return hasKey
}

func (s *session) setAccessToken(token string) bool {
	return s.set(UserAccessTokenKey, token)
}

func (s *session) accessToken() string {
	session, ok := s.get()
	if !ok {
		return ""
	}

	tok, _ := session.Values[UserAccessTokenKey].(string)
	return tok
}

func (s *session) setUserState(ustate string) bool {
	return s.set(UserOAuthStateKey, ustate)
}

// Get and delete user state from the session
func (s *session) userState() string {
	session, ok := s.get()
	if !ok {
		return ""
	}

	state, _ := session.Values[UserOAuthStateKey].(string)
	delete(session.Values, UserOAuthStateKey)
	session.Save(s.req, s.rw)
	return state
}
