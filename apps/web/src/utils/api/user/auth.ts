import type {
  LoginRequest,
  LoginResponse,
  SignupRequest,
  VerificationResponse,
} from "@/src/types/api/auth";
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
