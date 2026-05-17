package typings

import "github.com/gofiber/fiber/v2"

type FiberMiddleware func(ctx *fiber.Ctx) error

type TokenGen struct {
	Ctx       *fiber.Ctx
	Parameter string
}

type TokenGeneratorFunc func(tkg *TokenGen) string
