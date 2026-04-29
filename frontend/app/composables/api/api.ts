import type { APIResult } from "./api_type";
import { ref, type UnwrapRef } from "vue";

export function useAPI<T, M extends (...args: any[]) => Promise<T>>(
  thisArg: any,
  loadData: M,
) {
  const result = ref<APIResult<T, M>>({
    data: null,
    status: "idle",
    error: null,
    load: async (...args: Parameters<M>) => {
      result.value.status = "pending";
      try {
        const d = await loadData.apply(thisArg, args);
        result.value.data = d as UnwrapRef<Awaited<T>>;
        result.value.status = "success";
      } catch (err) {
        console.log(err);
        result.value.status = "error";
        result.value.error = err;
      }
    },
  });
  return result;
}
