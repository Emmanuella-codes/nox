package controllers

import (
	"github.com/emmanuella-codes/nox/follow/pipes"
	"github.com/emmanuella-codes/nox/middleware"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *FollowController) GetFollowers(ctx *fiber.Ctx) error {
	personaID, err := uuid.Parse(ctx.Params("personaID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	options := follow_repo.ListOptions{
		Limit:  queryLimit(ctx, 20),
		Offset: queryOffset(ctx),
	}

	var res *shared.PipeRes[pipes.FollowListResponse]
	viewerPersonaID := ctx.Query("viewer_persona_id")
	if viewerPersonaID != "" {
		userID, ok := middleware.CurrentUserID(ctx)
		if !ok {
			return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
		}
		parsedViewerPersonaID, err := uuid.Parse(viewerPersonaID)
		if err != nil {
			return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
		}
		res = c.pipe.FollowersForViewerPipe(ctx.Context(), userID, personaID, parsedViewerPersonaID, options)
	} else {
		res = c.pipe.FollowersPipe(ctx.Context(), personaID, options)
	}
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *FollowController) GetFollowing(ctx *fiber.Ctx) error {
	personaID, err := uuid.Parse(ctx.Params("personaID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	options := follow_repo.ListOptions{
		Limit:  queryLimit(ctx, 20),
		Offset: queryOffset(ctx),
	}

	var res *shared.PipeRes[pipes.FollowListResponse]
	viewerPersonaID := ctx.Query("viewer_persona_id")
	if viewerPersonaID != "" {
		userID, ok := middleware.CurrentUserID(ctx)
		if !ok {
			return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
		}
		parsedViewerPersonaID, err := uuid.Parse(viewerPersonaID)
		if err != nil {
			return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
		}
		res = c.pipe.FollowingForViewerPipe(ctx.Context(), userID, personaID, parsedViewerPersonaID, options)
	} else {
		res = c.pipe.FollowingPipe(ctx.Context(), personaID, options)
	}
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
