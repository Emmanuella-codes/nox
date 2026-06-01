package controllers

import (
	"strconv"

	"github.com/emmanuella-codes/nox/hashtag/messages"
	"github.com/emmanuella-codes/nox/hashtag/pipes"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/gofiber/fiber/v2"
)

type HashtagController struct {
	pipe *pipes.HashtagPipe
}

func NewHashtagController(pipe *pipes.HashtagPipe) *HashtagController {
	return &HashtagController{pipe: pipe}
}

func pipeSuccess[T any](ctx *fiber.Ctx, status int, message shared.PipeMessage, data *T) error {
	return ctx.Status(status).JSON(shared.PipeRes[T]{Success: true, Message: message, Data: data})
}

func pipeError(ctx *fiber.Ctx, status int, message shared.PipeMessage) error {
	return ctx.Status(status).JSON(shared.PipeRes[any]{Success: false, Message: message})
}

func pipeErrorStatus(message shared.PipeMessage) int {
	switch message {
	case messages.Invalid_Tag:
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
