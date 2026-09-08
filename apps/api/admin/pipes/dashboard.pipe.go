package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/shared"
	worker_runtime "github.com/emmanuella-codes/nox/workers/runtime"
	"github.com/google/uuid"
)

func (p *AdminPipe) DashboardPipe(ctx context.Context, actorID uuid.UUID) *shared.PipeRes[models.AdminDashboard] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminDashboard](shared.CreatePipeMessage(message))
	}
	dashboard, err := p.adminRepo.Dashboard(ctx)
	if err != nil {
		logInternalError(err, "dashboard.load")
		return shared.PipeError[models.AdminDashboard](messages.Internal_Error)
	}
	dashboard.Workers = dashboardWorkers(p.HealthPipe(ctx).Data)
	p.audit(ctx, actorID, "admin.dashboard.get", nil)
	return shared.PipeSuccess(messages.Admin_Dashboard_Loaded, dashboard)
}

func dashboardWorkers(health *HealthResponse) map[string]models.AdminDashboardWorker {
	workers := map[string]models.AdminDashboardWorker{
		worker_runtime.WorkerMediaKey:         {Status: "offline"},
		worker_runtime.WorkerNotificationsKey: {Status: "offline"},
		worker_runtime.WorkerStoriesKey:       {Status: "offline"},
		worker_runtime.WorkerCrewsKey:         {Status: "offline"},
		worker_runtime.WorkerModerationKey:    {Status: "offline"},
	}
	if health == nil {
		return workers
	}
	for name, worker := range health.Workers {
		workers[name] = models.AdminDashboardWorker{Status: worker.Status, LastSeen: worker.LastSeen}
	}
	return workers
}
