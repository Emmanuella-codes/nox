package controllers

import (
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/media/messages"
	"github.com/emmanuella-codes/nox/media/pipes"
	"github.com/emmanuella-codes/nox/shared"
	sharedapi "github.com/emmanuella-codes/nox/shared/api"
	"github.com/gofiber/fiber/v2"
)

type MediaController struct {
	pipe   *pipes.MediaPipe
	config *config.Config
}

func NewMediaController(pipe *pipes.MediaPipe, cfg *config.Config) *MediaController {
	return &MediaController{pipe: pipe, config: cfg}
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
	case messages.Invalid_Payload, messages.Invalid_Media:
		return fiber.StatusBadRequest
	case messages.Persona_Not_Found, messages.Media_Not_Found:
		return fiber.StatusNotFound
	case messages.Forbidden:
		return fiber.StatusForbidden
	case messages.Internal_Error:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusBadRequest
	}
}

func (c *MediaController) validProcessingSecret(ctx *fiber.Ctx) bool {
	if c.config == nil || c.config.MediaProcessingSecret == "" {
		return false
	}
	return ctx.Get("X-Media-Processing-Secret") == c.config.MediaProcessingSecret
}
