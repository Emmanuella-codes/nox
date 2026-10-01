package repositories

import (
	"github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/repositories/admission"
	"github.com/emmanuella-codes/nox/repositories/comment"
	"github.com/emmanuella-codes/nox/repositories/crew"
	"github.com/emmanuella-codes/nox/repositories/event"
	"github.com/emmanuella-codes/nox/repositories/follow"
	"github.com/emmanuella-codes/nox/repositories/guest"
	"github.com/emmanuella-codes/nox/repositories/hashtag"
	"github.com/emmanuella-codes/nox/repositories/like"
	"github.com/emmanuella-codes/nox/repositories/media"
	"github.com/emmanuella-codes/nox/repositories/messaging"
	"github.com/emmanuella-codes/nox/repositories/notification"
	"github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/repositories/preference"
	"github.com/emmanuella-codes/nox/repositories/search"
	"github.com/emmanuella-codes/nox/repositories/set"
	"github.com/emmanuella-codes/nox/repositories/staff"
	"github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/repositories/ticket"
	"github.com/emmanuella-codes/nox/repositories/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRepositories(db *pgxpool.Pool) {
	admin.InitAdminRepo(db)
	user.InitUserRepo(db)
	persona.InitPersonaRepo(db)
	post.InitPostRepo(db)
	preference.InitPreferenceRepo(db)
	comment.InitCommentRepo(db)
	crew.InitCrewRepo(db)
	like.InitLikeRepo(db)
	event.InitEventRepo(db)
	media.InitMediaRepo(db)
	messaging.InitMessagingRepo(db)
	notification.InitNotificationRepo(db)
	set.InitSetRepo(db)
	story.InitStoryRepo(db)
	search.InitSearchRepo(db)
	follow.InitFollowRepo(db)
	hashtag.InitHashtagRepo(db)
	ticket.InitTicketRepo(db)
	guest.InitGuestRepo(db)
	staff.InitStaffRepo(db)
	admission.InitAdmissionRepo(db, staff.StaffRepo)
}
