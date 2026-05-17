import type { ApiResponse } from "@/src/types/api/auth";

export class ApiRequestError extends Error {
  status: number;
  response?: ApiResponse;

  constructor(message: string, status: number, response?: ApiResponse) {
    super(message);
    this.name = "ApiRequestError";
    this.status = status;
    this.response = response;
  }
}

interface ApiRequestOptions<TBody = unknown> {
  method?: "GET" | "POST" | "PATCH" | "DELETE";
  body?: TBody;
  headers?: HeadersInit;
  token?: string;
}

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") ?? "http://localhost:4006/api/v1";

async function parseApiResponse<TResponse>(
  response: Response,
): Promise<ApiResponse<TResponse>> {
  return (await response.json()) as ApiResponse<TResponse>;
}

export async function apiRequest<TResponse, TBody = unknown>(
  path: string,
  options: ApiRequestOptions<TBody> = {},
): Promise<ApiResponse<TResponse>> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: options.method ?? "GET",
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
      ...(options.token ? { Authorization: `Bearer ${options.token}` } : {}),
    },
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  const parsed = await parseApiResponse<TResponse>(response);
  if (!response.ok || !parsed.success) {
    throw new ApiRequestError(parsed.message || "request_failed", response.status, parsed);
  }

  return parsed;
}
