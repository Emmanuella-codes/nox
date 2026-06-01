package follow

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestFollowRepositoryUpdatesCountsAndHandlesDuplicates(t *testing.T) {
	pool := testPool(t)
	repo := NewFollowRepository(pool)
	ctx := context.Background()

	followerID := createTestPersona(t, pool, "follow_follower")
	followingID := createTestPersona(t, pool, "follow_following")
	otherID := createTestPersona(t, pool, "follow_other")

	if err := repo.Follow(ctx, followerID, followingID); err != nil {
		t.Fatalf("expected follow success: %v", err)
	}
	if err := repo.Follow(ctx, followerID, followingID); !errors.Is(err, ErrAlreadyFollowing) {
		t.Fatalf("expected duplicate follow error, got %v", err)
	}

	assertFollowCounts(t, pool, followerID, 0, 1)
	assertFollowCounts(t, pool, followingID, 1, 0)

	following, err := repo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		t.Fatalf("expected follow status success: %v", err)
	}
	if !following {
		t.Fatal("expected follower to be following target")
	}

	followingIDs, err := repo.FindFollowingIDs(ctx, followerID, []uuid.UUID{followingID, otherID})
	if err != nil {
		t.Fatalf("expected following ids success: %v", err)
	}
	if !followingIDs[followingID] || followingIDs[otherID] {
		t.Fatalf("unexpected following ids: %v", followingIDs)
	}

	if err := repo.Unfollow(ctx, followerID, followingID); err != nil {
		t.Fatalf("expected unfollow success: %v", err)
	}
	if err := repo.Unfollow(ctx, followerID, followingID); !errors.Is(err, ErrNotFollowing) {
		t.Fatalf("expected not-following error, got %v", err)
	}

	assertFollowCounts(t, pool, followerID, 0, 0)
	assertFollowCounts(t, pool, followingID, 0, 0)
}

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("NOX_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("NOX_TEST_DATABASE_URL not set")
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func createTestPersona(t *testing.T, pool *pgxpool.Pool, handlePrefix string) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	userID := uuid.New()
	personaID := uuid.New()
	handle := handlePrefix + "_" + uuid.NewString()

	_, err := pool.Exec(ctx, `
		INSERT INTO users (id, fullname, email, password)
		VALUES ($1, $2, $3, $4)
	`, userID, "Test User", userID.String()+"@example.com", "password")
	if err != nil {
		t.Fatalf("insert test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userID)
	})

	_, err = pool.Exec(ctx, `
		INSERT INTO personas (id, user_id, handle, display_name, persona_type)
		VALUES ($1, $2, $3, $4, 'visible')
	`, personaID, userID, handle, "Test Persona")
	if err != nil {
		t.Fatalf("insert test persona: %v", err)
	}
	return personaID
}

func assertFollowCounts(t *testing.T, pool *pgxpool.Pool, personaID uuid.UUID, wantFollowers int, wantFollowing int) {
	t.Helper()
	var followers int
	var following int
	err := pool.QueryRow(context.Background(), `
		SELECT follower_count, following_count
		FROM personas
		WHERE id = $1
	`, personaID).Scan(&followers, &following)
	if err != nil {
		t.Fatalf("read follow counts: %v", err)
	}
	if followers != wantFollowers || following != wantFollowing {
		t.Fatalf("expected counts followers=%d following=%d, got followers=%d following=%d", wantFollowers, wantFollowing, followers, following)
	}
}
