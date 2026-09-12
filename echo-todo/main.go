package main

import (
	"database/sql"
	"errors"
	"net/http"

	"example.com/todo-comparison/common"
	"github.com/labstack/echo/v4"
)

var store *common.Store

func main() {
	var err error
	store, err = common.OpenStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()
	// Echo-specific: methods register a path and return an `echo.Context` to handlers.
	e := echo.New()
	e.GET("/health", health)
	e.GET("/todos", list)
	e.POST("/todos", create)
	e.GET("/todos/:id", get)
	e.PUT("/todos/:id", update)
	e.DELETE("/todos/:id", remove)
	e.Logger.Fatal(e.Start(":" + common.Env("PORT", "8082")))
}

func apiError(c echo.Context, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(404, echo.Map{"error": "todo not found"})
	}
	return c.JSON(500, echo.Map{"error": err.Error()})
}
func health(c echo.Context) error {
	if err := store.Ping(c.Request().Context()); err != nil {
		return c.JSON(503, echo.Map{"error": "database unavailable"})
	}
	return c.JSON(200, echo.Map{"status": "ok"})
}
func list(c echo.Context) error {
	todos, err := store.List(c.Request().Context(), c.QueryParam("user_id"))
	if err != nil {
		return apiError(c, err)
	}
	return c.JSON(200, todos)
}
func get(c echo.Context) error {
	// Echo-specific: Param obtains values captured with the `:id` syntax.
	todo, err := store.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return apiError(c, err)
	}
	return c.JSON(200, todo)
}
func input(c echo.Context) (common.TodoInput, error) {
	var value common.TodoInput
	// Echo-specific: Bind decodes the request body into a struct.
	if err := c.Bind(&value); err != nil {
		return value, err
	}
	return value, common.Validate(value)
}
func create(c echo.Context) error {
	value, err := input(c)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}
	todo, err := store.Create(c.Request().Context(), value)
	if err != nil {
		return apiError(c, err)
	}
	return c.JSON(http.StatusCreated, todo)
}
func update(c echo.Context) error {
	value, err := input(c)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}
	todo, err := store.Update(c.Request().Context(), c.Param("id"), value)
	if err != nil {
		return apiError(c, err)
	}
	return c.JSON(200, todo)
}
func remove(c echo.Context) error {
	if err := store.Delete(c.Request().Context(), c.Param("id")); err != nil {
		return apiError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}
