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
	Remove(name string) error
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
	store := recipes.NewMemStore()
    recipesHandler := NewRecipesHandler(store)
	router := mux.NewRouter()

	home := &homeHandler{}

	router.HandleFunc("/", home.ServeHTTP)
    router.HandleFunc("/recipes", recipesHandler.ListRecipes).Methods("GET")
    router.HandleFunc("/recipes", recipesHandler.CreateRecipe).Methods("POST")
    router.HandleFunc("/recipes/{id}", recipesHandler.GetRecipe).Methods("GET")
    router.HandleFunc("/recipes/{id}", recipesHandler.UpdateRecipe).Methods("PUT")
    router.HandleFunc("/recipes/{id}", recipesHandler.DeleteRecipe).Methods("DELETE")

	http.ListenAndServe(":3001", router)
}


func (h RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {}
func (h RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {}
func (h RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {}
func (h RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {}
func (h RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {}