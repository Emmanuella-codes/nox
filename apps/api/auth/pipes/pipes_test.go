package pipes

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/emmanuella-codes/nox/auth/dtos"
	"github.com/emmanuella-codes/nox/auth/messages"
	"github.com/emmanuella-codes/nox/auth/services"
	"github.com/emmanuella-codes/nox/config"
	"github.com/emmanuella-codes/nox/models"
	userrepo "github.com/emmanuella-codes/nox/repositories/user"
	"github.com/emmanuella-codes/nox/shared/mail"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestSignupPipeCreatesUserAndStoresVerificationOTP(t *testing.T) {
	ctx := context.Background()
	repo := newPipeTestUserRepo()
	emailProvider := &pipeTestMailProvider{}
	pipe, redisServer := newPipeTestAuthPipe(t, repo, emailProvider)

	res := pipe.SignupPipe(ctx, signupTestDTO(" Ada@Example.COM "))
	if !res.Success {
		t.Fatalf("expected signup success, got %q", res.Message)
	}
	if res.Message != messages.Verification_Sent {
		t.Fatalf("expected message %q, got %q", messages.Verification_Sent, res.Message)
	}
	if repo.createdFullname != "Ada Lovelace" {
		t.Fatalf("expected fullname to be stored, got %q", repo.createdFullname)
	}
	if repo.createdEmail != "ada@example.com" {
		t.Fatalf("expected normalized email to be stored, got %q", repo.createdEmail)
	}
	if len(emailProvider.messages) != 1 {
		t.Fatalf("expected one verification email, got %d", len(emailProvider.messages))
	}
	if !redisServer.Exists(emailVerificationKey(repo.createdUser.ID.String())) {
		t.Fatal("expected verification OTP hash to be stored in Redis")
	}
	if redisServer.Exists(emailVerificationAttemptsKey(repo.createdUser.ID.String())) {
		t.Fatal("expected previous verification attempts to be cleared")
	}
}

func TestSignupPipeReturnsAlreadyExistsWhenEmailExists(t *testing.T) {
	repo := newPipeTestUserRepo()
	repo.foundUser = pipeTestUser(t, "ada@example.com", "password123", true)
	pipe, _ := newPipeTestAuthPipe(t, repo, &pipeTestMailProvider{})

	res := pipe.SignupPipe(context.Background(), signupTestDTO("ada@example.com"))
	if res.Success {
		t.Fatal("expected signup to fail for existing user")
	}
	if res.Message != messages.User_Already_Exists {
		t.Fatalf("expected message %q, got %q", messages.User_Already_Exists, res.Message)
	}
}

func TestSignupPipeMapsUniqueConstraintRaceToAlreadyExists(t *testing.T) {
	repo := newPipeTestUserRepo()
	repo.createErr = userrepo.ErrUserAlreadyExists
	pipe, _ := newPipeTestAuthPipe(t, repo, &pipeTestMailProvider{})

	res := pipe.SignupPipe(context.Background(), signupTestDTO("ada@example.com"))
	if res.Success {
		t.Fatal("expected signup to fail when CreateUser reports duplicate email")
	}
	if res.Message != messages.User_Already_Exists {
		t.Fatalf("expected message %q, got %q", messages.User_Already_Exists, res.Message)
	}
}

