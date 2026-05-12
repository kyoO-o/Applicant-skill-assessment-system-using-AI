export interface ChatMessage {
  role: "user" | "assistant";
  content: string;
}

export class ChatAPI {
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

  send(message: string, history: ChatMessage[]) {
    return this.fetch<{ reply: string }>("/api/chat", {
      method: "POST",
      body: { message, history },
    });
  }
}

export const useChatAPI = () => {
  const config = useRuntimeConfig();
  return new ChatAPI(config.public.apiBase);
};
