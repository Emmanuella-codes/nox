import { HighlightsManageScreen } from "@/src/components/user/events/highlights-manage-screen";

interface HighlightsManagePageProps {
  params: Promise<{ eventID: string }>;
}

export default async function HighlightsManagePage({ params }: HighlightsManagePageProps) {
  const { eventID } = await params;
  return <HighlightsManageScreen eventID={eventID} />;
}
