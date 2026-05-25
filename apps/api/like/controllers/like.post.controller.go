package controllers

import "github.com/gofiber/fiber/v2"

func (c *LikeController) LikePost(ctx *fiber.Ctx) error {
	return c.likeAction(ctx, true)
}
