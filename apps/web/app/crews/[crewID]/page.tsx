import { CrewMapScreen } from "@/src/components/user/crew/crew-map-screen";

interface CrewPageProps {
  params: Promise<{ crewID: string }>;
}

export default async function CrewPage({ params }: CrewPageProps) {
  const { crewID } = await params;
  return <CrewMapScreen crewID={crewID} />;
}
