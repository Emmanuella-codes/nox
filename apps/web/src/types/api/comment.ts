export interface Comment {
  id: string;
  persona_id: string;
  post_id: string;
  body: string;
  parent_id?: string;
  like_count: number;
  created_at: string;
}

export interface CreateCommentRequest {
  persona_id: string;
  body: string;
  parent_id?: string;
}
