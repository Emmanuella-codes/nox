import { FollowListScreen } from "@/src/components/user/profile/follow-list-screen";

interface FollowingPageProps {
  params: Promise<{ personaID: string }>;
}

export default async function FollowingPage({ params }: FollowingPageProps) {
  const { personaID } = await params;
  return <FollowListScreen personaID={personaID} mode="following" />;
}
