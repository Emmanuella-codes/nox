package controllers

import (
	"github.com/emmanuella-codes/nox/guest/dtos"
	"github.com/emmanuella-codes/nox/guest/messages"
	"github.com/emmanuella-codes/nox/guest/pipes"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ pipe *pipes.GuestPipe }

func NewController(pipe *pipes.GuestPipe) *Controller { return &Controller{pipe: pipe} }
func fail(c *fiber.Ctx, msg shared.PipeMessage) error {
	return c.Status(fiber.StatusBadRequest).JSON(shared.PipeRes[any]{Success: false, Message: msg})
}
func (c *Controller) Configure(ctx *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return fail(ctx, messages.Forbidden)
	}
	e, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return fail(ctx, messages.Invalid)
	}
	var d dtos.ConfigureDTO
	if ctx.BodyParser(&d) != nil {
		return fail(ctx, messages.Invalid)
	}
	r := c.pipe.Configure(ctx.Context(), u, e, d)
	if !r.Success {
		return fail(ctx, r.Message)
	}
	return ctx.Status(201).JSON(r)
}
func (c *Controller) Quota(ctx *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return fail(ctx, messages.Forbidden)
	}
	e, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return fail(ctx, messages.Invalid)
	}
	var d dtos.QuotaDTO
	if ctx.BodyParser(&d) != nil {
		return fail(ctx, messages.Invalid)
	}
	r := c.pipe.Quota(ctx.Context(), u, e, d)
	if !r.Success {
		return fail(ctx, r.Message)
	}
	return ctx.Status(201).JSON(r)
}
func (c *Controller) Add(ctx *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return fail(ctx, messages.Forbidden)
	}
	e, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return fail(ctx, messages.Invalid)
	}
	q, err := uuid.Parse(ctx.Params("quotaID"))
	if err != nil {
		return fail(ctx, messages.Invalid)
	}
	var d dtos.GuestDTO
	if ctx.BodyParser(&d) != nil {
		return fail(ctx, messages.Invalid)
	}
	r := c.pipe.Add(ctx.Context(), u, e, q, d)
	if !r.Success {
		return fail(ctx, r.Message)
	}
	return ctx.Status(201).JSON(r)
}
func (c *Controller) Claim(ctx *fiber.Ctx) error {
	var u *uuid.UUID
	if id, ok := middleware.CurrentUserID(ctx); ok {
		u = &id
	}
	r := c.pipe.Claim(ctx.Context(), ctx.Query("token"), u)
	if !r.Success {
		return fail(ctx, r.Message)
	}
	return ctx.Status(200).JSON(r)
}
func (c *Controller) CheckIn(ctx *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return fail(ctx, messages.Forbidden)
	}
	e, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return fail(ctx, messages.Invalid)
	}
	var d struct {
		Credential string `json:"credential"`
	}
	if ctx.BodyParser(&d) != nil {
		return fail(ctx, messages.Invalid)
	}
	r := c.pipe.CheckIn(ctx.Context(), u, e, d.Credential)
	if !r.Success {
		return fail(ctx, r.Message)
	}
	return ctx.Status(200).JSON(r)
}
