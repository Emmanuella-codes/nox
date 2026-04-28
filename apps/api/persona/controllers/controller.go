package controllers

import (
	"strconv"

	"github.com/emmanuella-codes/nox/persona/messages"
	"github.com/emmanuella-codes/nox/persona/pipes"
	"github.com/emmanuella-codes/nox/shared"
	sharedapi "github.com/emmanuella-codes/nox/shared/api"
	"github.com/gofiber/fiber/v2"
)

type PersonaController struct {
	pipe *pipes.PersonaPipe
}

func NewPersonaController(pipe *pipes.PersonaPipe) *PersonaController {
	return &PersonaController{pipe: pipe}
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
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": messages.Invalid_Payload,
		"error":   err.Error(),
	})
}

func pipeSuccess[T any](ctx *fiber.Ctx, status int, message shared.PipeMessage, data *T) error {
	return ctx.Status(status).JSON(shared.PipeRes[T]{Success: true, Message: message, Data: data})
}

func pipeError(ctx *fiber.Ctx, status int, message shared.PipeMessage) error {
	return ctx.Status(status).JSON(shared.PipeRes[any]{Success: false, Message: message})
}

func pipeErrorStatus(message shared.PipeMessage) int {
	switch message {
	case messages.Handle_Already_Taken, messages.Ghost_Persona_Already_Set:
		return fiber.StatusConflict
	case messages.Invalid_Payload, messages.Invalid_Persona_Type:
		return fiber.StatusBadRequest
	case messages.Persona_Not_Found:
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
