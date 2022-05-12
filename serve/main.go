package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Response[T any] struct {
	Payload    T `json:"payload"`
	StatusCode int     `json:"status_code"`
}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Response[string]{
			Payload: "Welcome!",
			StatusCode: 200,
		})
	})

	http.ListenAndServe(":3137", r)
}
