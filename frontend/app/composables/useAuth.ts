import type { User } from "./types";
import type {
  AuthResponse,
  LoginPayload,
  LogoutResponse,
  RegisterPayload,
} from "./types/payload";

export function useAuth() {
  const config = useRuntimeConfig();
  const apiBase = config.public.apiBase;

  const user = useState<User | null>("auth-user", () => null);
  const isLoading = useState<boolean>("auth-loading", () => false);

  async function authFetch<T>(
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

    return await $fetch<T>(`${apiBase}${path}`, {
      ...options,
      credentials: "include",
      headers,
    });
  }

  async function me() {
    isLoading.value = true;

    try {
      const currentUser = await authFetch<User>("/api/me");
      user.value = currentUser;
      return currentUser;
    } catch {
      user.value = null;
      return null;
    } finally {
      isLoading.value = false;
    }
  }

  async function login(payload: LoginPayload) {
    isLoading.value = true;

    try {
      const response = await authFetch<AuthResponse>("/pub/login", {
        method: "POST",
        body: payload,
      });
      user.value = response.user;
      return response;
    } finally {
      isLoading.value = false;
    }
  }

  async function register(payload: RegisterPayload) {
    isLoading.value = true;

    try {
      const response = await authFetch<AuthResponse>("/pub/register", {
        method: "POST",
        body: payload,
      });
      user.value = response.user;
      return response;
    } finally {
      isLoading.value = false;
    }
  }

  async function logout() {
    isLoading.value = true;

    try {
      await authFetch<LogoutResponse>("/api/logout");
      user.value = null;
    } finally {
      isLoading.value = false;
    }
  }

  function clearUser() {
    user.value = null;
  }

  return {
    user,
    isLoading,
    me,
    login,
    register,
    logout,
    clearUser,
  };
}
