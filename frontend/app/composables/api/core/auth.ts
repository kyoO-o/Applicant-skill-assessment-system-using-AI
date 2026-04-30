import type { User } from "../../types";
import type {
  AuthResponse,
  LoginPayload,
  LogoutResponse,
  RegisterPayload,
} from "../../types/payload";

export class AuthAPI {
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

  me() {
    return this.fetch<User>("/api/me");
  }

  login(payload: LoginPayload) {
    return this.fetch<AuthResponse>("/pub/login", {
      method: "POST",
      body: payload,
    });
  }

  register(payload: RegisterPayload) {
    return this.fetch<AuthResponse>("/pub/register", {
      method: "POST",
      body: payload,
    });
  }

  logout() {
    return this.fetch<LogoutResponse>("/api/logout");
  }
}

export const useAuthAPI = () => {
  const config = useRuntimeConfig();
  return new AuthAPI(config.public.apiBase);
};
