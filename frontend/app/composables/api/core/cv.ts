import type { CVProfile } from "../../types";

export class CVAPI {
  baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async fetch<T>(path: string, options: Parameters<typeof $fetch<T>>[1] = {}) {
    const headers = new Headers(options?.headers as HeadersInit | undefined);
    if (import.meta.server) {
      const requestHeaders = useRequestHeaders(["cookie"]);
      if (requestHeaders.cookie && !headers.has("cookie")) {
        headers.set("cookie", requestHeaders.cookie);
      }
    }
    return $fetch<T>(`${this.baseURL}${path}`, { ...options, credentials: "include", headers });
  }

  get() {
    return this.fetch<CVProfile>("/api/cv-profile");
  }

  save(profile: CVProfile) {
    return this.fetch<CVProfile>("/api/cv-profile", {
      method: "PUT",
      body: profile,
    });
  }

  parse(file: File) {
    const form = new FormData();
    form.append("cv", file);
    return this.fetch<CVProfile>("/api/cv-profile/parse", {
      method: "POST",
      body: form,
    });
  }
}

export const useCVAPI = () => {
  const config = useRuntimeConfig();
  return new CVAPI(config.public.apiBase);
};
