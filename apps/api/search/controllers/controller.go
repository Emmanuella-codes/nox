package controllers

import (
	"strconv"

	"github.com/emmanuella-codes/nox/search/messages"
	"github.com/emmanuella-codes/nox/search/pipes"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/gofiber/fiber/v2"
)

type SearchController struct {
	pipe *pipes.SearchPipe
}

func NewSearchController(pipe *pipes.SearchPipe) *SearchController {
	return &SearchController{pipe: pipe}
}

func pipeSuccess[T any](ctx *fiber.Ctx, status int, message shared.PipeMessage, data *T) error {
	return ctx.Status(status).JSON(shared.PipeRes[T]{Success: true, Message: message, Data: data})
}

func pipeError(ctx *fiber.Ctx, status int, message shared.PipeMessage) error {
	return ctx.Status(status).JSON(shared.PipeRes[any]{Success: false, Message: message})
}

func pipeErrorStatus(message shared.PipeMessage) int {
	switch message {
	case messages.Invalid_Query:
		return fiber.StatusBadRequest
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
