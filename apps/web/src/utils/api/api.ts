export interface ApiResponse<T = unknown> {
    success: boolean;
    message: string;
    data?: T;
}
  
export class ApiRequestError extends Error {
    constructor(message: string) {
      super(message);
      this.name = "ApiRequestError";
    }
}

interface ApiRequestOptions {
    requiresAuth?: boolean;
    retryOnUnauthorized?: boolean;
}
  
async function parseApiResponse<TResponse>(
    response: Response,
): Promise<ApiResponse<TResponse>> {
    return (await response.json()) as ApiResponse<TResponse>;
}
