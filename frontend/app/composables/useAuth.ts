import type { User } from "./types";
import type {
  AuthResponse,
  LoginPayload,
  LogoutResponse,
  RegisterPayload,
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
      user.value = {
        ...currentUser,
        name: currentUser.name ?? currentUser.full_name,
      };
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
      user.value = {
        ...response.user,
        name: response.user.name ?? response.user.full_name,
      };
      return response;
    } finally {
      isLoading.value = false;
    }
  }

  async function register(payload: RegisterPayload) {
    isLoading.value = true;

    try {
      const response = await authAPI.register(payload);
      user.value = {
        ...response.user,
        name: response.user.name ?? response.user.full_name,
      };
      return response;
    } finally {
      isLoading.value = false;
    }
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
    logout,
    clearUser,
  };
}
