package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/middleware"
	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/service"

	"github.com/gorilla/mux"
)

func CreatePVZ(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(middleware.UserKey).(*middleware.Claims)
		if claims.Role != "moderator" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		var req struct{ City string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		pvz, err := svc.CreatePVZ(r.Context(), req.City)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(pvz)
	}
}

func ListPVZ(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(middleware.UserKey).(*middleware.Claims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		fromStr, toStr := r.URL.Query().Get("from"), r.URL.Query().Get("to")
		from, _ := time.Parse(time.RFC3339, fromStr)
		to, _ := time.Parse(time.RFC3339, toStr)
		list, err := svc.ListPVZ(r.Context(), from, to, 100, 0)
		if err != nil {
			http.Error(w, "error fetching", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(list)
	}
}

func StartSession(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(middleware.UserKey).(*middleware.Claims)
		if claims.Role != "pvz_employee" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		sess, err := svc.StartSession(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(sess)
	}
}

func AddItem(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(middleware.UserKey).(*middleware.Claims)
		if claims.Role != "pvz_employee" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		vars := mux.Vars(r)
		pvzID, _ := strconv.ParseInt(vars["id"], 10, 64)
		sid, _ := strconv.ParseInt(vars["sid"], 10, 64)
		var req struct{ Type string }
		json.NewDecoder(r.Body).Decode(&req)
		item, err := svc.AddItem(r.Context(), pvzID, sid, req.Type)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(item)
	}
}

func DeleteItem(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(middleware.UserKey).(*middleware.Claims)
		if claims.Role != "pvz_employee" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		vars := mux.Vars(r)
		pvzID, _ := strconv.ParseInt(vars["id"], 10, 64)
		sid, _ := strconv.ParseInt(vars["sid"], 10, 64)
		if err := svc.DeleteLastItem(r.Context(), pvzID, sid); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func CloseSession(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(middleware.UserKey).(*middleware.Claims)
		if claims.Role != "pvz_employee" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		vars := mux.Vars(r)
		pvzID, _ := strconv.ParseInt(vars["id"], 10, 64)
		sid, _ := strconv.ParseInt(vars["sid"], 10, 64)
		if err := svc.CloseSession(r.Context(), pvzID, sid); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
