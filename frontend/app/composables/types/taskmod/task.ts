import type { BaseModel } from "../base";

export interface TaskSubmission extends BaseModel {
  task_id: number;
  applicant_id: number;
  content: string;
  file_path?: string;
  has_file?: boolean;
  grade: number | null;
  feedback: string;
  status: "submitted" | "graded";
  task_title?: string;
  task_description?: string;
  applicant_name?: string;
  job_posting_id?: number | null;
  job_title?: string;
}

export interface Task extends BaseModel {
  job_posting_id: number | null;
  application_id: number | null;
  title: string;
  description: string;
  due_date: string | null;
  duration_days: number | null;
  status: "draft" | "sent" | "completed" | "graded";
  created_by_ai: boolean;
  job_title?: string;
  submissions?: TaskSubmission[];
}
