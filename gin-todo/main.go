package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/edwinsamodra/go-framework/common"
	"github.com/gin-gonic/gin"
)

var store *common.Store

func main() {
	var err error
	store, err = common.OpenStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()
	// Gin-specific: Default creates an Engine with Logger and Recovery middleware.
	router := gin.Default()
	router.GET("/health", health)
	router.GET("/todos", list)
	router.POST("/todos", create)
	router.GET("/todos/:id", get)
	router.PUT("/todos/:id", update)
	router.DELETE("/todos/:id", remove)
	router.Run(":" + common.Env("PORT", "8084"))
}

func apiError(c *gin.Context, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, gin.H{"error": "todo not found"})
		return
	}
	c.JSON(500, gin.H{"error": err.Error()})
}
func health(c *gin.Context) {
	if err := store.Ping(c.Request.Context()); err != nil {
		c.JSON(503, gin.H{"error": "database unavailable"})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}
func list(c *gin.Context) {
	todos, err := store.List(c.Request.Context(), c.Query("user_id"))
	if err != nil {
		apiError(c, err)
		return
	}
	c.JSON(200, todos)
}
func get(c *gin.Context) {
	// Gin-specific: Param reads a value captured by Gin's `:id` route syntax.
	todo, err := store.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		apiError(c, err)
		return
	}
	c.JSON(200, todo)
}
func input(c *gin.Context) (common.TodoInput, error) {
	var value common.TodoInput
	// Gin-specific: ShouldBindJSON uses Gin's JSON binding pipeline.
	if err := c.ShouldBindJSON(&value); err != nil {
		return value, err
	}
	return value, common.Validate(value)
}
func create(c *gin.Context) {
	value, err := input(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	todo, err := store.Create(c.Request.Context(), value)
	if err != nil {
		apiError(c, err)
		return
	}
	c.JSON(http.StatusCreated, todo)
}
func update(c *gin.Context) {
	value, err := input(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	todo, err := store.Update(c.Request.Context(), c.Param("id"), value)
	if err != nil {
		apiError(c, err)
		return
	}
	c.JSON(200, todo)
}
func remove(c *gin.Context) {
	if err := store.Delete(c.Request.Context(), c.Param("id")); err != nil {
		apiError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
