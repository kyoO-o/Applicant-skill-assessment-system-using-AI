import type { BaseModel } from "../base";
import type { JobStatus } from "../constants";

export interface Job extends BaseModel {
  recruiter_id: number;
  company_id: number;
  company_name?: string;
  company_logo_url?: string;
  company_profile_url?: string;
  title: string;
  location: string;
  additional_info: string;
  description: string;
  contact_info: string;
  type: string;
  employment_type: string;
  level: string;
  seniority: string;
  city?: string;
  district?: string;
  location_x?: number;
  location_y?: number;
  min_salary: number;
  max_salary: number;
  status: JobStatus;
  duties: string[];
  requirements: string[];
  skills: string[];
  bonuses: string[];
  applicants_count: number;
}
