package controllers

import "github.com/gofiber/fiber/v2"

func (c *PreferenceController) BlockUser(ctx *fiber.Ctx) error {
	return c.blockAction(ctx, true)
}

func (c *PreferenceController) UnblockUser(ctx *fiber.Ctx) error {
	return c.blockAction(ctx, false)
}

func (c *PreferenceController) MuteUser(ctx *fiber.Ctx) error {
	return c.muteAction(ctx, true)
}

func (c *PreferenceController) UnmuteUser(ctx *fiber.Ctx) error {
	return c.muteAction(ctx, false)
}

func (c *PreferenceController) AddDiscoverySuppression(ctx *fiber.Ctx) error {
	return c.discoverySuppressionAction(ctx, true)
}

func (c *PreferenceController) RemoveDiscoverySuppression(ctx *fiber.Ctx) error {
	return c.discoverySuppressionAction(ctx, false)
}
