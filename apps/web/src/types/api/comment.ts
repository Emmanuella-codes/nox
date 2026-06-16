export interface Comment {
  id: string;
  post_id: string;
  author: {
    mode: "public" | "anonymous";
    persona?: {
      id: string;
      handle: string;
      display_name: string;
      avatar_url: string;
    };
    anonymous?: {
      handle: string;
    };
  };
  body: string;
  parent_id?: string;
  like_count: number;
  created_at: string;
}

export interface CreateCommentRequest {
  persona_id: string;
  posting_mode?: "public" | "anonymous";
  body: string;
  parent_id?: string;
}
