import type { UserRole } from "../constants";

export interface SaveUserPayload {
  email: string;
  user_roles: UserRole[];
}
