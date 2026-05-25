package controllers

import "github.com/gofiber/fiber/v2"

func (c *LikeController) UnlikePost(ctx *fiber.Ctx) error {
	return c.likeAction(ctx, false)
}
