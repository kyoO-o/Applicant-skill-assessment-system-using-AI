import type { UserRole } from "../constants";
import type { User } from "../usermod";

export interface LoginPayload {
  email: string;
  password: string;
}

export interface RegisterPayload {
  name: string;
  company_name?: string;
  position?: string;
  email: string;
  password: string;
  role: UserRole;
}

export interface AuthResponse {
  user: User;
  message: string;
}

export interface LogoutResponse {
  message: string;
}
