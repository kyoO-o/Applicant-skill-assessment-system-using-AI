import type { Task, TaskFilter, TaskSubmission } from "../../types";

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

  list(jobID?: number, filter: TaskFilter = {}) {
    return this.fetch<Task[]>("/api/tasks", { query: { ...(jobID ? { job_id: jobID } : {}), ...filter } });
  }

  get(id: number) {
    return this.fetch<Task>(`/api/tasks/${id}`);
  }

  create(payload: { job_posting_id?: number; application_id?: number; title: string; description: string; due_date?: string; duration_days?: number }) {
    return this.fetch<Task>("/api/tasks", { method: "POST", body: payload });
  }

  update(id: number, payload: { title: string; description: string; due_date?: string; duration_days?: number }) {
    return this.fetch<Task>(`/api/tasks/${id}`, { method: "PUT", body: payload });
  }

  generate(jobPostingID: number) {
    return this.fetch<{ title: string; description: string }>("/api/tasks/generate", {
      method: "POST",
      body: { job_posting_id: jobPostingID },
    });
  }

  send(id: number, applicationID?: number, dueDate?: string) {
    return this.fetch<Task>(`/api/tasks/${id}/send`, {
      method: "PUT",
      body: { application_id: applicationID, due_date: dueDate },
    });
  }

  submit(id: number, content: string, file?: File) {
    const form = new FormData();
    if (content) form.append("content", content);
    if (file) form.append("file", file);
    return this.fetch<TaskSubmission>(`/api/tasks/${id}/submit`, {
      method: "POST",
      body: form,
    });
  }

  delete(id: number) {
    return this.fetch<void>(`/api/tasks/${id}`, { method: "DELETE" });
  }

  grade(submissionID: number, grade: number, feedback: string) {
    return this.fetch<TaskSubmission>(`/api/tasks/submissions/${submissionID}/grade`, {
      method: "PUT",
      body: { grade, feedback },
    });
  }

  aiGrade(submissionID: number) {
    return this.fetch<TaskSubmission>(`/api/tasks/submissions/${submissionID}/ai-grade`, {
      method: "POST",
    });
  }

  listAllSubmissions() {
    return this.fetch<TaskSubmission[]>("/api/tasks/submissions");
  }

  submissionFileURL(submissionID: number) {
    return `${this.baseURL}/api/tasks/submissions/${submissionID}/file`;
  }
}

export const useTasksAPI = () => {
  const config = useRuntimeConfig();
  return new TasksAPI(config.public.apiBase);
};
