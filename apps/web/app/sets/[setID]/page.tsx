import { SetDetailScreen } from "@/src/components/user/sets/set-detail-screen";

interface SetPageProps {
  params: Promise<{ setID: string }>;
}

export default async function SetPage({ params }: SetPageProps) {
  const { setID } = await params;
  return <SetDetailScreen setID={setID} />;
}
