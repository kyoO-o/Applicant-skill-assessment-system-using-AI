<script setup lang="ts">
import { JobStatus } from "../composables/types";
import type { Application, Task } from "../composables/types";
import { toast } from "vue-sonner";
import {
  BriefcaseBusiness,
  FileText,
  Sparkles,
  Calendar,
  ClipboardList,
  Users,
  ArrowRight,
  MapPin,
  Clock,
  Mail,
  NotebookPen,
} from "lucide-vue-next";

const { user } = useAuth();
const config = useRuntimeConfig();
const jobsAPI = useJobsAPI();
const applicationsAPI = useApplicationsAPI();
const tasksAPI = useTasksAPI();
const router = useRouter();

const recruiterJobs = ref<Job[]>([]);
const recruiterInterviews = ref<Application[]>([]);
const myApplications = ref<Application[]>([]);
const myTasks = ref<Task[]>([]);
const loading = ref(false);

function profileImgUrl(path?: string) {
  if (!path) return null;
  return `${config.public.apiBase}${path}`;
}

function formatInterviewDate(iso: string) {
  const d = new Date(iso);
  return d.toLocaleDateString("mn-MN", {
    month: "short",
    day: "numeric",
    weekday: "short",
    hour: "2-digit",
    minute: "2-digit",
  });
}

const recruiterStats = computed(() => {
  const totalJobs = recruiterJobs.value.length;
  const activeJobs = recruiterJobs.value.filter(
    (job) => job.status === JobStatus.Posted,
  ).length;
  const totalApplicants = recruiterJobs.value.reduce(
    (sum, job) => sum + (job.applicants_count ?? 0),
    0,
  );
  const newApplicantsThisWeek = recruiterJobs.value.reduce(
    (sum, job) => sum + (job.new_applicants_this_week ?? 0),
    0,
  );

  const now = new Date();
  const monday = new Date(now);
  monday.setHours(0, 0, 0, 0);
  monday.setDate(now.getDate() - ((now.getDay() + 6) % 7));
  const sunday = new Date(monday);
  sunday.setDate(monday.getDate() + 6);
  sunday.setHours(23, 59, 59, 999);

  const interviewsThisWeek = recruiterInterviews.value.filter((a) => {
    if (!a.interview_at) return false;
    const d = new Date(a.interview_at);
    return d >= monday && d <= sunday;
  }).length;

  const nearestInterview =
    recruiterInterviews.value.find(
      (a) => a.interview_at && new Date(a.interview_at) >= now,
    ) ?? null;

  return {
    totalJobs,
    activeJobs,
    totalApplicants,
    newApplicantsThisWeek,
    interviewsThisWeek,
    nearestInterview,
  };
});

const applicantStats = computed(() => {
  const totalApplications = myApplications.value.length;
  const scored = myApplications.value.filter((a) => a.overall_score > 0);
  const avgScore = scored.length
    ? Math.round(
        scored.reduce((sum, a) => sum + a.overall_score, 0) / scored.length,
      )
    : null;
  const interviews = myApplications.value.filter((a) => a.interview_at).length;
  const pendingTasks = myTasks.value.filter((t) => t.status === "sent").length;
  return { totalApplications, avgScore, interviews, pendingTasks };
});

const recruiterNeedsCompany = computed(
  () => user.value?.role === "recruiter" && !user.value?.company_id,
);

function openRecruiterJobs() {
  if (recruiterNeedsCompany.value) {
    toast.warning(
      "Ажлын байр удирдахаас өмнө компанийн мэдээллээ бүртгэнэ үү.",
    );
    router.push("/profile");
    return;
  }
  router.push("/jobs");
}

async function loadRecruiterDashboard() {
  if (user.value?.role !== "recruiter") return;
  loading.value = true;
  try {
    [recruiterJobs.value, recruiterInterviews.value] = await Promise.all([
      jobsAPI.list(),
      applicationsAPI.listInterviews(),
    ]);
  } finally {
    loading.value = false;
  }
}

