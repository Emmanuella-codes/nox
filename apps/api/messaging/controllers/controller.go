package controllers

import (
	"strconv"

	"github.com/emmanuella-codes/nox/messaging/messages"
	"github.com/emmanuella-codes/nox/messaging/pipes"
	"github.com/emmanuella-codes/nox/shared"
	sharedapi "github.com/emmanuella-codes/nox/shared/api"
	"github.com/gofiber/fiber/v2"
)

type MessagingController struct {
	pipe *pipes.MessagingPipe
}

// NewMessagingController builds the messaging HTTP controller.
func NewMessagingController(pipe *pipes.MessagingPipe) *MessagingController {
	return &MessagingController{pipe: pipe}
}

// parseAndValidate binds and validates one request payload.
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

// validationError returns one consistent validation failure response.
func validationError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": messages.Invalid_Payload, "error": err.Error()})
}

// pipeSuccess writes one successful pipe response to the client.
func pipeSuccess[T any](ctx *fiber.Ctx, status int, message shared.PipeMessage, data *T) error {
	return ctx.Status(status).JSON(shared.PipeRes[T]{Success: true, Message: message, Data: data})
}

// pipeError writes one failed pipe response to the client.
func pipeError(ctx *fiber.Ctx, status int, message shared.PipeMessage) error {
	return ctx.Status(status).JSON(shared.PipeRes[any]{Success: false, Message: message})
}

// pipeErrorStatus maps pipe messages into HTTP status codes.
func pipeErrorStatus(message shared.PipeMessage) int {
	switch message {
	case messages.Invalid_Payload:
		return fiber.StatusBadRequest
	case messages.Conversation_Not_Found, messages.Message_Not_Found, messages.Persona_Not_Found:
		return fiber.StatusNotFound
	case messages.Forbidden:
		return fiber.StatusForbidden
	case messages.Internal_Error:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusBadRequest
	}
}

// queryLimit parses one optional list limit query parameter.
func queryLimit(ctx *fiber.Ctx, fallback int) int {
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		return fallback
	}
	return limit
}

// queryOffset parses one optional list offset query parameter.
func queryOffset(ctx *fiber.Ctx) int {
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		return 0
	}
	return offset
}
