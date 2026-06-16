import { CrewHubScreen } from "@/src/components/user/crew/crew-hub-screen";

interface EventCrewPageProps {
  params: Promise<{ eventID: string }>;
}

export default async function EventCrewPage({ params }: EventCrewPageProps) {
  const { eventID } = await params;
  return <CrewHubScreen eventID={eventID} />;
}
