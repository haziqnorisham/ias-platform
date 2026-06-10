package http

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type session struct {
	user    User
	expires time.Time
}

var (
	sessions    sync.Map
	sessionTTL  = 24 * time.Hour
	cookieName  = "session_id"
	authCapable = false
)

func InitSessionStore() {
	authCapable = true
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			sessions.Range(func(key, value interface{}) bool {
				if s, ok := value.(*session); ok && now.After(s.expires) {
					sessions.Delete(key)
				}
				return true
			})
		}
	}()
	slog.Info("Session store initialized", "ttl", sessionTTL)
}

func CreateSession(user User) (string, *http.Cookie) {
	id := generateSessionID()
	sessions.Store(id, &session{
		user:    user,
		expires: time.Now().Add(sessionTTL),
	})
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   int(sessionTTL.Seconds()),
	}
	return id, cookie
}

func ValidateSession(sessionID string) *User {
	val, ok := sessions.Load(sessionID)
	if !ok {
		return nil
	}
	s := val.(*session)
	if time.Now().After(s.expires) {
		sessions.Delete(sessionID)
		return nil
	}
	s.expires = time.Now().Add(sessionTTL)
	return &s.user
}

func RefreshSessionCookie(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   int(sessionTTL.Seconds()),
	}
}

func DeleteSession(sessionID string) {
	sessions.Delete(sessionID)
}

func authEnabled() bool {
	return authCapable && os.Getenv("AUTH_ENABLED") == "true"
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
