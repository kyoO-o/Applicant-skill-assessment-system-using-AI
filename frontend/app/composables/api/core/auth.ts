import type { User } from "../../types";
import type {
  AuthResponse,
  LoginPayload,
  LogoutResponse,
  MessageResponse,
  RegisterPayload,
  RegisterResponse,
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
    return $fetch<T>(`${this.baseURL}${path}`, { ...options, credentials: "include", headers });
  }

  me() {
    return this.fetch<User>("/api/me");
  }

  login(payload: LoginPayload) {
    return this.fetch<AuthResponse>("/pub/login", { method: "POST", body: payload });
  }

  register(payload: RegisterPayload) {
    return this.fetch<RegisterResponse>("/pub/register", { method: "POST", body: payload });
  }

  verifyEmail(email: string, code: string) {
    return this.fetch<AuthResponse>("/pub/verify-email", { method: "POST", body: { email, code } });
  }

  resendVerification(email: string) {
    return this.fetch<MessageResponse>("/pub/resend-verification", { method: "POST", body: { email } });
  }

  forgotPassword(email: string) {
    return this.fetch<MessageResponse>("/pub/forgot-password", { method: "POST", body: { email } });
  }

  resetPassword(email: string, code: string, password: string) {
    return this.fetch<MessageResponse>("/pub/reset-password", {
      method: "POST",
      body: { email, code, password },
    });
  }

  updateMe(payload: { first_name: string; last_name: string }) {
    return this.fetch<User>("/api/me", { method: "PUT", body: payload });
  }

  uploadAvatar(file: File) {
    const form = new FormData();
    form.set("file", file);
    return this.fetch<User>("/api/me/avatar", { method: "POST", body: form });
  }

  initiateEmailChange(email: string) {
    return this.fetch<MessageResponse>("/api/me/email", { method: "PUT", body: { email } });
  }

  verifyEmailChange(code: string) {
    return this.fetch<User>("/api/me/verify-email-change", { method: "POST", body: { code } });
  }

  changePassword(currentPassword: string, newPassword: string) {
    return this.fetch<MessageResponse>("/api/me/password", {
      method: "PUT",
      body: { current_password: currentPassword, new_password: newPassword },
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
