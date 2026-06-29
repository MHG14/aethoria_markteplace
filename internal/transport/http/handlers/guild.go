package handlers

import (
	"github.com/MHG14/aethoria_marketplace/internal/application"
	"github.com/gofiber/fiber/v2"
)

type GuildHandler struct {
	app *application.App
}

func NewGuildHandler(app *application.App) *GuildHandler {
	return &GuildHandler{app: app}
}

func (h *GuildHandler) Create(c *fiber.Ctx) error {
	var req application.CreateGuildRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	result, err := h.app.CreateGuild(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *GuildHandler) Get(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	result, err := h.app.GetGuild(c.Context(), int64(id))
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (h *GuildHandler) GetWallet(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	result, err := h.app.GetWallet(c.Context(), int64(id))
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (h *GuildHandler) GetTransactions(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	result, err := h.app.GetTransactions(c.Context(), int64(id))
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (h *GuildHandler) TopUp(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	var req application.TopUpWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	req.GuildID = int64(id)
	result, err := h.app.TopUpWallet(c.Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(result)
}
