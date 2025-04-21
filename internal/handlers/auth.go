package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/config"
	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func DummyLogin(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Role string `json:"role"`
	}
	var rq req
	json.NewDecoder(r.Body).Decode(&rq)
	claims := jwt.MapClaims{"role": rq.Role, "exp": time.Now().Add(time.Hour).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ts, _ := token.SignedString([]byte(config.Load().JWTSecret))
	json.NewEncoder(w).Encode(map[string]string{"token": ts})
}

func Register(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Email, Password, Role string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if err := svc.CreateUser(r.Context(), req.Email, req.Password, req.Role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func Login(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Email, Password string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		user, err := svc.Authenticate(r.Context(), req.Email, req.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		claims := jwt.MapClaims{"user_id": user.ID, "role": user.Role, "exp": time.Now().Add(time.Hour).Unix()}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ts, _ := token.SignedString([]byte(config.Load().JWTSecret))
		json.NewEncoder(w).Encode(map[string]string{"token": ts})
	}
}
