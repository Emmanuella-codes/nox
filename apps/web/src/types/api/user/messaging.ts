import type { MediaAsset } from "@/src/types/api/user/media";

export type ConversationType = "direct" | "group";
export type ConversationRole = "member" | "admin";
export type MessageType = "text" | "image" | "video" | "system";

export interface ConversationMember {
  persona_id: string;
  persona?: {
    id: string;
    handle: string;
    display_name: string;
    avatar_url: string;
  };
  role: ConversationRole;
  last_read_message_id?: string;
  joined_at: string;
}

export interface Message {
  id: string;
  conversation_id: string;
  sender_persona_id: string;
  body: string;
  message_type: MessageType;
  media_asset_id?: string;
  media?: MediaAsset;
  deleted: boolean;
  created_at: string;
  edited_at?: string;
}

export interface Conversation {
  id: string;
  conversation_type: ConversationType;
  title: string;
  created_by: string;
  last_message_id?: string;
  members: ConversationMember[];
  last_message?: Message;
  unread_count: number;
  created_at: string;
  updated_at: string;
}

export interface CreateDirectConversationRequest {
  sender_persona_id: string;
  recipient_persona_id: string;
}

export interface CreateGroupConversationRequest {
  creator_persona_id: string;
  title: string;
  member_persona_ids: string[];
}

export interface SendMessageRequest {
  sender_persona_id: string;
  body: string;
  message_type: MessageType;
  media_asset_id?: string;
}

export interface MarkReadRequest {
  persona_id: string;
  message_id: string;
}

export interface AddMembersRequest {
  admin_persona_id: string;
  member_persona_ids: string[];
}

export interface RemoveMemberRequest {
  admin_persona_id: string;
}
