import { HashtagScreen } from "@/src/components/user/hashtag/hashtag-screen";

interface HashtagPageProps {
  params: Promise<{ tag: string }>;
}

export default async function HashtagPage({ params }: HashtagPageProps) {
  const { tag } = await params;
  return <HashtagScreen tag={decodeURIComponent(tag)} />;
}
