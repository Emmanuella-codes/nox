package controllers

import (
	"github.com/emmanuella-codes/nox/event/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *EventController) CreateEvent(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	var dto dtos.CreateEventDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.CreateEventPipe(ctx.Context(), userID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

func (c *EventController) UpdateEvent(ctx *fiber.Ctx) error  { return c.mutateEvent(ctx, "update") }
func (c *EventController) PublishEvent(ctx *fiber.Ctx) error { return c.mutateEvent(ctx, "publish") }
func (c *EventController) CloseEvent(ctx *fiber.Ctx) error   { return c.mutateEvent(ctx, "close") }
func (c *EventController) ArchiveEvent(ctx *fiber.Ctx) error { return c.mutateEvent(ctx, "archive") }

func (c *EventController) mutateEvent(ctx *fiber.Ctx, action string) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	eventID, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_event_id")
	}
	var res *shared.PipeRes[models.Event]
	if action == "update" {
		var dto dtos.UpdateEventDTO
		if err := parseAndValidate(ctx, &dto); err != nil {
			return validationError(ctx, err)
		}
		res = c.pipe.UpdateEventPipe(ctx.Context(), userID, eventID, dto)
	} else if action == "publish" {
		res = c.pipe.PublishEventPipe(ctx.Context(), userID, eventID)
	} else if action == "close" {
		res = c.pipe.CloseEventPipe(ctx.Context(), userID, eventID)
	} else {
		res = c.pipe.ArchiveEventPipe(ctx.Context(), userID, eventID)
	}
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
