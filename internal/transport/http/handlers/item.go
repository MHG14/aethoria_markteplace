package handlers

import (
	"github.com/MHG14/aethoria_marketplace/internal/application"
	"github.com/gofiber/fiber/v2"
)

type ItemHandler struct {
	app *application.App
}

func NewItemHandler(app *application.App) *ItemHandler {
	return &ItemHandler{app: app}
}

func (h *ItemHandler) Create(c *fiber.Ctx) error {
	var req application.CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	result, err := h.app.CreateItem(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *ItemHandler) List(c *fiber.Ctx) error {
	result, err := h.app.ListItems(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (h *ItemHandler) ListByOwner(c *fiber.Ctx) error {
	ownerID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	result, err := h.app.ListItemsByOwner(c.Context(), int64(ownerID))
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (h *ItemHandler) Get(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	result, err := h.app.GetItem(c.Context(), int64(id))
	if err != nil {
		return err
	}
	return c.JSON(result)
}
