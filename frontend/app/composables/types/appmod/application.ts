import type { BaseModel } from "../base";

export interface SkillResult {
  skill: string;
  explanation: string;
}

export interface Application extends BaseModel {
  job_posting_id: number;
  job_title: string;
  applicant_id: number;
  applicant_name: string;
  applicant_email: string;
  overall_score: number;
  summary: string;
  matched_skills: SkillResult[];
  missing_skills: SkillResult[];
  recommendations: string[];
  status: "pending" | "assessed" | "shortlisted" | "rejected";
  assessed_at: string | null;
}
