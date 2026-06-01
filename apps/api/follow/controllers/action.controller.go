package controllers

import "github.com/gofiber/fiber/v2"

func (c *FollowController) FollowPersona(ctx *fiber.Ctx) error {
	return c.followAction(ctx, true)
}

func (c *FollowController) UnfollowPersona(ctx *fiber.Ctx) error {
	return c.followAction(ctx, false)
}
