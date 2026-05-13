export class IntegrationsAPI {
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

  googleCalendarStatus() {
    return this.fetch<{ connected: boolean }>("/api/integrations/google-calendar/status");
  }

  googleCalendarDisconnect() {
    return this.fetch<{ message: string }>("/api/integrations/google-calendar/disconnect", {
      method: "DELETE",
    });
  }

  googleCalendarConnectURL() {
    return `${this.baseURL}/api/integrations/google-calendar/connect`;
  }
}

export const useIntegrationsAPI = () => {
  const config = useRuntimeConfig();
  return new IntegrationsAPI(config.public.apiBase);
};
