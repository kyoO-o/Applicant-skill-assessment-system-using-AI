export interface APIResult<T, M> {
  data: Awaited<T> | null;
  status: "idle" | "success" | "error" | "pending";
  load: (...args: Parameters<M>) => Promise<void>;
  error: any;
}
