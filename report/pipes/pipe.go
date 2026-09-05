package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/models"
	reportrepo "github.com/emmanuella-codes/nox/repositories/admin"
	personarepo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/report/dtos"
	"github.com/emmanuella-codes/nox/report/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type ReportPipe struct {
	reportRepo  reportrepo.AdminRepository
	personaRepo personarepo.PersonaRepository
}

func NewReportPipe(reportRepo reportrepo.AdminRepository, personaRepo personarepo.PersonaRepository) *ReportPipe {
	return &ReportPipe{reportRepo: reportRepo, personaRepo: personaRepo}
}

func (p *ReportPipe) CreateReportPipe(ctx context.Context, userID uuid.UUID, dto dtos.CreateReportDTO) *shared.PipeRes[models.Report] {
	dto.Reason = strings.TrimSpace(dto.Reason)
	dto.Description = strings.TrimSpace(dto.Description)
	if !validTarget(dto.TargetType) {
		return shared.PipeError[models.Report](messages.InvalidReportTarget)
	}
	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.ReporterPersonaID)
	if err != nil {
		if err == personarepo.ErrPersonaNotFound {
			return shared.PipeError[models.Report](messages.PersonaNotFound)
		}
		return shared.PipeError[models.Report](messages.InternalError)
	}
	if !persona.IsOwnedBy(userID) {
		return shared.PipeError[models.Report](messages.Forbidden)
	}
	report, err := p.reportRepo.CreateReport(ctx, reportrepo.CreateReportParams{
		ReporterUserID: userID, ReporterPersonaID: dto.ReporterPersonaID, TargetType: dto.TargetType,
		TargetID: dto.TargetID, Reason: dto.Reason, Description: dto.Description,
	})
	if err != nil {
		switch err {
		case reportrepo.ErrReportAlreadyExists:
			return shared.PipeError[models.Report](messages.ReportAlreadyExists)
		case reportrepo.ErrSelfReport:
			return shared.PipeError[models.Report](messages.CannotReportOwnData)
		case reportrepo.ErrModerationEntityNotFound:
			return shared.PipeError[models.Report](messages.ReportTargetNotFound)
		default:
			return shared.PipeError[models.Report](messages.InternalError)
		}
	}
	return shared.PipeSuccess(messages.ReportCreated, report)
}

func validTarget(target models.ReportTargetType) bool {
	switch target {
	case models.ReportTargetPersona, models.ReportTargetPost, models.ReportTargetComment,
		models.ReportTargetStory, models.ReportTargetStoryItem, models.ReportTargetSet, models.ReportTargetEvent:
		return true
	default:
		return false
	}
}
