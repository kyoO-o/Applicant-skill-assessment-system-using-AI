import type { Application, ApplicationFilter, AssessmentResult } from "../../types";

export class ApplicationsAPI {
  baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async fetch<T>(path: string, options: Parameters<typeof $fetch<T>>[1] = {}) {
    const headers = new Headers(options?.headers as HeadersInit | undefined);
    if (import.meta.server) {
      const requestHeaders = useRequestHeaders(["cookie"]);
      if (requestHeaders.cookie && !headers.has("cookie")) {
        headers.set("cookie", requestHeaders.cookie);
      }
    }
    return $fetch<T>(`${this.baseURL}${path}`, { ...options, credentials: "include", headers });
  }

  applyToJob(jobID: number, cvFile: File, assessment?: AssessmentResult | null) {
    const form = new FormData();
    form.append("cv", cvFile);
    if (assessment) {
      form.append("assessment", JSON.stringify(assessment));
    }
    return this.fetch<{ message: string; id: number }>(`/api/jobs/${jobID}/apply`, {
      method: "POST",
      body: form,
    });
  }

  applyFromProfile(jobID: number, assessment?: AssessmentResult | null) {
    return this.fetch<{ message: string; id: number }>(`/api/jobs/${jobID}/apply-from-profile`, {
      method: "POST",
      body: assessment ? { assessment } : {},
    });
  }

  analyzeJob(jobID: number, cvFile: File) {
    const form = new FormData();
    form.append("cv", cvFile);
    return this.fetch<AssessmentResult>(`/api/jobs/${jobID}/analyze`, {
      method: "POST",
      body: form,
    });
  }

  analyzeFromProfile(jobID: number) {
    return this.fetch<AssessmentResult>(`/api/jobs/${jobID}/analyze-from-profile`, {
      method: "POST",
    });
  }

  listMine(filter: ApplicationFilter = {}) {
    return this.fetch<Application[]>("/api/applications", { query: filter });
  }

  get(id: number) {
    return this.fetch<Application>(`/api/applications/${id}`);
  }

  listForJob(jobID: number, filter: ApplicationFilter = {}) {
    return this.fetch<Application[]>(`/api/jobs/${jobID}/applications`, { query: filter });
  }

  updateStatus(id: number, status: string) {
    return this.fetch<{ message: string }>(`/api/applications/${id}/status`, {
      method: "PUT",
      body: { status },
    });
  }

  scheduleInterview(id: number, interviewAt: string, location?: string, note?: string, generateMeet?: boolean) {
    return this.fetch<Application>(`/api/applications/${id}/interview`, {
      method: "PUT",
      body: { interview_at: interviewAt, interview_location: location ?? "", interview_note: note ?? "", generate_meet: generateMeet ?? false },
    });
  }

  listInterviews() {
    return this.fetch<Application[]>("/api/interviews");
  }
}

export const useApplicationsAPI = () => {
  const config = useRuntimeConfig();
  return new ApplicationsAPI(config.public.apiBase);
};