func TestLoginPipeIssuesTokensForVerifiedUser(t *testing.T) {
	ctx := context.Background()
	repo := newPipeTestUserRepo()
	repo.foundUser = pipeTestUser(t, "ada@example.com", "password123", true)
	pipe, _ := newPipeTestAuthPipe(t, repo, &pipeTestMailProvider{})

	res := pipe.LoginPipe(ctx, dtos.LoginDTO{Email: "ADA@example.com", Password: "password123"})
	if !res.Success {
		t.Fatalf("expected login success, got %q", res.Message)
	}
	if res.Data == nil {
		t.Fatal("expected auth response data")
	}
	if res.Data.Tokens.AccessToken == "" || res.Data.Tokens.RefreshToken == "" || res.Data.Tokens.RefreshTokenID == "" {
		t.Fatal("expected access token, refresh token, and refresh token id")
	}
	sessionUserID, err := pipe.redis.Get(ctx, refreshSessionKey(res.Data.Tokens.RefreshTokenID)).Result()
	if err != nil {
		t.Fatalf("expected refresh session to be stored: %v", err)
	}
	if sessionUserID != repo.foundUser.ID.String() {
		t.Fatalf("expected session owner %s, got %s", repo.foundUser.ID, sessionUserID)
	}
}

func TestLoginPipeRejectsMissingInvalidAndUnverifiedUsers(t *testing.T) {
	ctx := context.Background()
	repo := newPipeTestUserRepo()
	pipe, _ := newPipeTestAuthPipe(t, repo, &pipeTestMailProvider{})

	missing := pipe.LoginPipe(ctx, dtos.LoginDTO{Email: "missing@example.com", Password: "password123"})
	if missing.Message != messages.Invalid_Credentials {
		t.Fatalf("expected missing user to return %q, got %q", messages.Invalid_Credentials, missing.Message)
	}

	repo.foundUser = pipeTestUser(t, "ada@example.com", "password123", true)
	invalidPassword := pipe.LoginPipe(ctx, dtos.LoginDTO{Email: "ada@example.com", Password: "wrong-password"})
	if invalidPassword.Message != messages.Invalid_Credentials {
		t.Fatalf("expected invalid password to return %q, got %q", messages.Invalid_Credentials, invalidPassword.Message)
	}

	repo.foundUser = pipeTestUser(t, "ada@example.com", "password123", false)
	unverified := pipe.LoginPipe(ctx, dtos.LoginDTO{Email: "ada@example.com", Password: "password123"})
	if unverified.Message != messages.Email_Not_Verified {
		t.Fatalf("expected unverified user to return %q, got %q", messages.Email_Not_Verified, unverified.Message)
	}
}

func TestRefreshPipeRotatesRefreshTokenSession(t *testing.T) {
	ctx := context.Background()
	repo := newPipeTestUserRepo()
	pipe, _ := newPipeTestAuthPipe(t, repo, &pipeTestMailProvider{})
	userID := uuid.New()
	oldTokenID := "old-refresh-id"
	oldRefreshToken, err := pipe.tokenService.SignRefreshToken(userID, oldTokenID)
	if err != nil {
		t.Fatalf("sign refresh token: %v", err)
	}
	if err := pipe.storeRefreshSession(ctx, userID, oldTokenID); err != nil {
		t.Fatalf("store old refresh session: %v", err)
	}

	res := pipe.RefreshPipe(ctx, oldRefreshToken)
	if !res.Success {
		t.Fatalf("expected refresh success, got %q", res.Message)
	}
	if pipe.redis.Exists(ctx, refreshSessionKey(oldTokenID)).Val() != 0 {
		t.Fatal("expected old refresh session to be deleted")
	}
	if res.Data == nil || res.Data.RefreshTokenID == "" {
		t.Fatal("expected new refresh token id")
	}
	if pipe.redis.Exists(ctx, refreshSessionKey(res.Data.RefreshTokenID)).Val() != 1 {
		t.Fatal("expected new refresh session to be stored")
	}
}

