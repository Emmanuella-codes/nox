package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/like/dtos"
	"github.com/emmanuella-codes/nox/like/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *LikePipe) UnlikePostPipe(ctx context.Context, userID uuid.UUID, postID uuid.UUID, dto dtos.LikePostDTO) *shared.PipeRes[any] {
	if _, _, res := p.validatePersonaAndPost(ctx, userID, postID, dto.PersonaID); res != nil {
		return res
	}
	if err := p.likeRepo.UnlikePost(ctx, dto.PersonaID, postID); err != nil {
		return pipeInternalError[any](err, "unlike.post")
	}
	return shared.PipeSuccess[any](messages.Post_Unliked, nil)
}
