export interface User extends BaseModel {
  email: string;
  profile_url: string;
  name: string;
  phone_number: string;
  todu_id: number;
  role: UserRole;
}
