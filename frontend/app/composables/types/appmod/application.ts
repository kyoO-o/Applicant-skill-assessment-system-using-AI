import type { BaseModel } from "../base";

export interface SkillResult {
  skill: string;
  explanation: string;
}

export interface DutyAssessment {
  duty: string;
  status: "met" | "partial" | "not_met";
  explanation: string;
}

export interface RequirementAssessment {
  requirement: string;
  status: "met" | "partial" | "not_met";
  explanation: string;
}

export interface AssessmentResult {
  overall_score: number;
  summary: string;
  matched_skills: SkillResult[];
  missing_skills: SkillResult[];
  recommendations: string[];
  duty_assessments: DutyAssessment[];
  requirement_assessments: RequirementAssessment[];
}

export interface Application extends BaseModel {
  job_posting_id: number;
  job_title: string;
  applicant_id: number;
  applicant_name: string;
  applicant_email: string;
  applicant_profile_url?: string;
  overall_score: number;
  summary: string;
  matched_skills: SkillResult[];
  missing_skills: SkillResult[];
  recommendations: string[];
  duty_assessments: DutyAssessment[];
  requirement_assessments: RequirementAssessment[];
  status: "pending" | "assessed" | "shortlisted" | "rejected";
  assessed_at: string | null;
  interview_at: string | null;
  interview_location: string | null;
  interview_note: string | null;
}
