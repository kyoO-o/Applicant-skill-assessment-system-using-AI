import type { JobStatus } from "../constants";

export interface SaveJobPayload {
  title: string;
  location: string;
  employment_type: string;
  seniority: string;
  status: JobStatus;
  description: string;
  requirements: string[];
}

export interface DeleteJobResponse {
  message: string;
}
