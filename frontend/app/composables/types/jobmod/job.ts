import type { BaseModel } from "../base";
import type { JobStatus } from "../constants";

export interface Job extends BaseModel {
  recruiter_id: number;
  company_id: number;
  company_name?: string;
  title: string;
  location: string;
  employment_type: string;
  seniority: string;
  status: JobStatus;
  description: string;
  requirements: string[];
  applicants_count: number;
}
