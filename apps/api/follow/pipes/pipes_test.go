package pipes

import (
	"context"
	"testing"

	"github.com/emmanuella-codes/nox/follow/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/emmanuella-codes/nox/persona/dtos"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/google/uuid"
)

func TestFollowPersonaPipeRejectsSelfFollow(t *testing.T) {
	userID := uuid.New()
	personaID := uuid.New()
	pipe := NewFollowPipe(&followTestRepo{}, &followTestPersonaRepo{
		personas: map[string]*models.Persona{
			personaID.String(): {
				ID:          personaID,
				UserID:      userID,
				PersonaType: models.VisiblePersonaType,
			},
		},
	})

	res := pipe.FollowPersonaPipe(context.Background(), userID, personaID, personaID)
	if res.Message != messages.Self_Follow_Not_Allowed {
		t.Fatalf("expected %q, got %q", messages.Self_Follow_Not_Allowed, res.Message)
	}
}

func TestFollowPersonaPipeRejectsNonOwnedFollowerPersona(t *testing.T) {
	userID := uuid.New()
	followerID := uuid.New()
	targetID := uuid.New()
	pipe := NewFollowPipe(&followTestRepo{}, &followTestPersonaRepo{
		personas: map[string]*models.Persona{
			followerID.String(): {
				ID:          followerID,
				UserID:      uuid.New(),
				PersonaType: models.VisiblePersonaType,
			},
			targetID.String(): {
				ID:          targetID,
				UserID:      uuid.New(),
				PersonaType: models.VisiblePersonaType,
			},
		},
	})

	res := pipe.FollowPersonaPipe(context.Background(), userID, followerID, targetID)
	if res.Message != messages.Forbidden {
		t.Fatalf("expected %q, got %q", messages.Forbidden, res.Message)
	}
}

func TestFollowPersonaPipeMapsAlreadyFollowing(t *testing.T) {
	userID := uuid.New()
	followerID := uuid.New()
	targetID := uuid.New()
	pipe := NewFollowPipe(&followTestRepo{followErr: follow_repo.ErrAlreadyFollowing}, &followTestPersonaRepo{
		personas: map[string]*models.Persona{
			followerID.String(): {
				ID:          followerID,
				UserID:      userID,
				PersonaType: models.VisiblePersonaType,
			},
			targetID.String(): {
				ID:          targetID,
				UserID:      uuid.New(),
				PersonaType: models.VisiblePersonaType,
			},
		},
	})

	res := pipe.FollowPersonaPipe(context.Background(), userID, followerID, targetID)
	if res.Message != messages.Already_Following {
		t.Fatalf("expected %q, got %q", messages.Already_Following, res.Message)
	}
}

func TestFollowersPipeListsVisibleFollowers(t *testing.T) {
	targetID := uuid.New()
	followerID := uuid.New()
	pipe := NewFollowPipe(&followTestRepo{
		followers: []*models.Persona{{ID: followerID, PersonaType: models.VisiblePersonaType}},
	}, &followTestPersonaRepo{
		personas: map[string]*models.Persona{
			targetID.String(): {
				ID:          targetID,
				UserID:      uuid.New(),
				PersonaType: models.VisiblePersonaType,
			},
		},
	})

	res := pipe.FollowersPipe(context.Background(), targetID, 20)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if len(*res.Data) != 1 || (*res.Data)[0].ID != followerID {
		t.Fatalf("expected follower %s, got %+v", followerID, res.Data)
	}
}

type followTestRepo struct {
	followErr   error
	unfollowErr error
	followers   []*models.Persona
	following   []*models.Persona
	isFollowing bool
}

func (r *followTestRepo) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return r.followErr
}

func (r *followTestRepo) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return r.unfollowErr
}

func (r *followTestRepo) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	return r.isFollowing, nil
}

func (r *followTestRepo) FindFollowers(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Persona, error) {
	return r.followers, nil
}

func (r *followTestRepo) FindFollowing(ctx context.Context, personaID uuid.UUID, limit int) ([]*models.Persona, error) {
	return r.following, nil
}

type followTestPersonaRepo struct {
	personas map[string]*models.Persona
}

func (r *followTestPersonaRepo) CreatePersona(ctx context.Context, userID uuid.UUID, dto dtos.CreatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}

func (r *followTestPersonaRepo) FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error) {
	persona, ok := r.personas[personaID.String()]
	if !ok {
		return nil, persona_repo.ErrPersonaNotFound
	}
	return persona, nil
}

func (r *followTestPersonaRepo) FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error) {
	return nil, nil
}

func (r *followTestPersonaRepo) FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error) {
	return nil, persona_repo.ErrPersonaNotFound
}

func (r *followTestPersonaRepo) UpdatePersona(ctx context.Context, personaID uuid.UUID, dto dtos.UpdatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}