func TestRefreshPipeRejectsMissingSessionAndRevokesUserSessions(t *testing.T) {
	ctx := context.Background()
	repo := newPipeTestUserRepo()
	pipe, _ := newPipeTestAuthPipe(t, repo, &pipeTestMailProvider{})
	userID := uuid.New()
	missingTokenID := "missing-refresh-id"
	otherTokenID := "other-refresh-id"
	refreshToken, err := pipe.tokenService.SignRefreshToken(userID, missingTokenID)
	if err != nil {
		t.Fatalf("sign refresh token: %v", err)
	}
	if err := pipe.storeRefreshSession(ctx, userID, otherTokenID); err != nil {
		t.Fatalf("store other refresh session: %v", err)
	}

	res := pipe.RefreshPipe(ctx, refreshToken)
	if res.Message != messages.Invalid_Token {
		t.Fatalf("expected invalid token, got %q", res.Message)
	}
	if pipe.redis.Exists(ctx, refreshSessionKey(otherTokenID)).Val() != 0 {
		t.Fatal("expected user refresh sessions to be revoked after reuse signal")
	}
}

func TestLogoutPipeVerifiesTokenAndDeletesSession(t *testing.T) {
	ctx := context.Background()
	repo := newPipeTestUserRepo()
	pipe, _ := newPipeTestAuthPipe(t, repo, &pipeTestMailProvider{})
	userID := uuid.New()
	tokenID := "logout-refresh-id"
	refreshToken, err := pipe.tokenService.SignRefreshToken(userID, tokenID)
	if err != nil {
		t.Fatalf("sign refresh token: %v", err)
	}
	if err := pipe.storeRefreshSession(ctx, userID, tokenID); err != nil {
		t.Fatalf("store refresh session: %v", err)
	}

	res := pipe.LogoutPipe(ctx, refreshToken)
	if !res.Success {
		t.Fatalf("expected logout success, got %q", res.Message)
	}
	if pipe.redis.Exists(ctx, refreshSessionKey(tokenID)).Val() != 0 {
		t.Fatal("expected refresh session to be deleted on logout")
	}

	emptyToken := pipe.LogoutPipe(ctx, " ")
	if emptyToken.Message != messages.Invalid_Token {
		t.Fatalf("expected empty refresh token to return %q, got %q", messages.Invalid_Token, emptyToken.Message)
	}
}

func TestVerifyEmailPipeMarksUserVerifiedAndDeletesOTPState(t *testing.T) {
	ctx := context.Background()
	repo := newPipeTestUserRepo()
	repo.foundUser = pipeTestUser(t, "ada@example.com", "password123", false)
	pipe, _ := newPipeTestAuthPipe(t, repo, &pipeTestMailProvider{})
	otpHash, err := pipe.otpService.Hash("123456")
	if err != nil {
		t.Fatalf("hash otp: %v", err)
	}
	userID := repo.foundUser.ID.String()
	if err := pipe.redis.Set(ctx, emailVerificationKey(userID), otpHash, time.Minute).Err(); err != nil {
		t.Fatalf("store otp hash: %v", err)
	}
	if err := pipe.redis.Set(ctx, emailVerificationAttemptsKey(userID), "2", time.Minute).Err(); err != nil {
		t.Fatalf("store attempt count: %v", err)
	}

	res := pipe.VerifyEmailPipe(ctx, dtos.VerifyEmailDTO{Email: "ADA@example.com", OTP: "123456"})
	if !res.Success {
		t.Fatalf("expected verify email success, got %q", res.Message)
	}
	if repo.markedVerifiedUserID != userID {
		t.Fatalf("expected user %s to be marked verified, got %s", userID, repo.markedVerifiedUserID)
	}
	if pipe.redis.Exists(ctx, emailVerificationKey(userID), emailVerificationAttemptsKey(userID)).Val() != 0 {
		t.Fatal("expected OTP and attempt keys to be deleted")
	}
}

