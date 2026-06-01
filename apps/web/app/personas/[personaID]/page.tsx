import { PersonaProfileScreen } from "@/src/components/user/profile/persona-profile-screen";

interface PersonaPageProps {
  params: Promise<{ personaID: string }>;
}

export default async function PersonaPage({ params }: PersonaPageProps) {
  const { personaID } = await params;
  return <PersonaProfileScreen personaID={personaID} />;
}
