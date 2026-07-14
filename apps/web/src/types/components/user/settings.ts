export type PaymentStatus = "success" | "refunded";

export interface PaymentHistoryItem {
  id: string;
  event: string;
  amount: number;
  qty: number;
  date: string;
  ref: string;
  status: PaymentStatus;
}
