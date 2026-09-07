package controllers

import (
	"context"

	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/admin/pipes"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/gofiber/fiber/v2"
)

type AdminController struct {
	pipe *pipes.AdminPipe
}

func NewAdminController(pipe *pipes.AdminPipe) *AdminController {
	return &AdminController{pipe: pipe}
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
	case messages.Invalid_Credentials, messages.Invalid_Token:
		return fiber.StatusUnauthorized
	case messages.Admin_Access_Denied:
		return fiber.StatusForbidden
	case messages.Admin_User_Not_Found:
		return fiber.StatusNotFound
	case messages.Admin_User_Exists:
		return fiber.StatusConflict
	case messages.User_Not_Found:
		return fiber.StatusNotFound
	case messages.User_Already_Verified:
		return fiber.StatusConflict
	case messages.Moderation_Entity_Not_Found:
		return fiber.StatusNotFound
	case messages.Hashtag_Not_Found, messages.Search_Suppression_Not_Found:
		return fiber.StatusNotFound
	case messages.Media_Not_Found:
		return fiber.StatusNotFound
	case messages.Report_Not_Found:
		return fiber.StatusNotFound
	case messages.Crew_Not_Found:
		return fiber.StatusNotFound
	case messages.Invalid_Moderation_Status:
		return fiber.StatusBadRequest
	case messages.Invalid_Report_Status, messages.Invalid_Report_Action:
		return fiber.StatusBadRequest
	case messages.Admin_Dashboard_Loaded:
		return fiber.StatusOK
	case messages.Invalid_Payload:
		return fiber.StatusBadRequest
	case messages.Invalid_Admin_Role, messages.Invalid_User_Status, messages.Last_Super_Admin_Required:
		return fiber.StatusBadRequest
	case messages.Internal_Error:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusBadRequest
	}
}

func requestContext(ctx *fiber.Ctx) context.Context {
	var base context.Context = ctx.Context()
	base = context.WithValue(base, pipes.RequestIDContextKey(), localString(ctx, middleware.RequestIDLocalKey))
	base = context.WithValue(base, pipes.IPAddressContextKey(), ctx.IP())
	base = context.WithValue(base, pipes.UserAgentContextKey(), string(ctx.Request().Header.UserAgent()))
	return base
}

func localString(ctx *fiber.Ctx, key string) string {
	value, _ := ctx.Locals(key).(string)
	return value
}
