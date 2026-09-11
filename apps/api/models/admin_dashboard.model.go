package models

import "time"

type AdminDashboard struct {
	Users       AdminDashboardUsers             `json:"users"`
	Content     AdminDashboardContent           `json:"content"`
	Moderation  AdminDashboardModeration        `json:"moderation"`
	Operations  AdminDashboardOperations        `json:"operations"`
	Workers     map[string]AdminDashboardWorker `json:"workers"`
	GeneratedAt time.Time                       `json:"generated_at"`
}

type AdminDashboardWorker struct {
	Status   string     `json:"status"`
	LastSeen *time.Time `json:"last_seen,omitempty"`
}

type AdminDashboardUsers struct {
	Total       int64 `json:"total"`
	NewToday    int64 `json:"new_today"`
	NewThisWeek int64 `json:"new_this_week"`
	Suspended   int64 `json:"suspended"`
}

type AdminDashboardContent struct {
	Posts         int64 `json:"posts"`
	Comments      int64 `json:"comments"`
	Personas      int64 `json:"personas"`
	Events        int64 `json:"events"`
	Sets          int64 `json:"sets"`
	ActiveStories int64 `json:"active_stories"`
}

type AdminDashboardModeration struct {
	OpenReports      int64 `json:"open_reports"`
	ReviewingReports int64 `json:"reviewing_reports"`
	RecentActions    int64 `json:"recent_actions"`
}

type AdminDashboardOperations struct {
	FailedNotifications int64 `json:"failed_notifications"`
	PendingMedia        int64 `json:"pending_media"`
}
