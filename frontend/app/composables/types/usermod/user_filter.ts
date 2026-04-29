export interface UserFilter {
  keyword?: string;
  ids?: number[];
  role?: UserRole;
  email?: string;
  emails?: string[];
}
