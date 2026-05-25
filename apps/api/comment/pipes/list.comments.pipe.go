package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/comment/messages"
	"github.com/emmanuella-codes/nox/models"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *CommentPipe) ListCommentsPipe(ctx context.Context, postID uuid.UUID, limit int) *shared.PipeRes[[]*models.Comment] {
	if _, err := p.postRepo.FindPostByID(ctx, postID); err != nil {
		if err == post_repo.ErrPostNotFound {
			return pipeError[[]*models.Comment](messages.Post_Not_Found)
		}
		return pipeInternalError[[]*models.Comment](err, "comment.find_post")
	}

	comments, err := p.commentRepo.FindCommentsByPostID(ctx, postID, limit)
	if err != nil {
		return pipeInternalError[[]*models.Comment](err, "comment.list")
	}

	return pipeSuccess(messages.Comments_Listed, &comments)
}
