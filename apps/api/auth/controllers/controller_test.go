package controllers

import (
	"testing"

	"github.com/emmanuella-codes/nox/auth/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/gofiber/fiber/v2"
)

func TestPipeErrorStatusMapsKnownAuthErrors(t *testing.T) {
	tests := []struct {
		message shared.PipeMessage
		status  int
	}{
		{messages.User_Already_Exists, fiber.StatusConflict},
		{messages.Invalid_Credentials, fiber.StatusUnauthorized},
		{messages.Invalid_Token, fiber.StatusUnauthorized},
		{messages.Email_Not_Verified, fiber.StatusForbidden},
		{messages.Email_Already_Verified, fiber.StatusConflict},
		{messages.Invalid_OTP, fiber.StatusUnauthorized},
		{messages.OTP_Expired, fiber.StatusUnauthorized},
		{messages.OTP_Locked, fiber.StatusTooManyRequests},
		{messages.Internal_Error, fiber.StatusInternalServerError},
	}

	for _, tt := range tests {
		if got := pipeErrorStatus(tt.message); got != tt.status {
			t.Fatalf("pipeErrorStatus(%q) = %d, want %d", tt.message, got, tt.status)
		}
	}
}

func TestPipeErrorStatusDefaultsToBadRequest(t *testing.T) {
	if got := pipeErrorStatus(shared.CreatePipeMessage("future_business_error")); got != fiber.StatusBadRequest {
		t.Fatalf("expected unknown pipe errors to map to %d, got %d", fiber.StatusBadRequest, got)
	}
}
