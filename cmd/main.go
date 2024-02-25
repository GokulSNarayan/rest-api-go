package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/GokulSNarayan/rest-api-go/pkg/recipes"
	"github.com/gosimple/slug"
)

var (
	RecipeRe       = regexp.MustCompile(`^/recipes/*$`)
	RecipeReWithId = regexp.MustCompile(`^/recipes/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	List() (map[string]recipes.Recipe, error)
	Remove(name string) error
}

func main() {
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)

	mux := http.NewServeMux()

	mux.Handle("/", &homeHandler{})
	mux.Handle("/recipes", recipesHandler)
	mux.Handle("/recipes/", recipesHandler)

	http.ListenAndServe(":3000", mux)
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqDump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("REQUEST:\n%s", string(reqDump))
	w.Write([]byte("This is my home page!"))

}

func (h *RecipesHandler) CreateRecipes(w http.ResponseWriter, r *http.Request) {
	var recipe recipes.Recipe
	// Decode the request body into the recipe struct
	if err := json.NewDecoder((r.Body)).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w,r)
		return
	}
	// Use the store to add the recipe
	resourceId := slug.Make(recipe.Name)
	if err := h.store.Add(resourceId, recipe); err != nil {
		InternalServerErrorHandler(w,r)
		return
	}
	// Write the response
	w.WriteHeader(http.StatusCreated)
}
func (h *RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request)   {
	recipes, err := h.store.List()

	if err != nil{
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
func (h *RecipesHandler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	matches := RecipeReWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w,r)
		return
	}
	recipe, err := h.store.Get(matches[1])
	if err != nil {
		if err == recipes.ErrNotFound {
			NotFoundHandler(w,r)
			return
		}
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
func (h *RecipesHandler) UpdateRecipes(w http.ResponseWriter, r *http.Request) {
	matches := RecipeReWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {	
		InternalServerErrorHandler(w,r)
		return
	}
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w,r)
		return
	}

	if err := h.store.Update(matches[1],recipe); err != nil {
		if err == recipes.ErrNotFound {
			NotFoundHandler(w,r)
			return
		}
		InternalServerErrorHandler(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)

}
func (h *RecipesHandler) DeleteRecipes(w http.ResponseWriter, r *http.Request) {
	matches := RecipeReWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w,r)
		return
	}
	if err := h.store.Remove(matches[1]); err != nil {
		if err == recipes.ErrNotFound {
			NotFoundHandler(w,r)
			return
		}
		InternalServerErrorHandler(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type RecipesHandler struct {
	store recipeStore
}

func NewRecipesHandler(store recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: store,
	}
}

func (h *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.CreateRecipes(w, r)
		return
	case r.Method == http.MethodGet && RecipeRe.MatchString(r.URL.Path):
		h.ListRecipes(w, r)
		return
	case r.Method == http.MethodGet && RecipeReWithId.MatchString(r.URL.Path):
		h.GetRecipes(w, r)
		return
	case r.Method == http.MethodPut && RecipeReWithId.MatchString(r.URL.Path):
		h.UpdateRecipes(w, r)
		return
	case r.Method == http.MethodDelete && RecipeReWithId.MatchString(r.URL.Path):
		h.DeleteRecipes(w, r)
		return
	default:
		return
	}
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Not Found"))
}