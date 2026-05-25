package controllers

import (
	"github.com/emmanuella-codes/nox/like/dtos"
	"github.com/emmanuella-codes/nox/like/messages"
	"github.com/emmanuella-codes/nox/like/pipes"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared"
	sharedapi "github.com/emmanuella-codes/nox/shared/api"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LikeController struct {
	pipe *pipes.LikePipe
}

func NewLikeController(pipe *pipes.LikePipe) *LikeController {
	return &LikeController{pipe: pipe}
}

func (c *LikeController) likeAction(ctx *fiber.Ctx, like bool) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	postID, err := uuid.Parse(ctx.Params("postID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_post_id")
	}

	var dto dtos.LikePostDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	var res *shared.PipeRes[any]
	if like {
		res = c.pipe.LikePostPipe(ctx.Context(), userID, postID, dto)
	} else {
		res = c.pipe.UnlikePostPipe(ctx.Context(), userID, postID, dto)
	}
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
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
	case messages.Invalid_Payload:
		return fiber.StatusBadRequest
	case messages.Post_Not_Found, messages.Persona_Not_Found:
		return fiber.StatusNotFound
	case messages.Forbidden:
		return fiber.StatusForbidden
	case messages.Internal_Error:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusBadRequest
	}
}