func TestVerifyEmailPipeLocksAfterMaxInvalidOTPAttempts(t *testing.T) {
	ctx := context.Background()
	repo := newPipeTestUserRepo()
	repo.foundUser = pipeTestUser(t, "ada@example.com", "password123", false)
	pipe, _ := newPipeTestAuthPipe(t, repo, &pipeTestMailProvider{})
	otpHash, err := pipe.otpService.Hash("123456")
	if err != nil {
		t.Fatalf("hash otp: %v", err)
	}
	userID := repo.foundUser.ID.String()
	if err := pipe.redis.Set(ctx, emailVerificationKey(userID), otpHash, time.Minute).Err(); err != nil {
		t.Fatalf("store otp hash: %v", err)
	}

	for attempt := 1; attempt < maxEmailVerificationAttempts; attempt++ {
		res := pipe.VerifyEmailPipe(ctx, dtos.VerifyEmailDTO{Email: "ada@example.com", OTP: "000000"})
		if res.Message != messages.Invalid_OTP {
			t.Fatalf("attempt %d expected %q, got %q", attempt, messages.Invalid_OTP, res.Message)
		}
	}

	res := pipe.VerifyEmailPipe(ctx, dtos.VerifyEmailDTO{Email: "ada@example.com", OTP: "000000"})
	if res.Message != messages.OTP_Locked {
		t.Fatalf("expected locked OTP to return %q, got %q", messages.OTP_Locked, res.Message)
	}
	if pipe.redis.Exists(ctx, emailVerificationKey(userID), emailVerificationAttemptsKey(userID)).Val() != 0 {
		t.Fatal("expected OTP and attempt keys to be deleted after lockout")
	}
}

func TestVerifyEmailPipeHandlesMissingAlreadyVerifiedAndExpiredOTP(t *testing.T) {
	ctx := context.Background()
	repo := newPipeTestUserRepo()
	pipe, _ := newPipeTestAuthPipe(t, repo, &pipeTestMailProvider{})

	missingUser := pipe.VerifyEmailPipe(ctx, dtos.VerifyEmailDTO{Email: "missing@example.com", OTP: "123456"})
	if missingUser.Message != messages.Invalid_Credentials {
		t.Fatalf("expected missing user to return %q, got %q", messages.Invalid_Credentials, missingUser.Message)
	}

	repo.foundUser = pipeTestUser(t, "ada@example.com", "password123", true)
	alreadyVerified := pipe.VerifyEmailPipe(ctx, dtos.VerifyEmailDTO{Email: "ada@example.com", OTP: "123456"})
	if alreadyVerified.Message != messages.Email_Already_Verified {
		t.Fatalf("expected already verified user to return %q, got %q", messages.Email_Already_Verified, alreadyVerified.Message)
	}

	repo.foundUser = pipeTestUser(t, "ada@example.com", "password123", false)
	expired := pipe.VerifyEmailPipe(ctx, dtos.VerifyEmailDTO{Email: "ada@example.com", OTP: "123456"})
	if expired.Message != messages.OTP_Expired {
		t.Fatalf("expected missing OTP to return %q, got %q", messages.OTP_Expired, expired.Message)
	}
}

func TestResendVerificationPipeStoresNewOTPAndSendsEmail(t *testing.T) {
	ctx := context.Background()
	repo := newPipeTestUserRepo()
	repo.foundUser = pipeTestUser(t, "ada@example.com", "password123", false)
	emailProvider := &pipeTestMailProvider{}
	pipe, redisServer := newPipeTestAuthPipe(t, repo, emailProvider)

	res := pipe.ResendVerificationPipe(ctx, dtos.ResendVerificationDTO{Email: "ADA@example.com"})
	if !res.Success {
		t.Fatalf("expected resend verification success, got %q", res.Message)
	}
	if len(emailProvider.messages) != 1 {
		t.Fatalf("expected one verification email, got %d", len(emailProvider.messages))
	}
	if !redisServer.Exists(emailVerificationKey(repo.foundUser.ID.String())) {
		t.Fatal("expected new OTP hash to be stored")
	}
}

