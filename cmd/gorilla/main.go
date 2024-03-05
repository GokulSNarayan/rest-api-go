package main

import (
	"encoding/json"
	"net/http"

	"github.com/GokulSNarayan/rest-api-go/pkg/recipes"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
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


func (h RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe recipes.Recipe
	if err:= json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w,r)
		return
	}
	resourceId := slug.Make(recipe.Name)
	if err := h.store.Add(resourceId, recipe); err != nil {
		InternalServerErrorHandler(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resourceId))
}
func (h RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.store.List()
	if err != nil {
		InternalServerErrorHandler(w,r)
		return
	}
	jsonBytes, err := json.Marshal(recipes)
	if err != nil {
		InternalServerErrorHandler(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	recipe, err := h.store.Get(id)
	if err != nil {
		InternalServerErrorHandler(w,r)
		return
	}
	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var newRecipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&newRecipe); err != nil {
		InternalServerErrorHandler(w,r)
		return
	}

	if err := h.store.Update(id, newRecipe); err != nil {
		if err == recipes.ErrNotFound {
			NotFoundHandler(w,r)
			return
		}
		InternalServerErrorHandler(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"] 
	if err := h.store.Remove(id); err != nil {
		if err == recipes.ErrNotFound {
			NotFoundHandler(w,r)
			return
		}
		InternalServerErrorHandler(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Internal server error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Not found"))
}