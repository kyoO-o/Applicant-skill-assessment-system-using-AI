import type { BaseModel } from "../base";

export interface CompanyBenefit extends BaseModel {
  company_id: number;
  description: string;
}

export interface Company extends BaseModel {
  name: string;
  description: string;
  register_id: string;
  contact_info: string;
  city: string;
  district: string;
  location_x: number;
  location_y: number;
  logo_url?: string;
  benefits: CompanyBenefit[];
}
