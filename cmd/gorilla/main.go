package main

import (
	"net/http"

	"github.com/GokulSNarayan/rest-api-go/pkg/recipes"
	"github.com/gorilla/mux"
)

type homeHandler struct{}
type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	List() (map[string]recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	Delete(name string) error
}

type RecipesHandler struct {
	store recipeStore
}

func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{store: s}
}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

func main() {
	router := mux.NewRouter()

	home := &homeHandler{}

	router.Handle(("/"), home)

	http.ListenAndServe(":3001", router)
}
