import type { Job } from "../../types";
import type { DeleteJobResponse, SaveJobPayload } from "../../types/payload";

export class JobsAPI {
  baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async fetch<T>(
    path: string,
    options: Parameters<typeof $fetch<T>>[1] = {},
  ) {
    const headers = new Headers(options?.headers as HeadersInit | undefined);

    if (import.meta.server) {
      const requestHeaders = useRequestHeaders(["cookie"]);

      if (requestHeaders.cookie && !headers.has("cookie")) {
        headers.set("cookie", requestHeaders.cookie);
      }
    }

    return $fetch<T>(`${this.baseURL}${path}`, {
      ...options,
      credentials: "include",
      headers,
    });
  }

  list() {
    return this.fetch<Job[]>("/api/jobs");
  }

  listByCompany(companyID: number) {
    return this.fetch<Job[]>(`/api/jobs/${companyID}`);
  }

  get(companyID: number, id: number) {
    return this.fetch<Job>(`/api/jobs/${companyID}/${id}`);
  }

  create(payload: SaveJobPayload) {
    return this.fetch<Job>("/api/jobs", {
      method: "POST",
      body: payload,
    });
  }

  update(companyID: number, id: number, payload: SaveJobPayload) {
    return this.fetch<Job>(`/api/jobs/${companyID}/${id}`, {
      method: "PUT",
      body: payload,
    });
  }

  delete(companyID: number, id: number) {
    return this.fetch<DeleteJobResponse>(`/api/jobs/${companyID}/${id}`, {
      method: "DELETE",
    });
  }
}

export const useJobsAPI = () => {
  const config = useRuntimeConfig();
  return new JobsAPI(config.public.apiBase);
};
