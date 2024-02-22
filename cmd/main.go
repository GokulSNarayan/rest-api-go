package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", &homeHandler{})
	mux.Handle("/recipes", &RecipesHandler{})
	mux.Handle("/recipes/", &homeHandler{})

	http.ListenAndServe(":3000", mux)
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page!"))
}

type RecipesHandler struct {}

func (h *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my recipes page!"))
}