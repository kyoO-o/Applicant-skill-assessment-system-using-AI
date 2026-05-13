import type { Application, AssessmentResult } from "../../types";

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

  applyToJob(jobID: number, cvFile: File) {
    const form = new FormData();
    form.append("cv", cvFile);
    return this.fetch<{ message: string; id: number }>(`/api/jobs/${jobID}/apply`, {
      method: "POST",
      body: form,
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

  listMine() {
    return this.fetch<Application[]>("/api/applications");
  }

  get(id: number) {
    return this.fetch<Application>(`/api/applications/${id}`);
  }

  listForJob(jobID: number) {
    return this.fetch<Application[]>(`/api/jobs/${jobID}/applications`);
  }

  updateStatus(id: number, status: string) {
    return this.fetch<{ message: string }>(`/api/applications/${id}/status`, {
      method: "PUT",
      body: { status },
    });
  }

  scheduleInterview(id: number, interviewAt: string) {
    return this.fetch<Application>(`/api/applications/${id}/interview`, {
      method: "PUT",
      body: { interview_at: interviewAt },
    });
  }
}

export const useApplicationsAPI = () => {
  const config = useRuntimeConfig();
  return new ApplicationsAPI(config.public.apiBase);
};
