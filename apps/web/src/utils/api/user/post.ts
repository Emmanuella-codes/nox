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

export function getPersonaFeed(personaID: string) {
  return apiRequest<Post[]>(`/posts/persona/${personaID}/feed`);
}

export function getPersonaPosts(personaID: string) {
  return apiRequest<Post[]>(`/posts/persona/${personaID}`);
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