func TestResendVerificationPipeRejectsMissingAndAlreadyVerifiedUsers(t *testing.T) {
	ctx := context.Background()
	repo := newPipeTestUserRepo()
	pipe, _ := newPipeTestAuthPipe(t, repo, &pipeTestMailProvider{})

	missing := pipe.ResendVerificationPipe(ctx, dtos.ResendVerificationDTO{Email: "missing@example.com"})
	if missing.Message != messages.Invalid_Credentials {
		t.Fatalf("expected missing user to return %q, got %q", messages.Invalid_Credentials, missing.Message)
	}

	repo.foundUser = pipeTestUser(t, "ada@example.com", "password123", true)
	alreadyVerified := pipe.ResendVerificationPipe(ctx, dtos.ResendVerificationDTO{Email: "ada@example.com"})
	if alreadyVerified.Message != messages.Email_Already_Verified {
		t.Fatalf("expected already verified user to return %q, got %q", messages.Email_Already_Verified, alreadyVerified.Message)
	}
}

func newPipeTestAuthPipe(t *testing.T, repo *pipeTestUserRepo, emailProvider *pipeTestMailProvider) (*AuthPipe, *miniredis.Miniredis) {
	t.Helper()

	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: redisServer.Addr()})
	t.Cleanup(func() {
		_ = redisClient.Close()
	})

	cfg := &config.Config{
		JWTAccessSecret:  "access-secret",
		JWTRefreshSecret: "refresh-secret",
		JWTIssuer:        "nox-api",
		JWTAudience:      "nox-client",
		JWTAccessTTL:     time.Minute,
		JWTRefreshTTL:    time.Hour,
		EmailOTPTTL:      10 * time.Minute,
	}

	hashService := services.NewHashService()
	return NewAuthPipe(AuthPipeDeps{
		UserRepo:     repo,
		HashService:  hashService,
		OTPService:   services.NewOTPService(),
		EmailService: services.NewEmailService(emailProvider),
		TokenService: services.NewTokenService(cfg),
		Redis:        redisClient,
		Config:       cfg,
	}), redisServer
}

func pipeTestUser(t *testing.T, email string, password string, verified bool) *models.User {
	t.Helper()
	hash, err := services.NewHashService().HashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	return &models.User{
		ID:            uuid.New(),
		Email:         email,
		Password:      hash,
		EmailVerified: verified,
	}
}

func signupTestDTO(email string) dtos.SignupDTO {
	return dtos.SignupDTO{
		Firstname: "Ada",
		Lastname:  "Lovelace",
		Email:     email,
		Password:  "password123",
	}
}

type pipeTestUserRepo struct {
	foundUser            *models.User
	findErr              error
	createErr            error
	markErr              error
	createdUser          *models.User
	createdFullname      string
	createdEmail         string
	markedVerifiedUserID string
}

func newPipeTestUserRepo() *pipeTestUserRepo {
	return &pipeTestUserRepo{}
}

func (r *pipeTestUserRepo) CreateUser(ctx context.Context, fullname string, email string, passwordHash string) (*models.User, error) {
	if r.createErr != nil {
		return nil, r.createErr
	}
	if r.createdUser == nil {
		r.createdUser = &models.User{
			ID:            uuid.New(),
			Fullname:      fullname,
			Email:         email,
			Password:      passwordHash,
			EmailVerified: false,
		}
	}
	r.createdFullname = fullname
	r.createdEmail = email
	return r.createdUser, nil
}

func (r *pipeTestUserRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	return r.foundUser, nil
}

func (r *pipeTestUserRepo) MarkEmailVerified(ctx context.Context, userID string) error {
	if r.markErr != nil {
		return r.markErr
	}
	if r.foundUser == nil || r.foundUser.ID.String() != userID {
		return errors.New("unexpected user id")
	}
	r.markedVerifiedUserID = userID
	r.foundUser.EmailVerified = true
	return nil
}

type pipeTestMailProvider struct {
	messages []mail.Message
	err      error
}

func (p *pipeTestMailProvider) Send(ctx context.Context, message mail.Message) error {
	if p.err != nil {
		return p.err
	}
	p.messages = append(p.messages, message)
	return nil
}
