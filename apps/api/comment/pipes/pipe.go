package pipes

import (
	"github.com/emmanuella-codes/nox/comment/messages"
	comment_repo "github.com/emmanuella-codes/nox/repositories/comment"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
)

type CommentPipe struct {
	commentRepo comment_repo.CommentRepository
	personaRepo persona_repo.PersonaRepository
	postRepo    post_repo.PostRepository
}

func NewCommentPipe(commentRepo comment_repo.CommentRepository, personaRepo persona_repo.PersonaRepository, postRepo post_repo.PostRepository) *CommentPipe {
	return &CommentPipe{commentRepo: commentRepo, personaRepo: personaRepo, postRepo: postRepo}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "comment", operation, messages.Internal_Error)
}
