export interface SaveCompanyPayload {
  name: string;
  description: string;
  register_id: string;
  contact_info: string;
  profile_url?: string;
  city: string;
  district: string;
  location_x?: number;
  location_y?: number;
  location_description?: string;
  benefits: string[];
}
