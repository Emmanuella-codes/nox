import { FollowListScreen } from "@/src/components/user/profile/follow-list-screen";

interface FollowersPageProps {
  params: Promise<{ personaID: string }>;
}

export default async function FollowersPage({ params }: FollowersPageProps) {
  const { personaID } = await params;
  return <FollowListScreen personaID={personaID} mode="followers" />;
}
