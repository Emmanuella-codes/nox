package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *AdminPipe) ListReportsPipe(ctx context.Context, actorID uuid.UUID, params adminrepo.ListReportsParams) *shared.PipeRes[[]models.Report] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[[]models.Report](shared.CreatePipeMessage(message))
	}
	reports, err := p.adminRepo.ListReports(ctx, params)
	if err != nil {
		logInternalError(err, "reports.list")
		return shared.PipeError[[]models.Report](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.reports.list", nil)
	return shared.PipeSuccess(messages.Reports_Loaded, &reports)
}

func (p *AdminPipe) GetReportPipe(ctx context.Context, actorID, reportID uuid.UUID) *shared.PipeRes[models.Report] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.Report](shared.CreatePipeMessage(message))
	}
	report, err := p.adminRepo.FindReportByID(ctx, reportID)
	if err != nil {
		if err == adminrepo.ErrReportNotFound {
			return shared.PipeError[models.Report](messages.Report_Not_Found)
		}
		logInternalError(err, "reports.get")
		return shared.PipeError[models.Report](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.reports.get", map[string]any{"report_id": reportID.String()})
	return shared.PipeSuccess(messages.Report_Loaded, report)
}

func (p *AdminPipe) UpdateReportPipe(ctx context.Context, actorID, reportID uuid.UUID, dto dtos.UpdateReportDTO) *shared.PipeRes[models.Report] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.Report](shared.CreatePipeMessage(message))
	}
	if !validReportStatus(dto.Status) {
		return shared.PipeError[models.Report](messages.Invalid_Report_Status)
	}
	dto.ResolutionNote = strings.TrimSpace(dto.ResolutionNote)
	report, err := p.adminRepo.UpdateReport(ctx, adminrepo.UpdateReportParams{
		ReportID: reportID, Status: dto.Status, AssignedAdminID: dto.AssignedAdminID, ResolutionNote: dto.ResolutionNote,
	})
	if err != nil {
		if err == adminrepo.ErrReportNotFound {
			return shared.PipeError[models.Report](messages.Report_Not_Found)
		}
		logInternalError(err, "reports.update")
		return shared.PipeError[models.Report](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.reports.update", map[string]any{"report_id": reportID.String(), "status": dto.Status})
	return shared.PipeSuccess(messages.Report_Updated, report)
}

func (p *AdminPipe) ResolveReportPipe(ctx context.Context, actorID, reportID uuid.UUID, dto dtos.ReportActionDTO) *shared.PipeRes[models.Report] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.Report](shared.CreatePipeMessage(message))
	}
	dto.Action = strings.ToLower(strings.TrimSpace(dto.Action))
	dto.Reason = strings.TrimSpace(dto.Reason)
	dto.Resolution = strings.TrimSpace(dto.Resolution)
	if dto.Action == "dismiss" {
		if dto.Resolution == "" {
			return shared.PipeError[models.Report](messages.Invalid_Payload)
		}
		report, err := p.adminRepo.UpdateReport(ctx, adminrepo.UpdateReportParams{ReportID: reportID, Status: models.ReportStatusDismissed, ResolutionNote: dto.Resolution, AssignedAdminID: &actorID})
		if err != nil {
			return reportError[models.Report](err)
		}
		p.audit(ctx, actorID, "admin.reports.dismiss", map[string]any{"report_id": reportID.String()})
		return shared.PipeSuccess(messages.Report_Resolved, report)
	}
	status, ok := reportModerationStatus(dto.Action)
	if !ok || (status != models.ModerationStatusActive && dto.Reason == "") {
		return shared.PipeError[models.Report](messages.Invalid_Report_Action)
	}
	report, err := p.adminRepo.ResolveReport(ctx, adminrepo.ResolveReportParams{ReportID: reportID, Status: status, Reason: dto.Reason, AdminID: actorID, Resolution: dto.Resolution})
	if err != nil {
		return reportError[models.Report](err)
	}
	p.audit(ctx, actorID, "admin.reports.resolve", map[string]any{"report_id": reportID.String(), "action": dto.Action})
	return shared.PipeSuccess(messages.Report_Resolved, report)
}

func reportModerationStatus(action string) (models.ModerationStatus, bool) {
	switch action {
	case "hide":
		return models.ModerationStatusHidden, true
	case "remove":
		return models.ModerationStatusRemoved, true
	case "restore":
		return models.ModerationStatusActive, true
	default:
		return "", false
	}
}

func validReportStatus(status models.ReportStatus) bool {
	return status == models.ReportStatusOpen || status == models.ReportStatusReviewing || status == models.ReportStatusResolved || status == models.ReportStatusDismissed
}

func reportError[T any](err error) *shared.PipeRes[T] {
	if err == adminrepo.ErrReportNotFound || err == adminrepo.ErrModerationEntityNotFound {
		return shared.PipeError[T](messages.Report_Not_Found)
	}
	logInternalError(err, "reports.resolve")
	return shared.PipeError[T](messages.Internal_Error)
}
