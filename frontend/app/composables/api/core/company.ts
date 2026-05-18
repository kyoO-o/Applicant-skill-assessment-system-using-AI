import type { Company } from "../../types";
import type { SaveCompanyPayload } from "../../types/payload";

export class CompanyAPI {
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

  get() {
    return this.fetch<Company>("/api/company");
  }

  save(payload: SaveCompanyPayload) {
    return this.fetch<Company>("/api/company", {
      method: "PUT",
      body: payload,
    });
  }

  uploadLogo(file: File) {
    const form = new FormData();
    form.set("file", file);
    return this.fetch<Company>("/api/company/logo", { method: "POST", body: form });
  }
}

export const useCompanyAPI = () => {
  const config = useRuntimeConfig();
  return new CompanyAPI(config.public.apiBase);
};