async function loadApplicantDashboard() {
  if (user.value?.role === "recruiter") return;
  loading.value = true;
  try {
    [myApplications.value, myTasks.value] = await Promise.all([
      applicationsAPI.listMine(),
      tasksAPI.list(),
    ]);
  } finally {
    loading.value = false;
  }
}

await Promise.all([loadRecruiterDashboard(), loadApplicantDashboard()]);
</script>

<template>
  <!-- ── RECRUITER DASHBOARD ──────────────────────────────── -->
  <div v-if="user?.role === 'recruiter'" class="space-y-5">
    <!-- Welcome -->
    <div class="flex items-end justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-semibold tracking-[-0.6px]">
          Сайн байна уу, {{ user?.name }}
        </h1>
        <p class="mt-1.5 text-[14px] text-muted-foreground">
          Өнөөдрийн ажлын урсгал.
        </p>
      </div>
      <button
        class="inline-flex items-center gap-2 rounded-full px-4 py-2.5 text-[13.5px] font-semibold text-white transition hover:opacity-90"
        style="
          background: linear-gradient(
            135deg,
            var(--primary),
            oklch(0.348 0.106 295)
          );
        "
        @click="openRecruiterJobs"
      >
        <BriefcaseBusiness class="h-4 w-4" />
        Ажлын байр удирдах
      </button>
    </div>

    <!-- Stats row -->
    <div class="grid grid-cols-2 gap-3.5 lg:grid-cols-4">
      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted text-foreground"
          >
            <BriefcaseBusiness class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Нийт ажлын байр</p>
            <p
              class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]"
            >
              {{ loading ? "…" : recruiterStats.totalJobs }}
            </p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-success-foreground">
          {{ recruiterStats.activeJobs }} идэвхтэй
        </p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted text-foreground"
          >
            <Users class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Нийт горилогч</p>
            <p
              class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]"
            >
              {{ loading ? "…" : recruiterStats.totalApplicants }}
            </p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-success-foreground">
          <template v-if="!loading && recruiterStats.newApplicantsThisWeek > 0">
            +{{ recruiterStats.newApplicantsThisWeek }} энэ долоо хоногт
          </template>
          <template v-else-if="!loading">
            Энэ долоо хоногт шинэ горилогч байхгүй
          </template>
        </p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl"
            style="
              background: oklch(0.408 0.124 295 / 8%);
              color: var(--primary);
            "
          >
            <FileText class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Энэ долоо хоногт</p>
            <p
              class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]"
            >
              {{ loading ? "…" : recruiterStats.newApplicantsThisWeek }}
            </p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-primary">
          Шинэ илгээлт
        </p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted text-foreground"
          >
            <Calendar class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Ярилцлага</p>
            <p
              class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]"
            >
              {{ loading ? "…" : recruiterStats.interviewsThisWeek }}
            </p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-muted-foreground">
          Энэ долоо хоногт
        </p>
      </div>
    </div>

    <!-- Jobs list + Activity -->
    <div class="grid gap-4 lg:grid-cols-[1.3fr_0.9fr]">
      <!-- Jobs -->
      <div class="rounded-2xl border border-border bg-card">
        <div
          class="flex items-center justify-between border-b border-border px-5 py-4"
        >
          <div>
            <p class="text-[15px] font-semibold tracking-[-0.2px]">
              Ажлын байрны хөдөлгөөн
            </p>
            <p class="mt-0.5 text-[12px] text-muted-foreground">
              {{ loading ? "Ачааллаж байна…" : "Таны сүүлийн зарууд" }}
            </p>
          </div>
          <button
            class="inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-[12.5px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground"
            @click="openRecruiterJobs"
          >
            Бүгд <ArrowRight class="h-3.5 w-3.5" />
          </button>
        </div>
        <div>
          <div
            v-for="job in recruiterJobs.slice(0, 4)"
            :key="job.id"
            class="flex items-center gap-4 border-b border-border px-5 py-4 last:border-0"
          >
            <div
              class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-xl bg-muted text-sm font-bold text-muted-foreground"
            >
              {{ (job.company_name || job.title || "?")[0]?.toUpperCase() }}
            </div>
            <div class="flex-1 min-w-0">
              <p class="truncate text-[14.5px] font-semibold">
                {{ job.title }}
              </p>
              <p class="mt-0.5 text-[12.5px] text-muted-foreground">
                {{ job.location }} ·
                {{ job.employment_type || "Тогтоогдоогүй" }}
              </p>
            </div>
            <span
              class="inline-flex items-center rounded-full px-2.5 py-1 text-[11.5px] font-semibold capitalize"
              :class="
                job.status === JobStatus.Posted
                  ? 'badge-success'
                  : 'bg-muted text-muted-foreground'
              "
            >
              <span
                class="mr-1.5 h-1.5 w-1.5 rounded-full"
                :class="
                  job.status === JobStatus.Posted
                    ? 'bg-success'
                    : 'bg-muted-foreground'
                "
              />
              {{ job.status }}
            </span>
          </div>
          <div
            v-if="!recruiterJobs.length && !loading"
            class="px-5 py-12 text-center text-[13.5px] text-muted-foreground"
          >
            Эхний ажлын зараа оруулбал энд харагдана.
          </div>
        </div>
      </div>

      <!-- Nearest interview -->
      <div class="rounded-2xl border border-border bg-card">
        <div
          class="flex items-center justify-between border-b border-border px-5 py-4"
        >
          <div>
            <p class="text-[15px] font-semibold tracking-[-0.2px]">
              Ойрын ярилцлага
            </p>
            <p class="mt-0.5 text-[12px] text-muted-foreground">
              Дараагийн товлогдсон уулзалт
            </p>
          </div>
          <button
            class="inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-[12.5px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground"
            @click="router.push('/interviews')"
          >
            Бүгд <ArrowRight class="h-3.5 w-3.5" />
          </button>
        </div>
        <div class="p-4">
          <template v-if="recruiterStats.nearestInterview">
            <div
              class="overflow-hidden rounded-xl border border-border bg-card"
            >
              <!-- Profile header -->
              <div class="flex items-center gap-3.5 border-b border-border p-4">
                <!-- Avatar -->
                <div class="relative flex-shrink-0">
                  <img
                    v-if="
                      profileImgUrl(
                        recruiterStats.nearestInterview.applicant_profile_url,
                      )
                    "
                    :src="
                      profileImgUrl(
                        recruiterStats.nearestInterview.applicant_profile_url,
                      )!
                    "
                    class="h-12 w-12 rounded-full object-cover"
                  />
                  <div
                    v-else
                    class="flex h-12 w-12 items-center justify-center rounded-full text-[17px] font-bold"
                    style="
                      background: oklch(0.408 0.124 295 / 12%);
                      color: var(--primary);
                    "
                  >
                    {{
                      (recruiterStats.nearestInterview.applicant_name ||
                        "?")[0]?.toUpperCase()
                    }}
                  </div>
                </div>
                <!-- Name + email + job -->
                <div class="flex-1 min-w-0">
                  <p class="truncate text-[15px] font-semibold leading-tight">
                    {{
                      recruiterStats.nearestInterview.applicant_name ||
                      "Горилогч"
                    }}
                  </p>
                  <div
                    class="mt-0.5 flex items-center gap-1.5 text-[12px] text-muted-foreground"
                  >
                    <Mail class="h-3 w-3 flex-shrink-0" />
                    <span class="truncate">{{
                      recruiterStats.nearestInterview.applicant_email || "—"
                    }}</span>
                  </div>
                  <div class="mt-1.5 flex items-center gap-2 flex-wrap">
                    <span
                      class="inline-flex items-center rounded-full bg-muted px-2 py-0.5 text-[11px] font-medium text-muted-foreground"
                    >
                      {{ recruiterStats.nearestInterview.job_title }}
                    </span>
                    <span
                      v-if="recruiterStats.nearestInterview.overall_score > 0"
                      class="inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-semibold"
                      :style="
                        recruiterStats.nearestInterview.overall_score >= 70
                          ? 'background: oklch(0.78 0.15 145 / 15%); color: oklch(0.45 0.15 145)'
                          : recruiterStats.nearestInterview.overall_score >= 40
                            ? 'background: oklch(0.85 0.16 75 / 15%); color: oklch(0.5 0.16 75)'
                            : 'background: oklch(0.7 0.18 25 / 15%); color: oklch(0.45 0.18 25)'
                      "
                    >
                      Үнэлгээ
                      {{ recruiterStats.nearestInterview.overall_score }}%
                    </span>
                  </div>
                </div>
              </div>
              <!-- Interview details -->
              <div class="space-y-2.5 p-4">
                <div class="flex items-center gap-2.5 text-[13px]">
                  <Clock class="h-3.5 w-3.5 flex-shrink-0 text-primary" />
                  <span class="font-medium">{{
                    formatInterviewDate(
                      recruiterStats.nearestInterview.interview_at!,
                    )
                  }}</span>
                </div>
                <div
                  v-if="recruiterStats.nearestInterview.interview_location"
                  class="flex items-center gap-2.5 text-[13px] text-muted-foreground"
                >
                  <MapPin class="h-3.5 w-3.5 flex-shrink-0" />
                  <span class="truncate">{{
                    recruiterStats.nearestInterview.interview_location
                  }}</span>
                </div>
                <div
                  v-if="recruiterStats.nearestInterview.interview_note"
                  class="flex items-start gap-2.5 text-[13px] text-muted-foreground"
                >
                  <NotebookPen class="mt-0.5 h-3.5 w-3.5 flex-shrink-0" />
                  <span class="line-clamp-2">{{
                    recruiterStats.nearestInterview.interview_note
                  }}</span>
                </div>
              </div>
            </div>
          </template>
          <div
            v-else
            class="px-4 py-10 text-center text-[13px] text-muted-foreground"
          >
            Товлогдсон ярилцлага байхгүй байна.
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- ── APPLICANT DASHBOARD ──────────────────────────────── -->
  <div v-else class="space-y-5">
    <!-- Welcome -->
    <div class="flex items-end justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-semibold tracking-[-0.6px]">
          Тавтай морилно уу, {{ user?.name }}
        </h1>
        <p class="mt-1.5 text-[14px] text-muted-foreground">
          Энэ долоо хоногт таны хайлтын байдал.
        </p>
      </div>
      <button
        class="inline-flex items-center gap-2 rounded-full px-4 py-2.5 text-[13.5px] font-semibold text-white transition hover:opacity-90"
        style="
          background: linear-gradient(
            135deg,
            var(--primary),
            oklch(0.348 0.106 295)
          );
        "
        @click="router.push('/chat')"
      >
        <Sparkles class="h-4 w-4" />
        AI туслагчаас асуух
      </button>
    </div>

    <!-- Stats row -->
    <div class="grid grid-cols-2 gap-3.5 lg:grid-cols-4">
      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted"
          >
            <BriefcaseBusiness class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Идэвхтэй анкет</p>
            <p
              class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]"
            >
              {{ loading ? "…" : applicantStats.totalApplications }}
            </p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-success-foreground">
          Илгээсэн
        </p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl"
            style="
              background: oklch(0.408 0.124 295 / 8%);
              color: var(--primary);
            "
          >
            <Sparkles class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Дундаж AI оноо</p>
            <p
              class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]"
            >
              {{
                loading
                  ? "…"
                  : applicantStats.avgScore !== null
                    ? applicantStats.avgScore
                    : "—"
              }}
            </p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-primary">AI үнэлгээ</p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted"
          >
            <Calendar class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Ярилцлага</p>
            <p
              class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]"
            >
              {{ loading ? "…" : applicantStats.interviews }}
            </p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-muted-foreground">
          Товлогдсон
        </p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted"
          >
            <ClipboardList class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Даалгавар</p>
            <p
              class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]"
            >
              {{ loading ? "…" : applicantStats.pendingTasks }}
            </p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-warning-foreground">
          Хүлээгдэж байгаа
        </p>
      </div>
    </div>

    <!-- Feature cards + AI suggestions -->
    <div class="grid gap-4 lg:grid-cols-3">
      <div
        class="flex flex-col gap-3 rounded-2xl border border-border bg-card p-5 cursor-pointer transition hover:border-primary/30"
        @click="router.push('/jobs')"
      >
        <div
          class="flex h-10 w-10 items-center justify-center rounded-xl bg-muted"
        >
          <BriefcaseBusiness class="h-5 w-5 text-muted-foreground" />
        </div>
        <div>
          <p class="text-[15px] font-semibold">Тохирох ажлыг олоорой</p>
          <p class="mt-1.5 text-[13px] leading-[1.5] text-muted-foreground">
            Ажил олгогчдийн зарласан нээлттэй байр дундаас өөрт тохирохыг нь
            олоорой.
          </p>
        </div>
        <div
          class="mt-auto flex items-center gap-1 text-[13px] font-medium text-primary"
        >
          Ажил хайх <ArrowRight class="h-3.5 w-3.5" />
        </div>
      </div>

      <div
        class="flex flex-col gap-3 rounded-2xl border border-border bg-card p-5 cursor-pointer transition hover:border-primary/30"
        @click="router.push('/applications')"
      >
        <div
          class="flex h-10 w-10 items-center justify-center rounded-xl bg-muted"
        >
          <FileText class="h-5 w-5 text-muted-foreground" />
        </div>
        <div>
          <p class="text-[15px] font-semibold">Бэлэн байдлаа хянаарай</p>
          <p class="mt-1.5 text-[13px] leading-[1.5] text-muted-foreground">
            Үнэлгээ, профайл мэдээлэл болон ажлын нээлтийг нэг дор хадгалаарай.
          </p>
        </div>
        <div
          class="mt-auto flex items-center gap-1 text-[13px] font-medium text-primary"
        >
          Анкет харах <ArrowRight class="h-3.5 w-3.5" />
        </div>
      </div>

      <!-- AI suggestions card -->
      <div
        class="flex flex-col gap-3 rounded-2xl border p-5 cursor-pointer transition"
        style="
          background: oklch(0.408 0.124 295 / 4.5%);
          border-color: oklch(0.408 0.124 295 / 13%);
        "
        @click="router.push('/chat')"
      >
        <div class="flex items-center gap-3">
          <div
            class="flex h-10 w-10 items-center justify-center rounded-xl"
            style="
              background: linear-gradient(
                135deg,
                var(--primary),
                oklch(0.348 0.106 295)
              );
              color: #fff;
            "
          >
            <Sparkles class="h-5 w-5" />
          </div>
          <div>
            <p class="text-[15px] font-semibold">AI туслагч</p>
            <p class="text-[12px] text-muted-foreground">
              Таны CV-д тулгуурласан
            </p>
          </div>
        </div>
        <p class="text-[13px] leading-[1.5] text-muted-foreground">
          Ярилцлагад бэлэн байгаарай. Анкетаас эхлээд дүгнэлт хүртэлх бүх шатыг
          энэ платформоор хянаарай.
        </p>
        <div
          class="mt-auto flex items-center gap-1 text-[13px] font-medium text-primary"
        >
          Асуух <ArrowRight class="h-3.5 w-3.5" />
        </div>
      </div>
    </div>
  </div>
</template>
