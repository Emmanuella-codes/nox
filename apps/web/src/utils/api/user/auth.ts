import type {
  LoginRequest,
  LoginResponse,
  ResendVerificationRequest,
  SignupRequest,
  VerificationResponse,
  VerifyEmailRequest,
} from "@/src/types/api/user/auth";
import { apiRequest } from "@/src/utils/api/api";

export function login(payload: LoginRequest) {
  return apiRequest<LoginResponse, LoginRequest>("/auth/login", {
    method: "POST",
    body: payload,
  });
}

export function signup(payload: SignupRequest) {
  return apiRequest<VerificationResponse, SignupRequest>("/auth/register", {
    method: "POST",
    body: payload,
  });
}

export function verifyEmail(payload: VerifyEmailRequest) {
  return apiRequest<null, VerifyEmailRequest>("/auth/verify-email", {
    method: "POST",
    body: payload,
  });
}

export function resendVerification(payload: ResendVerificationRequest) {
  return apiRequest<VerificationResponse, ResendVerificationRequest>(
    "/auth/resend-verification",
    { method: "POST", body: payload },
  );
}
