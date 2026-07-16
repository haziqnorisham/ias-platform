package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func AuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, `{"error":"username and password are required"}`, http.StatusBadRequest)
		return
	}

	slog.Info("Login attempt", "username", req.Username)

	user, err := LdapAuthenticate(req.Username, req.Password)
	if err != nil {
		slog.Warn("Login failed", "username", req.Username, "error", err)
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	_, cookie := CreateSession(*user)
	http.SetCookie(w, cookie)

	resp := map[string]string{
		"username": user.Username,
		"role":     user.Role,
	}
	jsonData, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func AuthLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie(cookieName)
	if err == nil {
		DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	jsonData, _ := json.Marshal(map[string]string{"message": "logged out"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func AuthSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		http.Error(w, `{"error":"not authenticated"}`, http.StatusUnauthorized)
		return
	}

	user := ValidateSession(cookie.Value)
	if user == nil {
		http.Error(w, `{"error":"not authenticated"}`, http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, RefreshSessionCookie(cookie.Value))

	resp := map[string]string{
		"username": user.Username,
		"role":     user.Role,
	}
	jsonData, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
