export type NotificationType = "like" | "follow" | "comment";

export interface SampleNotification {
  id: string;
  type: NotificationType;
  actor: string;
  actorID: string;
  text: string;
  time: string;
}
