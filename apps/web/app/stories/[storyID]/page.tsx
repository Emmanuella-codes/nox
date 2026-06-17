import { StoryViewerScreen } from "@/src/components/user/stories/story-viewer-screen";

interface StoryPageProps {
  params: Promise<{ storyID: string }>;
}

export default async function StoryPage({ params }: StoryPageProps) {
  const { storyID } = await params;
  return <StoryViewerScreen storyID={storyID} />;
}
