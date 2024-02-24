package main

import (
	"net/http"
)
var(
	RecipeRe = regexp.MustCompile(`^/recipes/*$`)
	RecipeReWithId = regexp.MustCompile(`^/recipes/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
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
switch {
case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.PAth):
	h.CreeateREcipe(w,r)
	return
case r.Method == http.MethodGet && RecipeRe.MatchString(r.URL.Path):
	h.ListRecipes(w,r)
	return
case r.Method == http.MethodGet && RecipeRe.MatchString(r.URL.Path):
	h.GetRecipe(w,r)
	return
case r.Method == http.MethodPut && RecipeRe.MatchString(r.URL.Path):
	h.UpdateRecipe(w,r)
	return
	case r.Method == http.MethodDelete && RecipeRe.MatchString(r.URL.PAth):
	h.DeleteRecipe(w,r)
	return
default:
	return
}

func (h *RecipesHandler) CreateRecipes(w http.ResponseWriter, r *http.Request){}
func (h *RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request){}
func (h *RecipesHandler) GetRecipes(w http.ResponseWriter, r *http.Request){}
func (h *RecipesHandler) UpdateRecipes(w http.ResponseWriter, r *http.Request){}
func (h *RecipesHandler) DeleteRecipes(w http.ResponseWriter, r *http.Request){}

