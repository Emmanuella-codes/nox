package controllers

import (
	"github.com/emmanuella-codes/nox/report/messages"
	"github.com/emmanuella-codes/nox/report/pipes"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/shared/api"
	"github.com/gofiber/fiber/v2"
)

type ReportController struct {
	pipe *pipes.ReportPipe
}

func NewReportController(pipe *pipes.ReportPipe) *ReportController {
	return &ReportController{pipe: pipe}
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

func (c *ReportController) CreateReport(ctx *fiber.Ctx) error {
	userID, ok := currentUserID(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(shared.PipeRes[any]{Success: false, Message: "invalid_token"})
	}
	var dto dtos.CreateReportDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(shared.PipeRes[any]{Success: false, Message: messages.InvalidPayload})
	}
	res := c.pipe.CreateReportPipe(ctx.Context(), userID, dto)
	if res.Success {
		return ctx.Status(fiber.StatusCreated).JSON(res)
	}
	return ctx.Status(reportErrorStatus(res.Message)).JSON(res)
}

func currentUserID(ctx *fiber.Ctx) (uuid.UUID, bool) {
	return middleware.CurrentUserID(ctx)
}

func reportErrorStatus(message shared.PipeMessage) int {
	switch message {
	case messages.ReportTargetNotFound, messages.PersonaNotFound:
		return fiber.StatusNotFound
	case messages.CannotReportOwnData, messages.Forbidden:
		return fiber.StatusForbidden
	case messages.ReportAlreadyExists:
		return fiber.StatusConflict
	case messages.InternalError:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusBadRequest
	}
}
