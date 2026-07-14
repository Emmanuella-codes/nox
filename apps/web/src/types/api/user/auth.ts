export interface ApiResponse<T = unknown> {
  success: boolean;
  message: string;
  data?: T;
  token?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface SignupRequest {
  firstname: string;
  lastname: string;
  email: string;
  password: string;
  category?: "patron" | "dj" | "organizer";
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
}

export interface LoginResponse {
  tokens: AuthTokens;
}

export interface VerificationResponse {
  expires_in_seconds: number;
}

export interface VerifyEmailRequest {
  email: string;
  otp: string;
}

export interface ResendVerificationRequest {
  email: string;
}
