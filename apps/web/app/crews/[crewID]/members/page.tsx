import { CrewMembersScreen } from "@/src/components/user/crew/crew-members-screen";

interface CrewMembersPageProps {
  params: Promise<{ crewID: string }>;
}

export default async function CrewMembersPage({ params }: CrewMembersPageProps) {
  const { crewID } = await params;
  return <CrewMembersScreen crewID={crewID} />;
}
