package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/preference/dtos"
	"github.com/emmanuella-codes/nox/preference/messages"
	"github.com/emmanuella-codes/nox/preference/pipes"
	"github.com/emmanuella-codes/nox/shared"
	sharedapi "github.com/emmanuella-codes/nox/shared/api"
	"github.com/gofiber/fiber/v2"
)

type PreferenceController struct {
	pipe *pipes.PreferencePipe
}

func NewPreferenceController(pipe *pipes.PreferencePipe) *PreferenceController {
	return &PreferenceController{pipe: pipe}
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

func pipeError(ctx *fiber.Ctx, status int, message shared.PipeMessage) error {
	return ctx.Status(status).JSON(shared.PipeRes[any]{Success: false, Message: message})
}

func pipeSuccess(ctx *fiber.Ctx, status int, message shared.PipeMessage) error {
	return ctx.Status(status).JSON(shared.PipeRes[any]{Success: true, Message: message})
}

func pipeErrorStatus(message shared.PipeMessage) int {
	switch message {
	case messages.Persona_Not_Found, messages.Target_Not_Found:
		return fiber.StatusNotFound
	case messages.Forbidden:
		return fiber.StatusForbidden
	case messages.Internal_Error:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusBadRequest
	}
}

func (c *PreferenceController) blockAction(ctx *fiber.Ctx, block bool) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	var dto dtos.PersonaTargetDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	var res *shared.PipeRes[any]
	if block {
		res = c.pipe.BlockUserPipe(ctx.Context(), userID, dto)
	} else {
		res = c.pipe.UnblockUserPipe(ctx.Context(), userID, dto)
	}
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message)
}

func (c *PreferenceController) muteAction(ctx *fiber.Ctx, mute bool) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	var dto dtos.PersonaTargetDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	var res *shared.PipeRes[any]
	if mute {
		res = c.pipe.MuteUserPipe(ctx.Context(), userID, dto)
	} else {
		res = c.pipe.UnmuteUserPipe(ctx.Context(), userID, dto)
	}
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message)
}

func (c *PreferenceController) discoverySuppressionAction(ctx *fiber.Ctx, add bool) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	var dto dtos.DiscoverySuppressionDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	var res *shared.PipeRes[any]
	if add {
		res = c.pipe.AddDiscoverySuppressionPipe(ctx.Context(), userID, dto)
	} else {
		res = c.pipe.RemoveDiscoverySuppressionPipe(ctx.Context(), userID, dto)
	}
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message)
}
