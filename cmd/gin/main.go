package main

import (
	"net/http"

	"github.com/GokulSNarayan/rest-api-go/pkg/recipes"
	"github.com/gin-gonic/gin"
)

func main() {
    // Create Gin router
    router := gin.Default()

    // Instantiate recipe Handler and provide a data store implementation
    store := recipes.NewMemStore()
    recipesHandler := NewRecipesHandler(store)

    // Register Routes
    router.GET("/", homePage)
    router.GET("/recipes", recipesHandler.ListRecipes)
    router.POST("/recipes", recipesHandler.CreateRecipe)
    router.GET("/recipes/:id", recipesHandler.GetRecipe)
    router.PUT("/recipes/:id", recipesHandler.UpdateRecipe)
    router.DELETE("/recipes/:id", recipesHandler.DeleteRecipe)

    // Start the server
    router.Run(":3002")
}

func homePage(c *gin.Context) {
    c.String(http.StatusOK, "This is my home page")
}

type RecipesHandler struct {
    store recipeStore
}

type recipeStore interface {
    Add(name string, recipe recipes.Recipe) error
    Get(name string) (recipes.Recipe, error)
    List() (map[string]recipes.Recipe, error)
    Update(name string, recipe recipes.Recipe) error
    Remove(name string) error
}

func (h RecipesHandler) CreateRecipe(c *gin.Context){}
func (h RecipesHandler) ListRecipes(c *gin.Context){}
func (h RecipesHandler) GetRecipe(c *gin.Context){}
func (h RecipesHandler) UpdateRecipe(c *gin.Context){}
func (h RecipesHandler) DeleteRecipe(c *gin.Context){}


// NewRecipesHandler is a constructor for RecipesHandler
func NewRecipesHandler(s recipeStore) *RecipesHandler {
    return &RecipesHandler{
        store: s,
    }
}