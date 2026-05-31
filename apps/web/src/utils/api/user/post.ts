import type { CreatePostRequest, Post } from "@/src/types/api/post";
import { apiRequest } from "@/src/utils/api/api";

export function createPost(payload: CreatePostRequest, token: string) {
  return apiRequest<Post, CreatePostRequest>("/posts", {
    method: "POST",
    body: payload,
    token,
  });
}

export function getPost(postID: string) {
  return apiRequest<Post>(`/posts/${postID}`);
}

export function getPostForViewer(postID: string, personaID: string, token: string) {
  return apiRequest<Post>(`/posts/${postID}/viewer?persona_id=${encodeURIComponent(personaID)}`, {
    token,
  });
}

export function getPersonaFeed(personaID: string, token: string) {
  return apiRequest<Post[]>(`/posts/persona/${personaID}/feed`, { token });
}

export function getFollowingFeed(personaID: string, token: string) {
  return apiRequest<Post[]>(`/posts/persona/${personaID}/following-feed`, { token });
}

export function getPersonaPosts(personaID: string, token?: string, viewerPersonaID?: string) {
  const params = new URLSearchParams();
  if (viewerPersonaID) {
    params.set("viewer_persona_id", viewerPersonaID);
  }
  const query = params.toString();
  return apiRequest<Post[]>(`/posts/persona/${personaID}${query ? `?${query}` : ""}`, { token });
}

export function likePost(postID: string, personaID: string, token: string) {
  return apiRequest<null, { persona_id: string }>(`/posts/${postID}/likes`, {
    method: "POST",
    body: { persona_id: personaID },
    token,
  });
}

export function unlikePost(postID: string, personaID: string, token: string) {
  return apiRequest<null, { persona_id: string }>(`/posts/${postID}/likes`, {
    method: "DELETE",
    body: { persona_id: personaID },
    token,
  });
}
