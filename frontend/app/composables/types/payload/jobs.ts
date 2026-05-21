import type { JobStatus } from "../constants";

export interface SaveJobPayload {
  title: string;
  location: string;
  additional_info: string;
  contact_info: string;
  type: string;
  level: string;
  department: string;
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
}

export interface DeleteJobResponse {
  message: string;
}
