package repositories

import (
	admin_repo "github.com/emmanuella-codes/nox/repositories/admin"
	admission_repo "github.com/emmanuella-codes/nox/repositories/admission"
	comment_repo "github.com/emmanuella-codes/nox/repositories/comment"
	crew_repo "github.com/emmanuella-codes/nox/repositories/crew"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	guest_repo "github.com/emmanuella-codes/nox/repositories/guest"
	hashtag_repo "github.com/emmanuella-codes/nox/repositories/hashtag"
	like_repo "github.com/emmanuella-codes/nox/repositories/like"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	messaging_repo "github.com/emmanuella-codes/nox/repositories/messaging"
	notification_repo "github.com/emmanuella-codes/nox/repositories/notification"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	preference_repo "github.com/emmanuella-codes/nox/repositories/preference"
	search_repo "github.com/emmanuella-codes/nox/repositories/search"
	set_repo "github.com/emmanuella-codes/nox/repositories/set"
	staff_repo "github.com/emmanuella-codes/nox/repositories/staff"
	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	ticket_repo "github.com/emmanuella-codes/nox/repositories/ticket"
	user_repo "github.com/emmanuella-codes/nox/repositories/user"
)

type AppRepositories struct {
	Admin        admin_repo.AdminRepository
	Admission    admission_repo.AdmissionRepository
	Comment      comment_repo.CommentRepository
	Crew         crew_repo.CrewRepository
	Event        event_repo.EventRepository
	Follow       follow_repo.FollowRepository
	Guest        guest_repo.GuestRepository
	Hashtag      hashtag_repo.HashtagRepository
	Like         like_repo.LikeRepository
	Media        media_repo.MediaRepository
	Messaging    messaging_repo.MessagingRepository
	Notification notification_repo.NotificationRepository
	Persona      persona_repo.PersonaRepository
	Post         post_repo.PostRepository
	Preference   preference_repo.PreferenceRepository
	Search       search_repo.SearchRepository
	Set          set_repo.SetRepository
	Staff        staff_repo.StaffRepository
	Story        story_repo.StoryRepository
	Ticket       ticket_repo.TicketRepository
	User         user_repo.UserRepository
}

func Current() AppRepositories {
	return AppRepositories{
		Admin:        admin_repo.AdminRepo,
		Admission:    admission_repo.AdmissionRepo,
		Comment:      comment_repo.CommentRepo,
		Crew:         crew_repo.CrewRepo,
		Event:        event_repo.EventRepo,
		Follow:       follow_repo.FollowRepo,
		Guest:        guest_repo.GuestRepo,
		Hashtag:      hashtag_repo.HashtagRepo,
		Like:         like_repo.LikeRepo,
		Media:        media_repo.MediaRepo,
		Messaging:    messaging_repo.MessagingRepo,
		Notification: notification_repo.NotificationRepo,
		Persona:      persona_repo.PersonaRepo,
		Post:         post_repo.PostRepo,
		Preference:   preference_repo.PreferenceRepo,
		Search:       search_repo.SearchRepo,
		Set:          set_repo.SetRepo,
		Staff:        staff_repo.StaffRepo,
		Story:        story_repo.StoryRepo,
		Ticket:       ticket_repo.TicketRepo,
		User:         user_repo.UserRepo,
	}
}
