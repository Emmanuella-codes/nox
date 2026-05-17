package middleware

import (
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const personaLocalKey = "persona"

func RequirePersonaOwner(repo persona.PersonaRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := CurrentUserID(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(shared.PipeRes[any]{
				Success: false,
				Message: "invalid_token",
			})
		}

		personaID, err := resolvePersonaID(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(shared.PipeRes[any]{
				Success: false,
				Message: "invalid_persona_id",
			})
		}

		foundPersona, err := repo.FindPersonaByID(c.Context(), personaID)
		if err == persona.ErrPersonaNotFound {
			return c.Status(fiber.StatusNotFound).JSON(shared.PipeRes[any]{
				Success: false,
				Message: "persona_not_found",
			})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(shared.PipeRes[any]{
				Success: false,
				Message: "internal_error",
			})
		}
		if foundPersona.UserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(shared.PipeRes[any]{
				Success: false,
				Message: "forbidden",
			})
		}

		c.Locals(personaLocalKey, foundPersona)
		return c.Next()
	}
}

func CurrentPersona(c *fiber.Ctx) (*models.Persona, bool) {
	foundPersona, ok := c.Locals(personaLocalKey).(*models.Persona)
	return foundPersona, ok
}

func resolvePersonaID(c *fiber.Ctx) (uuid.UUID, error) {
	if rawID := c.Params("personaID"); rawID != "" {
		return uuid.Parse(rawID)
	}

	var payload struct {
		PersonaID uuid.UUID `json:"persona_id"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return uuid.Nil, err
	}

	return payload.PersonaID, nil
}
