import { EventDetailScreen } from "@/src/components/user/events/event-detail-screen";

interface EventPageProps {
  params: Promise<{ eventID: string }>;
}

export default async function EventPage({ params }: EventPageProps) {
  const { eventID } = await params;
  return <EventDetailScreen eventID={eventID} />;
}
