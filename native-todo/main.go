package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/edwinsamodra/go-framework/common"
)

var store *common.Store

func main() {
	var err error
	store, err = common.OpenStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()
	// Go 1.22+ native routing: method/path patterns and `{id}` are provided by http.ServeMux.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /todos", list)
	mux.HandleFunc("POST /todos", create)
	mux.HandleFunc("GET /todos/{id}", get)
	mux.HandleFunc("PUT /todos/{id}", update)
	mux.HandleFunc("DELETE /todos/{id}", remove)
	http.ListenAndServe(":"+common.Env("PORT", "8085"), mux)
}

func jsonResponse(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
func apiError(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		jsonResponse(w, 404, map[string]string{"error": "todo not found"})
		return
	}
	jsonResponse(w, 500, map[string]string{"error": err.Error()})
}
func health(w http.ResponseWriter, r *http.Request) {
	if err := store.Ping(r.Context()); err != nil {
		jsonResponse(w, 503, map[string]string{"error": "database unavailable"})
		return
	}
	jsonResponse(w, 200, map[string]string{"status": "ok"})
}
func list(w http.ResponseWriter, r *http.Request) {
	todos, err := store.List(r.Context(), r.URL.Query().Get("user_id"))
	if err != nil {
		apiError(w, err)
		return
	}
	jsonResponse(w, 200, todos)
}
func get(w http.ResponseWriter, r *http.Request) {
	// Native net/http-specific: PathValue reads `{id}` from ServeMux's matched pattern.
	todo, err := store.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		apiError(w, err)
		return
	}
	jsonResponse(w, 200, todo)
}
func input(r *http.Request) (common.TodoInput, error) {
	var value common.TodoInput
	// Native net/http-specific: JSON decoding is explicit; there is no binding helper.
	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		return value, err
	}
	return value, common.Validate(value)
}
func create(w http.ResponseWriter, r *http.Request) {
	value, err := input(r)
	if err != nil {
		jsonResponse(w, 400, map[string]string{"error": err.Error()})
		return
	}
	todo, err := store.Create(r.Context(), value)
	if err != nil {
		apiError(w, err)
		return
	}
	jsonResponse(w, 201, todo)
}
func update(w http.ResponseWriter, r *http.Request) {
	value, err := input(r)
	if err != nil {
		jsonResponse(w, 400, map[string]string{"error": err.Error()})
		return
	}
	todo, err := store.Update(r.Context(), r.PathValue("id"), value)
	if err != nil {
		apiError(w, err)
		return
	}
	jsonResponse(w, 200, todo)
}
func remove(w http.ResponseWriter, r *http.Request) {
	if err := store.Delete(r.Context(), r.PathValue("id")); err != nil {
		apiError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
