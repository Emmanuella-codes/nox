import type { Comment, CreateCommentRequest } from "@/src/types/api/comment";
import { apiRequest } from "@/src/utils/api/api";

export function getPostComments(postID: string) {
  return apiRequest<Comment[]>(`/posts/${postID}/comments`);
}

export function createComment(postID: string, payload: CreateCommentRequest, token: string) {
  return apiRequest<Comment, CreateCommentRequest>(`/posts/${postID}/comments`, {
    method: "POST",
    body: payload,
    token,
  });
}
