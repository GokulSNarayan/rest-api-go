package main

import (
	"net/http"

	"github.com/GokulSNarayan/rest-api-go/pkg/recipes"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
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

func (h RecipesHandler) CreateRecipe(c *gin.Context) {
	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := slug.Make(recipe.Name)
	if err := h.store.Add(id, recipe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
func (h RecipesHandler) ListRecipes(c *gin.Context) {
	recipeList, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipeList)
}
func (h RecipesHandler) GetRecipe(c *gin.Context) {
	id := c.Param("id")

	recipe, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}
func (h RecipesHandler) UpdateRecipe(c *gin.Context) {
	var newRecipe recipes.Recipe
	if err := c.ShouldBindJSON(&newRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	err := h.store.Update(id, newRecipe)
	if err != nil {
		if err == recipes.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h RecipesHandler) DeleteRecipe(c *gin.Context) {
	id := c.Param("id")
	err := h.store.Remove(id)
	if err != nil {
		if err == recipes.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// NewRecipesHandler is a constructor for RecipesHandler
func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}
