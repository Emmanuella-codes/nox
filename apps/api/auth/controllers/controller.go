package controllers

import (
	"github.com/emmanuella-codes/nox/auth/messages"
	"github.com/emmanuella-codes/nox/auth/pipes"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	pipe *pipes.AuthPipe
}

func NewAuthController(pipe *pipes.AuthPipe) *AuthController {
	return &AuthController{pipe: pipe}
}

func parseAndValidate(ctx *fiber.Ctx, dto any) error {
	if err := ctx.BodyParser(dto); err != nil {
		return err
	}

	success, err := api.ValidateAPIData(dto)
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
	return ctx.Status(status).JSON(shared.PipeRes[T]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func pipeError(ctx *fiber.Ctx, status int, message shared.PipeMessage) error {
	return ctx.Status(status).JSON(shared.PipeRes[any]{
		Success: false,
		Message: message,
	})
}

func pipeErrorStatus(message shared.PipeMessage) int {
	switch message {
	case messages.User_Already_Exists:
		return fiber.StatusConflict
	case messages.Invalid_Credentials, messages.Invalid_Token:
		return fiber.StatusUnauthorized
	case messages.Invalid_Payload:
		return fiber.StatusBadRequest
	case messages.Email_Not_Verified:
		return fiber.StatusForbidden
	case messages.Email_Already_Verified:
		return fiber.StatusConflict
	case messages.Invalid_OTP, messages.OTP_Expired:
		return fiber.StatusUnauthorized
	case messages.OTP_Locked:
		return fiber.StatusTooManyRequests
	case messages.Internal_Error:
		return fiber.StatusInternalServerError
	case messages.Unknown_Error:
		return fiber.StatusBadRequest
	default:
		return fiber.StatusBadRequest
	}
}
