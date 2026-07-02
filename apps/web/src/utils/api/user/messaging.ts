import type {
  AddMembersRequest,
  Conversation,
  CreateDirectConversationRequest,
  CreateGroupConversationRequest,
  MarkReadRequest,
  Message,
  RemoveMemberRequest,
  SendMessageRequest,
} from "@/src/types/api/messaging";
import { apiRequest } from "@/src/utils/api/api";

export function createDirectConversation(payload: CreateDirectConversationRequest, token: string) {
  return apiRequest<Conversation, CreateDirectConversationRequest>("/conversations/direct", {
    method: "POST",
    body: payload,
    token,
  });
}

export function createGroupConversation(payload: CreateGroupConversationRequest, token: string) {
  return apiRequest<Conversation, CreateGroupConversationRequest>("/conversations/group", {
    method: "POST",
    body: payload,
    token,
  });
}

export function listConversations(personaID: string, token: string, limit = 20, offset = 0) {
  const params = new URLSearchParams({
    persona_id: personaID,
    limit: String(limit),
    offset: String(offset),
  });
  return apiRequest<Conversation[]>(`/conversations?${params.toString()}`, { token });
}

export function getConversation(conversationID: string, token: string) {
  return apiRequest<Conversation>(`/conversations/${conversationID}`, { token });
}

export function listMessages(conversationID: string, personaID: string, token: string, limit = 30, offset = 0) {
  const params = new URLSearchParams({
    persona_id: personaID,
    limit: String(limit),
    offset: String(offset),
  });
  return apiRequest<Message[]>(`/conversations/${conversationID}/messages?${params.toString()}`, { token });
}

export function sendMessage(conversationID: string, payload: SendMessageRequest, token: string) {
  return apiRequest<Message, SendMessageRequest>(`/conversations/${conversationID}/messages`, {
    method: "POST",
    body: payload,
    token,
  });
}

export function markConversationRead(conversationID: string, payload: MarkReadRequest, token: string) {
  return apiRequest(`/conversations/${conversationID}/read`, {
    method: "POST",
    body: payload,
    token,
  });
}

export function addConversationMembers(conversationID: string, payload: AddMembersRequest, token: string) {
  return apiRequest(`/conversations/${conversationID}/members`, {
    method: "POST",
    body: payload,
    token,
  });
}

export function removeConversationMember(conversationID: string, personaID: string, payload: RemoveMemberRequest, token: string) {
  return apiRequest(`/conversations/${conversationID}/members/${personaID}`, {
    method: "DELETE",
    body: payload,
    token,
  });
}

export function deleteMessage(messageID: string, token: string) {
  return apiRequest<Message>(`/messages/${messageID}`, {
    method: "DELETE",
    token,
  });
}
