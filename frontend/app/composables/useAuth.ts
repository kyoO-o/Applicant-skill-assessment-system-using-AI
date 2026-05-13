import type { User } from "./types";
import type {
  AuthResponse,
  LoginPayload,
  LogoutResponse,
  RegisterPayload,
  RegisterResponse,
} from "./types/payload";
import { useAuthAPI } from "./api";

export function useAuth() {
  const user = useState<User | null>("auth-user", () => null);
  const isLoading = useState<boolean>("auth-loading", () => false);
  const initialized = useState<boolean>("auth-initialized", () => false);
  const authAPI = useAuthAPI();

  async function me() {
    isLoading.value = true;
    try {
      const currentUser = await authAPI.me();
      user.value = { ...currentUser, name: currentUser.name ?? currentUser.full_name };
      return currentUser;
    } catch {
      user.value = null;
      return null;
    } finally {
      initialized.value = true;
      isLoading.value = false;
    }
  }

  async function login(payload: LoginPayload) {
    isLoading.value = true;
    try {
      const response = await authAPI.login(payload);
      user.value = { ...response.user, name: response.user.name ?? response.user.full_name };
      return response;
    } finally {
      isLoading.value = false;
    }
  }

  // register no longer creates a session — user must verify email first
  async function register(payload: RegisterPayload): Promise<RegisterResponse> {
    isLoading.value = true;
    try {
      return await authAPI.register(payload);
    } finally {
      isLoading.value = false;
    }
  }

  async function setUserFromAuthResponse(response: AuthResponse) {
    user.value = { ...response.user, name: response.user.name ?? response.user.full_name };
  }

  async function logout() {
    isLoading.value = true;
    try {
      await authAPI.logout();
      user.value = null;
    } finally {
      initialized.value = true;
      isLoading.value = false;
    }
  }

  function clearUser() {
    user.value = null;
    initialized.value = true;
  }

  return {
    user,
    isLoading,
    initialized,
    me,
    login,
    register,
    setUserFromAuthResponse,
    logout,
    clearUser,
  };
}
