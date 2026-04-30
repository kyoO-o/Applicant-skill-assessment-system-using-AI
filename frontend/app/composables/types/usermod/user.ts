export interface User extends BaseModel {
  first_name?: string;
  last_name?: string;
  full_name: string;
  email: string;
  profile_url?: string;
  name?: string;
  company_id?: number;
  company_name?: string;
  position?: string;
  phone_number: string;
  todu_id: number;
  role: UserRole;
}
