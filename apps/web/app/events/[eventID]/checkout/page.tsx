import { CheckoutScreen } from "@/src/components/user/events/checkout-screen";

interface CheckoutPageProps {
  params: Promise<{ eventID: string }>;
}

export default async function CheckoutPage({ params }: CheckoutPageProps) {
  const { eventID } = await params;
  return <CheckoutScreen eventID={eventID} />;
}
