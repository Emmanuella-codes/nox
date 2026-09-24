package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/ticket/dtos"
	"github.com/emmanuella-codes/nox/ticket/messages"
	"github.com/emmanuella-codes/nox/ticket/pipes"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ pipe *pipes.TicketPipe }

func NewController(pipe *pipes.TicketPipe) *Controller { return &Controller{pipe: pipe} }
func success[T any](c *fiber.Ctx, status int, msg shared.PipeMessage, data *T) error {
	return c.Status(status).JSON(shared.PipeRes[T]{Success: true, Message: msg, Data: data})
}
func failure(c *fiber.Ctx, msg shared.PipeMessage) error {
	return c.Status(fiber.StatusBadRequest).JSON(shared.PipeRes[any]{Success: false, Message: msg})
}
func (c *Controller) Products(ctx *fiber.Ctx) error {
	id, e := uuid.Parse(ctx.Params("eventID"))
	if e != nil {
		return failure(ctx, messages.Invalid)
	}
	r := c.pipe.Products(ctx.Context(), id)
	if !r.Success {
		return failure(ctx, r.Message)
	}
	return success(ctx, 200, r.Message, r.Data)
}
func (c *Controller) CreateWave(ctx *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return failure(ctx, messages.Forbidden)
	}
	e, er := uuid.Parse(ctx.Params("eventID"))
	if er != nil {
		return failure(ctx, messages.Invalid)
	}
	var d dtos.CreateWaveDTO
	if ctx.BodyParser(&d) != nil {
		return failure(ctx, messages.Invalid)
	}
	r := c.pipe.CreateWave(ctx.Context(), u, e, d)
	if !r.Success {
		return failure(ctx, r.Message)
	}
	return success(ctx, 201, r.Message, r.Data)
}
func (c *Controller) CreateTier(ctx *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return failure(ctx, messages.Forbidden)
	}
	w, er := uuid.Parse(ctx.Params("waveID"))
	if er != nil {
		return failure(ctx, messages.Invalid)
	}
	var d dtos.CreateTierDTO
	if ctx.BodyParser(&d) != nil {
		return failure(ctx, messages.Invalid)
	}
	r := c.pipe.CreateTier(ctx.Context(), u, w, d)
	if !r.Success {
		return failure(ctx, r.Message)
	}
	return success(ctx, 201, r.Message, r.Data)
}
func (c *Controller) ReleaseWave(ctx *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return failure(ctx, messages.Forbidden)
	}
	w, err := uuid.Parse(ctx.Params("waveID"))
	if err != nil {
		return failure(ctx, messages.Invalid)
	}
	r := c.pipe.ReleaseWave(ctx.Context(), u, w)
	if !r.Success {
		return failure(ctx, r.Message)
	}
	return success(ctx, fiber.StatusOK, r.Message, r.Data)
}
func (c *Controller) Reserve(ctx *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return failure(ctx, messages.Forbidden)
	}
	t, er := uuid.Parse(ctx.Params("tierID"))
	if er != nil {
		return failure(ctx, messages.Invalid)
	}
	var d dtos.ReserveDTO
	if ctx.BodyParser(&d) != nil {
		return failure(ctx, messages.Invalid)
	}
	r := c.pipe.Reserve(ctx.Context(), u, t, d)
	if !r.Success {
		return failure(ctx, r.Message)
	}
	return success(ctx, 201, r.Message, r.Data)
}

func (c *Controller) Order(ctx *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return failure(ctx, messages.Forbidden)
	}
	o, err := uuid.Parse(ctx.Params("orderID"))
	if err != nil {
		return failure(ctx, messages.Invalid)
	}
	r := c.pipe.Order(ctx.Context(), u, o)
	if !r.Success {
		return failure(ctx, r.Message)
	}
	return success(ctx, fiber.StatusOK, r.Message, r.Data)
}

func (c *Controller) InitializePayment(ctx *fiber.Ctx) error {
	u, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return failure(ctx, messages.Forbidden)
	}
	o, err := uuid.Parse(ctx.Params("orderID"))
	if err != nil {
		return failure(ctx, messages.Invalid)
	}
	r := c.pipe.InitializePayment(ctx.Context(), u, o)
	if !r.Success {
		return failure(ctx, r.Message)
	}
	return success(ctx, fiber.StatusOK, r.Message, r.Data)
}

func (c *Controller) PaystackWebhook(ctx *fiber.Ctx) error {
	signature := ctx.Get("x-paystack-signature")
	if !c.pipe.VerifyWebhook(ctx.Body(), signature) {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}
	if err := c.pipe.HandleWebhook(ctx.Context(), ctx.Body()); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	return ctx.SendStatus(fiber.StatusOK)
}
