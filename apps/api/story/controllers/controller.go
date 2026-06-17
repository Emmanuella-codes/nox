package controllers

import (
	"strconv"

	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared"
	sharedapi "github.com/emmanuella-codes/nox/shared/api"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/emmanuella-codes/nox/story/pipes"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StoryController struct {
	pipe *pipes.StoryPipe
}

func NewStoryController(pipe *pipes.StoryPipe) *StoryController {
	return &StoryController{pipe: pipe}
}

func parseAndValidate(ctx *fiber.Ctx, dto any) error {
	if err := ctx.BodyParser(dto); err != nil {
		return err
	}
	success, err := sharedapi.ValidateAPIData(dto)
	if !success {
		return err
	}
	return nil
}

func validationError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": messages.Invalid_Payload, "error": err.Error()})
}

func pipeSuccess[T any](ctx *fiber.Ctx, status int, message shared.PipeMessage, data *T) error {
	return ctx.Status(status).JSON(shared.PipeRes[T]{Success: true, Message: message, Data: data})
}

func pipeError(ctx *fiber.Ctx, status int, message shared.PipeMessage) error {
	return ctx.Status(status).JSON(shared.PipeRes[any]{Success: false, Message: message})
}

func pipeErrorStatus(message shared.PipeMessage) int {
	switch message {
	case messages.Invalid_Payload, messages.Invalid_Story, messages.Story_Duration_Limit_Exceeded, messages.Media_Asset_In_Use:
		return fiber.StatusBadRequest
	case messages.Story_Not_Found, messages.Story_Item_Not_Found, messages.Event_Not_Found, messages.Persona_Not_Found, messages.Media_Asset_Not_Found:
		return fiber.StatusNotFound
	case messages.Forbidden:
		return fiber.StatusForbidden
	case messages.Internal_Error:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusBadRequest
	}
}

func queryLimit(ctx *fiber.Ctx, fallback int) int {
	value := ctx.Query("limit")
	if value == "" {
		return fallback
	}
	limit, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return limit
}

func queryOffset(ctx *fiber.Ctx) int {
	value := ctx.Query("offset")
	if value == "" {
		return 0
	}
	offset, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return offset
}

func optionalViewerPersonaID(ctx *fiber.Ctx) (*uuid.UUID, error) {
	value := ctx.Query("viewer_persona_id")
	if value == "" {
		return nil, nil
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func optionalViewerContext(ctx *fiber.Ctx) (*uuid.UUID, *uuid.UUID, error) {
	viewerPersonaID, err := optionalViewerPersonaID(ctx)
	if err != nil {
		return nil, nil, err
	}
	if viewerPersonaID == nil {
		return nil, nil, nil
	}
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return nil, nil, fiber.ErrUnauthorized
	}
	return &userID, viewerPersonaID, nil
}

func requiredPersonaQuery(ctx *fiber.Ctx, key string) (uuid.UUID, error) {
	return uuid.Parse(ctx.Query(key))
}
