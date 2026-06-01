package hashtag

import (
	"context"
	"os"
	"testing"

	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestHashtagRepositorySyncsSearchesAndDeletesCounts(t *testing.T) {
	pool := testPool(t)
	repo := NewHashtagRepository(pool)
	ctx := context.Background()

	postID := createTestPost(t, pool, "first #Amapiano post")
	secondPostID := createTestPost(t, pool, "second #Amapiano post")

	if err := repo.SyncPostHashtags(ctx, postID, []string{"Amapiano", "afro-house", "amapiano"}); err != nil {
		t.Fatalf("expected hashtag sync success: %v", err)
	}
	if err := repo.SyncPostHashtags(ctx, secondPostID, []string{"amapiano"}); err != nil {
		t.Fatalf("expected second hashtag sync success: %v", err)
	}

	amapiano, err := repo.FindByTag(ctx, "amapiano")
	if err != nil {
		t.Fatalf("expected find hashtag success: %v", err)
	}
	if amapiano == nil || amapiano.PostCount != 2 {
		t.Fatalf("expected amapiano count 2, got %+v", amapiano)
	}

	tagsByPost, err := repo.FindTagsByPostIDs(ctx, []uuid.UUID{postID})
	if err != nil {
		t.Fatalf("expected tags by post success: %v", err)
	}
	if got := tagsByPost[postID]; len(got) != 2 || got[0] != "amapiano" || got[1] != "afro-house" {
		t.Fatalf("unexpected post tags: %v", got)
	}

	posts, err := repo.FindPostsByTag(ctx, "amapiano", 1, 0)
	if err != nil {
		t.Fatalf("expected posts by tag success: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected paginated single post, got %d", len(posts))
	}

	results, err := repo.Search(ctx, "#ama", 10, 0)
	if err != nil {
		t.Fatalf("expected hashtag search success: %v", err)
	}
	if len(results) == 0 || results[0].Tag != "amapiano" {
		t.Fatalf("expected amapiano search result, got %+v", results)
	}

	if err := repo.DeletePostHashtags(ctx, postID); err != nil {
		t.Fatalf("expected delete hashtag links success: %v", err)
	}
	amapiano, err = repo.FindByTag(ctx, "amapiano")
	if err != nil {
		t.Fatalf("expected find hashtag after delete success: %v", err)
	}
	if amapiano == nil || amapiano.PostCount != 1 {
		t.Fatalf("expected amapiano count 1 after cleanup, got %+v", amapiano)
	}
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

func createTestPost(t *testing.T, pool *pgxpool.Pool, body string) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	userID := uuid.New()
	postID := uuid.New()

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
		INSERT INTO posts (id, author_user_id, posting_mode, body, post_type)
		VALUES ($1, $2, 'anonymous', $3, $4)
	`, postID, userID, body, models.TextPostType)
	if err != nil {
		t.Fatalf("insert test post: %v", err)
	}
	return postID
}
