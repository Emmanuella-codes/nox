package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/staff/dtos"
	"github.com/emmanuella-codes/nox/staff/messages"
	"github.com/emmanuella-codes/nox/staff/pipes"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ pipe *pipes.StaffPipe }

func NewController(pipe *pipes.StaffPipe) *Controller { return &Controller{pipe: pipe} }
func fail(c *fiber.Ctx, m shared.PipeMessage) error {
	return c.Status(fiber.StatusBadRequest).JSON(shared.PipeRes[any]{Success: false, Message: m})
}
func user(c *fiber.Ctx) (uuid.UUID, bool) { return middleware.CurrentUserID(c) }
func (c *Controller) Invite(x *fiber.Ctx) error {
	u, ok := user(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	e, er := uuid.Parse(x.Params("eventID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	var d dtos.InviteDTO
	if x.BodyParser(&d) != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.Invite(x.Context(), u, e, d)
	if !r.Success {
		return fail(x, r.Message)
	}
	return x.Status(201).JSON(r)
}
func (c *Controller) Roles(x *fiber.Ctx) error {
	u, ok := user(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	m, er := uuid.Parse(x.Params("membershipID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	var d dtos.RolesDTO
	if x.BodyParser(&d) != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.Roles(x.Context(), u, m, d)
	if !r.Success {
		return fail(x, r.Message)
	}
	return x.Status(200).JSON(r)
}
func (c *Controller) Accept(x *fiber.Ctx) error {
	u, ok := user(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	m, er := uuid.Parse(x.Params("membershipID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.Accept(x.Context(), u, m)
	if !r.Success {
		return fail(x, r.Message)
	}
	return x.Status(200).JSON(r)
}
func (c *Controller) Revoke(x *fiber.Ctx) error {
	u, ok := user(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	m, er := uuid.Parse(x.Params("membershipID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.Revoke(x.Context(), u, m)
	if !r.Success {
		return fail(x, r.Message)
	}
	return x.Status(200).JSON(r)
}
func (c *Controller) Shift(x *fiber.Ctx) error {
	u, ok := user(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	m, er := uuid.Parse(x.Params("membershipID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	var d dtos.ShiftDTO
	if x.BodyParser(&d) != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.Shift(x.Context(), u, m, d)
	if !r.Success {
		return fail(x, r.Message)
	}
	return x.Status(201).JSON(r)
}
func (c *Controller) List(x *fiber.Ctx) error {
	u, ok := user(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	e, er := uuid.Parse(x.Params("eventID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.List(x.Context(), u, e)
	if !r.Success {
		return fail(x, r.Message)
	}
	return x.Status(200).JSON(r)
}
