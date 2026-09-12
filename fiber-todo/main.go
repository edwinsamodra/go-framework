package main

import (
	"database/sql"
	"errors"

	"example.com/todo-comparison/common"
	"github.com/gofiber/fiber/v2"
)

var store *common.Store

func main() {
	var err error
	store, err = common.OpenStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()
	// Fiber-specific: handlers use *fiber.Ctx, backed by fasthttp rather than net/http.
	app := fiber.New()
	app.Get("/health", health)
	app.Get("/todos", list)
	app.Post("/todos", create)
	app.Get("/todos/:id", get)
	app.Put("/todos/:id", update)
	app.Delete("/todos/:id", remove)
	app.Listen(":" + common.Env("PORT", "8083"))
}

func apiError(c *fiber.Ctx, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(404).JSON(fiber.Map{"error": "todo not found"})
	}
	return c.Status(500).JSON(fiber.Map{"error": err.Error()})
}
func health(c *fiber.Ctx) error {
	if err := store.Ping(c.Context()); err != nil {
		return c.Status(503).JSON(fiber.Map{"error": "database unavailable"})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}
func list(c *fiber.Ctx) error {
	todos, err := store.List(c.Context(), c.Query("user_id"))
	if err != nil {
		return apiError(c, err)
	}
	return c.JSON(todos)
}
func get(c *fiber.Ctx) error {
	// Fiber-specific: Params obtains route values captured as `:id`.
	todo, err := store.Get(c.Context(), c.Params("id"))
	if err != nil {
		return apiError(c, err)
	}
	return c.JSON(todo)
}
func input(c *fiber.Ctx) (common.TodoInput, error) {
	var value common.TodoInput
	// Fiber-specific: BodyParser performs JSON decoding from fasthttp's request body.
	if err := c.BodyParser(&value); err != nil {
		return value, err
	}
	return value, common.Validate(value)
}
func create(c *fiber.Ctx) error {
	value, err := input(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	todo, err := store.Create(c.Context(), value)
	if err != nil {
		return apiError(c, err)
	}
	return c.Status(201).JSON(todo)
}
func update(c *fiber.Ctx) error {
	value, err := input(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	todo, err := store.Update(c.Context(), c.Params("id"), value)
	if err != nil {
		return apiError(c, err)
	}
	return c.JSON(todo)
}
func remove(c *fiber.Ctx) error {
	if err := store.Delete(c.Context(), c.Params("id")); err != nil {
		return apiError(c, err)
	}
	return c.SendStatus(204)
}
