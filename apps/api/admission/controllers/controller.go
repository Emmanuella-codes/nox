package controllers

import (
	"github.com/emmanuella-codes/nox/admission/dtos"
	"github.com/emmanuella-codes/nox/admission/messages"
	"github.com/emmanuella-codes/nox/admission/pipes"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ pipe *pipes.AdmissionPipe }

func NewController(p *pipes.AdmissionPipe) *Controller { return &Controller{pipe: p} }

func (c *Controller) Create(x *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	e, er := uuid.Parse(x.Params("eventID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	var d dtos.CreateAdmissionDTO
	if x.BodyParser(&d) != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.Create(x.Context(), u, e, d)
	if !r.Success {
		return fail(x, r.Message)
	}
	return x.Status(201).JSON(r)
}
func fail(c *fiber.Ctx, m shared.PipeMessage) error {
	return c.Status(fiber.StatusBadRequest).JSON(shared.PipeRes[any]{Success: false, Message: m})
}
func (c *Controller) QR(x *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	e, er := uuid.Parse(x.Params("eventID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	var d struct {
		Credential string `json:"credential"`
	}
	if x.BodyParser(&d) != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.QR(x.Context(), u, e, d.Credential)
	if !r.Success {
		return x.Status(fiber.StatusConflict).JSON(r)
	}
	return x.Status(200).JSON(r)
}
func (c *Controller) Manual(x *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	e, er := uuid.Parse(x.Params("eventID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	var d dtos.ManualCheckInDTO
	if x.BodyParser(&d) != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.Manual(x.Context(), u, e, d)
	if !r.Success {
		return x.Status(fiber.StatusConflict).JSON(r)
	}
	return x.Status(200).JSON(r)
}
func (c *Controller) Override(x *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	a, er := uuid.Parse(x.Params("admissionID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	var d dtos.OverrideDTO
	if x.BodyParser(&d) != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.Override(x.Context(), u, a, d)
	if !r.Success {
		return x.Status(fiber.StatusConflict).JSON(r)
	}
	return x.Status(200).JSON(r)
}
func (c *Controller) Search(x *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	e, er := uuid.Parse(x.Params("eventID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.Search(x.Context(), u, e, x.Query("q"))
	if !r.Success {
		return fail(x, r.Message)
	}
	return x.Status(200).JSON(r)
}
func (c *Controller) Reconcile(x *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(x)
	if !ok {
		return fail(x, messages.Forbidden)
	}
	e, er := uuid.Parse(x.Params("eventID"))
	if er != nil {
		return fail(x, messages.Invalid)
	}
	r := c.pipe.Reconcile(x.Context(), u, e)
	if !r.Success {
		return fail(x, r.Message)
	}
	return x.Status(200).JSON(r)
}
