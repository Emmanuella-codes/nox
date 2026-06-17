import { ThreadScreen } from "@/src/components/user/messages/thread-screen";

interface ThreadPageProps {
  params: Promise<{ conversationID: string }>;
}

export default async function ThreadPage({ params }: ThreadPageProps) {
  const { conversationID } = await params;
  return <ThreadScreen conversationID={conversationID} />;
}
