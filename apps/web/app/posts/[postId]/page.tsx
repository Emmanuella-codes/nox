import { SinglePostScreen } from "@/src/components/user/feed/single-post-screen";

interface PostPageProps {
  params: Promise<{ postId: string }>;
}

export default async function PostPage({ params }: PostPageProps) {
  const { postId } = await params;
  return <SinglePostScreen postId={postId} />;
}
