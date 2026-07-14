export interface MediaAsset {
  id: string;
  owner_persona_id: string;
  media_kind: "image" | "video";
  playback_url: string;
  thumbnail_url: string;
  mime_type: string;
  duration_seconds: number;
  size_bytes: number;
  processing_status: "pending" | "ready" | "failed";
  created_at: string;
  updated_at: string;
}

export interface InitiatePostMediaUploadRequest {
  owner_persona_id: string;
  media_kind: "image" | "video";
  mime_type: string;
  size_bytes: number;
}

export interface CloudinaryPostUploadResponse {
  cloud_name: string;
  api_key: string;
  resource_type: "image" | "video";
  upload_url: string;
  public_id: string;
  folder: string;
  timestamp: number;
  signature: string;
  params: Record<string, string>;
}

export interface ConfirmPostMediaUploadRequest {
  owner_persona_id: string;
  media_kind: "image" | "video";
  public_id: string;
  secure_url: string;
  thumbnail_url?: string;
  mime_type: string;
  duration_seconds?: number;
  size_bytes: number;
}
