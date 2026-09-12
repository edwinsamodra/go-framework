package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"example.com/todo-comparison/common"
	"github.com/go-chi/chi/v5"
)

var store *common.Store

func main() {
	var err error
	store, err = common.OpenStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()

	// Chi-specific: a Router is an http.Handler with expressive method/path helpers.
	router := chi.NewRouter()
	router.Get("/health", health)
	router.Get("/todos", list)
	router.Post("/todos", create)
	router.Get("/todos/{id}", get)
	router.Put("/todos/{id}", update)
	router.Delete("/todos/{id}", remove)

	http.ListenAndServe(":"+common.Env("PORT", "8081"), router)
}

func jsonResponse(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
func errorResponse(w http.ResponseWriter, err error) {
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
		errorResponse(w, err)
		return
	}
	jsonResponse(w, 200, todos)
}

func get(w http.ResponseWriter, r *http.Request) {
	// Chi-specific: URLParam reads `{id}` from Chi's route context.
	todo, err := store.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		errorResponse(w, err)
		return
	}
	jsonResponse(w, 200, todo)
}

func decode(r *http.Request) (common.TodoInput, error) {
	var input common.TodoInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err == nil {
		err = common.Validate(input)
	}
	return input, err
}
func create(w http.ResponseWriter, r *http.Request) {
	input, err := decode(r)
	if err != nil {
		jsonResponse(w, 400, map[string]string{"error": err.Error()})
		return
	}
	todo, err := store.Create(r.Context(), input)
	if err != nil {
		errorResponse(w, err)
		return
	}
	jsonResponse(w, 201, todo)
}
func update(w http.ResponseWriter, r *http.Request) {
	input, err := decode(r)
	if err != nil {
		jsonResponse(w, 400, map[string]string{"error": err.Error()})
		return
	}
	todo, err := store.Update(r.Context(), chi.URLParam(r, "id"), input)
	if err != nil {
		errorResponse(w, err)
		return
	}
	jsonResponse(w, 200, todo)
}
func remove(w http.ResponseWriter, r *http.Request) {
	if err := store.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(204)
}
