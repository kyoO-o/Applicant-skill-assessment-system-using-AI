import type { BaseModel } from "../base";

export interface TaskSubmission extends BaseModel {
  task_id: number;
  applicant_id: number;
  content: string;
  grade: number | null;
  feedback: string;
  status: "submitted" | "graded";
}

export interface Task extends BaseModel {
  job_posting_id: number;
  application_id: number | null;
  title: string;
  description: string;
  due_date: string | null;
  status: "draft" | "sent" | "completed" | "graded";
  created_by_ai: boolean;
  submissions?: TaskSubmission[];
}
