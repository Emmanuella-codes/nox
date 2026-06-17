import type {
  CloudinaryPostUploadResponse,
  ConfirmPostMediaUploadRequest,
  InitiatePostMediaUploadRequest,
  MediaAsset,
} from "@/src/types/api/media";
import { apiRequest } from "@/src/utils/api/api";

export function initiatePostMediaUpload(payload: InitiatePostMediaUploadRequest, token: string) {
  return apiRequest<CloudinaryPostUploadResponse, InitiatePostMediaUploadRequest>("/media-assets/post/uploads", {
    method: "POST",
    body: payload,
    token,
  });
}

export function confirmPostMediaUpload(payload: ConfirmPostMediaUploadRequest, token: string) {
  return apiRequest<MediaAsset, ConfirmPostMediaUploadRequest>("/media-assets/post/confirm", {
    method: "POST",
    body: payload,
    token,
  });
}

export async function uploadToCloudinary(file: File, upload: CloudinaryPostUploadResponse) {
  const form = new FormData();
  form.set("file", file);
  for (const [key, value] of Object.entries(upload.params)) {
    form.set(key, value);
  }

  const response = await fetch(upload.upload_url, {
    method: "POST",
    body: form,
  });
  if (!response.ok) {
    throw new Error("cloudinary_upload_failed");
  }
  return response.json() as Promise<{
    public_id: string;
    secure_url: string;
    resource_type: "image" | "video";
    format: string;
    bytes: number;
    duration?: number;
    thumbnail_url?: string;
  }>;
}
