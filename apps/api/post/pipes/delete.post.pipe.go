package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/post/messages"
	post_repo "github.com/emmanuella-codes/nox/repositories/post"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *PostPipe) DeletePostPipe(ctx context.Context, userID uuid.UUID, postID uuid.UUID) *shared.PipeRes[any] {
	post, err := p.postRepo.FindPostByID(ctx, postID)
	if err != nil {
		if err == post_repo.ErrPostNotFound {
			return shared.PipeError[any](messages.Post_Not_Found)
		}
		return pipeInternalError[any](err, "post.find_for_delete")
	}

	if post.AuthorUserID != userID {
		return shared.PipeError[any](messages.Forbidden)
	}

	if p.hashtagRepo != nil {
		if err := p.hashtagRepo.DeletePostHashtags(ctx, postID); err != nil {
			return pipeInternalError[any](err, "post.delete_hashtags")
		}
	}

	if err := p.postRepo.DeletePost(ctx, postID); err != nil {
		if err == post_repo.ErrPostNotFound {
			return shared.PipeError[any](messages.Post_Not_Found)
		}
		return pipeInternalError[any](err, "post.delete")
	}

	return shared.PipeSuccess[any](messages.Post_Deleted, nil)
}
