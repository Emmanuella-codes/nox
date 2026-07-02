import type { Conversation } from "@/src/types/api/messaging";

export function otherMemberID(conversation: Conversation, activePersonaID: string) {
  return conversation.members.find((member) => member.persona_id !== activePersonaID)?.persona_id ?? "";
}

export function conversationName(conversation: Conversation, activePersonaID: string) {
	if (conversation.conversation_type === "group") {
		return conversation.title || "group chat";
	}

	const otherID = otherMemberID(conversation, activePersonaID);
	const persona = conversation.members.find((member) => member.persona_id === otherID)?.persona;
	return persona?.display_name || persona?.handle || "direct message";
}

export function conversationHandle(conversation: Conversation, activePersonaID: string) {
	if (conversation.conversation_type === "group") {
		return `${conversation.members.length} members`;
	}

	const otherID = otherMemberID(conversation, activePersonaID);
	const persona = conversation.members.find((member) => member.persona_id === otherID)?.persona;
	return persona?.handle ? `@${persona.handle}` : "conversation";
}

export function lastMessagePreview(conversation: Conversation) {
  if (!conversation.last_message) {
    return "No messages yet";
  }
  if (conversation.last_message.deleted) {
    return "Message deleted";
  }
  if (conversation.last_message.body) {
    return conversation.last_message.body;
  }
  return conversation.last_message.message_type === "video" ? "Video" : "Image";
}
