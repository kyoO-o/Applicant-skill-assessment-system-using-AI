import type { Task, TaskSubmission } from "../../types";

export class TasksAPI {
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

  list(jobID?: number) {
    const query = jobID ? `?job_id=${jobID}` : "";
    return this.fetch<Task[]>(`/api/tasks${query}`);
  }

  get(id: number) {
    return this.fetch<Task>(`/api/tasks/${id}`);
  }

  create(payload: { job_posting_id: number; application_id?: number; title: string; description: string; due_date?: string }) {
    return this.fetch<Task>("/api/tasks", { method: "POST", body: payload });
  }

  generate(jobPostingID: number) {
    return this.fetch<{ title: string; description: string }>("/api/tasks/generate", {
      method: "POST",
      body: { job_posting_id: jobPostingID },
    });
  }

  send(id: number) {
    return this.fetch<Task>(`/api/tasks/${id}/send`, { method: "PUT" });
  }

  submit(id: number, content: string) {
    return this.fetch<TaskSubmission>(`/api/tasks/${id}/submit`, {
      method: "POST",
      body: { content },
    });
  }

  grade(submissionID: number, grade: number, feedback: string) {
    return this.fetch<TaskSubmission>(`/api/tasks/submissions/${submissionID}/grade`, {
      method: "PUT",
      body: { grade, feedback },
    });
  }
}

export const useTasksAPI = () => {
  const config = useRuntimeConfig();
  return new TasksAPI(config.public.apiBase);
};
