<script setup lang="ts">
import type { Job } from "../../composables/types";
import { JobStatus } from "../../composables/types";
import { toast } from "vue-sonner";
import {
  BriefcaseBusiness,
  Eye,
  Plus,
  Trash2,
  Users,
  MapPin,
  Loader2,
  LayoutGrid,
  LayoutList,
  Search,
  Heart,
  SlidersHorizontal,
  ChevronLeft,
  ChevronRight,
} from "lucide-vue-next";

const { user } = useAuth();
const jobsAPI = useJobsAPI();
const {
  public: { apiBase },
} = useRuntimeConfig();

function logoURL(url: string | undefined): string | undefined {
  if (!url) return undefined;
  if (url.startsWith("http")) return url;
  return `${apiBase}${url}`;
}
const router = useRouter();

const jobs = ref<Job[]>([]);
const selectedJob = ref<Job | null>(null);
const loading = ref(false);
const deleteDialogOpen = ref(false);
const isDeleting = ref(false);

const isRecruiter = computed(() => user.value?.role === "recruiter");
const recruiterCompanyID = computed(() => user.value?.company_id ?? 0);
const hasRecruiterCompany = computed(() => recruiterCompanyID.value > 0);

const recPage = ref(1);
const REC_PAGE_SIZE = 10;
const recTotalPages = computed(() =>
  Math.max(1, Math.ceil(jobs.value.length / REC_PAGE_SIZE)),
);
const paginatedRecruiterJobs = computed(() => {
  const start = (recPage.value - 1) * REC_PAGE_SIZE;
  return jobs.value.slice(start, start + REC_PAGE_SIZE);
});

const recruiterStats = computed(() => {
  const totalJobs = jobs.value.length;
  const activeJobs = jobs.value.filter(
    (job) => job.status === JobStatus.Posted,
  ).length;
  const totalApplicants = jobs.value.reduce(
    (sum, job) => sum + job.applicants_count,
    0,
  );
  return { totalJobs, activeJobs, totalApplicants };
});

const applicantJobs = computed(() =>
  jobs.value.filter((job) => job.status === JobStatus.Posted),
);

function statusLabel(status: Job["status"]) {
  if (status === JobStatus.Posted) return "Нийтлэгдсэн";
  if (status === JobStatus.Closed) return "Хаагдсан";
  return "Ноорог";
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("mn-MN", { dateStyle: "medium" }).format(
    new Date(value),
  );
}

function openCreatePage() {
  if (!hasRecruiterCompany.value) {
    toast.warning("Ажлын байр нэмэхийн өмнө компанийн мэдээллээ бүртгэнэ үү.");
    return;
  }
  router.push("/jobs/new");
}

function openDeleteDialog(job: Job) {
  selectedJob.value = job;
  deleteDialogOpen.value = true;
}

async function loadJobs() {
  loading.value = true;
  try {
    if (isRecruiter.value && hasRecruiterCompany.value) {
      jobs.value = await jobsAPI.listByCompany(recruiterCompanyID.value);
      return;
    }
    jobs.value = await jobsAPI.list();
  } finally {
    loading.value = false;
  }
}

async function confirmDelete() {
  if (!selectedJob.value) return;

  isDeleting.value = true;
  try {
    await jobsAPI.delete(recruiterCompanyID.value, selectedJob.value.id);
    deleteDialogOpen.value = false;
    selectedJob.value = null;
    await loadJobs();
  } catch (error: any) {
    toast.error(
      error?.data?.message ||
        error?.message ||
        "Ажлын байр устгахад алдаа гарлаа.",
    );
  } finally {
    isDeleting.value = false;
  }
}

await loadJobs();

// ── Applicant browse state ─────────────────────────────────
const query = ref("");
const FILTER_OPTIONS = [
  { value: "all", name: "бүгд" },
  { value: "saved", name: "хадгалсан" },
] as const;
type FilterOption = (typeof FILTER_OPTIONS)[number]["value"];
const filter = ref<FilterOption>("all");
const view = ref<"grid" | "list">("grid");
const page = ref(1);
const PAGE_SIZE = 9;

const savedJobIds = ref<Set<number>>(
  new Set(
    import.meta.client
      ? JSON.parse(localStorage.getItem("saved-jobs") || "[]")
      : [],
  ),
);

function toggleSave(job: Job, e: Event) {
  e.stopPropagation();
  if (savedJobIds.value.has(job.id)) {
    savedJobIds.value.delete(job.id);
  } else {
    savedJobIds.value.add(job.id);
  }
  savedJobIds.value = new Set(savedJobIds.value);
  if (import.meta.client) {
    localStorage.setItem("saved-jobs", JSON.stringify([...savedJobIds.value]));
  }
}

const filteredJobs = computed(() => {
  const q = query.value.toLowerCase().trim();
  return applicantJobs.value.filter((job) => {
    if (filter.value === "saved" && !savedJobIds.value.has(job.id))
      return false;
    // if (filter.value === "remote") {
    //   const isRemote =
    //     job.type?.toLowerCase().includes("remote") ||
    //     job.employment_type?.toLowerCase().includes("remote");
    //   if (!isRemote) return false;
    // }
    if (!q) return true;
    return (
      job.title.toLowerCase().includes(q) ||
      (job.company_name || "").toLowerCase().includes(q) ||
      job.location.toLowerCase().includes(q)
    );
  });
});

const totalPages = computed(() =>
  Math.max(1, Math.ceil(filteredJobs.value.length / PAGE_SIZE)),
);

const paginatedJobs = computed(() => {
  const start = (page.value - 1) * PAGE_SIZE;
  return filteredJobs.value.slice(start, start + PAGE_SIZE);
});

watch([query, filter], () => {
  page.value = 1;
});

const GRADIENTS = [
  "linear-gradient(135deg, #fde68a, #fb923c)",
  "linear-gradient(135deg, #bae6fd, #818cf8)",
  "linear-gradient(135deg, #fecaca, #f472b6)",
  "linear-gradient(135deg, #fdba74, #fb7185)",
  "linear-gradient(135deg, #a7f3d0, #67e8f9)",
  "linear-gradient(135deg, #ddd6fe, #c4b5fd)",
];

function companyGradient(name: string): string {
  let h = 0;
  for (const c of name) h = (h * 31 + c.charCodeAt(0)) & 0xffffffff;
  return GRADIENTS[Math.abs(h) % GRADIENTS.length];
}

function formatRelativeDate(value: string): string {
  const d = Math.floor((Date.now() - new Date(value).getTime()) / 86400000);
  if (d === 0) return "Today";
  if (d === 1) return "1d ago";
  if (d < 7) return `${d}d ago`;
  if (d < 30) return `${Math.floor(d / 7)}w ago`;
  return `${Math.floor(d / 30)}mo ago`;
}

function formatSalary(job: Job): string {
  const fmt = (n: number) =>
    n >= 1_000_000
      ? `${(n / 1_000_000).toFixed(n % 1_000_000 ? 1 : 0)}M`
      : n >= 1_000
        ? `${Math.round(n / 1_000)}K`
        : String(n);
  if (job.min_salary && job.max_salary)
    return `${fmt(job.min_salary)}–${fmt(job.max_salary)}₮`;
  if (job.min_salary) return `${fmt(job.min_salary)}₮+`;
  if (job.max_salary) return `up to ${fmt(job.max_salary)}₮`;
  return "";
}
</script>

<template>
  <!-- ── RECRUITER VIEW ──────────────────────────────────── -->
  <div v-if="isRecruiter" class="space-y-5">
    <!-- Page header -->
    <div class="flex items-center justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-semibold tracking-[-0.6px]">
          Ажлын байрны удирдлага
        </h1>
        <p class="mt-1 text-[13.5px] text-muted-foreground">
          Ажлын зар, анкетуудыг нэг дор хянаарай.
        </p>
      </div>
      <button
        class="inline-flex shrink-0 items-center gap-2 rounded-full px-4 py-2.5 text-[13.5px] font-semibold text-white transition hover:opacity-90"
        style="
          background: linear-gradient(
            135deg,
            var(--primary),
            oklch(0.348 0.106 295)
          );
        "
        @click="openCreatePage"
      >
        <Plus class="h-4 w-4" />
        Ажлын байр нэмэх
      </button>
    </div>

    <!-- No company banner -->
    <div
      v-if="!hasRecruiterCompany"
      class="rounded-2xl border border-warning/40 bg-warning-bg px-5 py-4 text-[13.5px] text-warning-foreground"
    >
      Ажлын байр нэмэх эсвэл удирдахаас өмнө профайл хэсгээс компанийн мэдээллээ
      бүртгэнэ үү.
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-3 gap-3.5">
      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted"
          >
            <BriefcaseBusiness class="h-4 w-4 text-muted-foreground" />
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
            class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted"
          >
            <Users class="h-4 w-4 text-muted-foreground" />
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
        <p class="mt-2.5 text-[11.5px] font-medium text-muted-foreground">
          Нийт илгээлт
        </p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted"
          >
            <Eye class="h-4 w-4 text-muted-foreground" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Идэвхтэй зар</p>
            <p
              class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]"
            >
              {{ loading ? "…" : recruiterStats.activeJobs }}
            </p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-primary">Нийтлэгдсэн</p>
      </div>
    </div>

    <!-- Jobs list -->
    <div class="rounded-2xl border border-border bg-card">
      <div
        class="flex items-center justify-between border-b border-border px-5 py-4"
      >
        <div>
          <p class="text-[15px] font-semibold tracking-[-0.2px]">
            Таны ажлын байрууд
          </p>
          <p class="mt-0.5 text-[12px] text-muted-foreground">
            {{ loading ? "Ачааллаж байна…" : `${jobs.length} ажлын байр` }}
          </p>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-10">
        <Loader2 class="h-6 w-6 animate-spin text-muted-foreground" />
      </div>

      <!-- Empty -->
      <div
        v-else-if="!jobs.length"
        class="flex flex-col items-center py-16 text-center"
      >
        <div
          class="mb-3 flex h-12 w-12 items-center justify-center rounded-xl bg-muted"
        >
          <BriefcaseBusiness class="h-5 w-5 text-muted-foreground" />
        </div>
        <p class="text-[15px] font-semibold">Ажлын байр байхгүй</p>
        <p class="mt-1.5 text-[13px] text-muted-foreground">
          Эхний ажлын байраа нэмнэ үү.
        </p>
        <button
          class="mt-5 inline-flex items-center gap-2 rounded-full px-4 py-2 text-[13px] font-semibold text-white transition hover:opacity-90"
          style="
            background: linear-gradient(
              135deg,
              var(--primary),
              oklch(0.348 0.106 295)
            );
          "
          @click="openCreatePage"
        >
          <Plus class="h-3.5 w-3.5" />
          Ажлын байр нэмэх
        </button>
      </div>

      <!-- Job rows -->
      <div v-else>
        <div
          v-for="job in paginatedRecruiterJobs"
          :key="job.id"
          class="flex items-center gap-4 border-b border-border px-5 py-4 last:border-0 cursor-pointer transition hover:bg-muted/30"
          @click="router.push(`/jobs/${job.id}/applicants`)"
        >
          <!-- Logo placeholder -->
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-muted text-sm font-bold text-muted-foreground"
          >
            {{ (job.company_name || job.title || "?")[0]?.toUpperCase() }}
          </div>

          <!-- Job info -->
          <div class="flex-1 min-w-0">
            <p class="truncate text-[14.5px] font-semibold">{{ job.title }}</p>
            <p
              class="mt-0.5 flex items-center gap-1.5 text-[12px] text-muted-foreground"
            >
              <MapPin class="h-3 w-3 shrink-0" />
              {{ job.location }} ·
              {{ job.employment_type || job.type || "Тогтоогдоогүй" }}
            </p>
          </div>

          <!-- Status badge -->
          <span
            class="inline-flex shrink-0 items-center rounded-full px-2.5 py-1 text-[11.5px] font-semibold"
            :class="
              job.status === JobStatus.Posted
                ? 'badge-success'
                : job.status === JobStatus.Closed
                  ? 'bg-muted text-muted-foreground'
                  : 'badge-warning'
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
            {{ statusLabel(job.status) }}
          </span>

          <!-- Applicant count -->
          <div
            class="hidden shrink-0 items-center gap-1 text-[13px] text-muted-foreground sm:flex"
          >
            <Users class="h-3.5 w-3.5" />
            {{ job.applicants_count }}
          </div>

          <!-- Date -->
          <p
            class="hidden shrink-0 text-[12.5px] text-muted-foreground lg:block"
          >
            {{ formatDate(job.created_at) }}
          </p>

          <!-- Actions -->
          <div class="flex shrink-0 items-center gap-1.5" @click.stop>
            <button
              class="flex h-8 w-8 items-center justify-center rounded-lg border border-border text-muted-foreground transition hover:bg-muted hover:text-foreground"
              title="Горилогчид"
              @click="router.push(`/jobs/${job.id}/applicants`)"
            >
              <Users class="h-3.5 w-3.5" />
            </button>
            <button
              class="flex h-8 w-8 items-center justify-center rounded-lg border border-border text-muted-foreground transition hover:bg-muted hover:text-foreground"
              title="Харах / Засах"
              @click="router.push(`/jobs/${job.id}`)"
            >
              <Eye class="h-3.5 w-3.5" />
            </button>
            <button
              class="flex h-8 w-8 items-center justify-center rounded-lg border border-border text-muted-foreground transition hover:bg-destructive/10 hover:text-destructive hover:border-destructive/30"
              title="Устгах"
              @click="openDeleteDialog(job)"
            >
              <Trash2 class="h-3.5 w-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>
    <!-- Recruiter jobs pagination -->
    <div class="flex items-center justify-center gap-1 px-5 border-border">
      <button
        class="flex h-8 w-8 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
        :disabled="recPage === 1"
        @click="recPage--"
      >
        <ChevronLeft class="h-4 w-4" />
      </button>
      <button
        v-for="p in recTotalPages"
        :key="p"
        :class="[
          'h-8 w-8 rounded-full text-[13px] font-medium transition',
          p === recPage
            ? 'bg-primary text-background'
            : 'text-muted-foreground hover:bg-muted',
        ]"
        @click="recPage = p"
      >
        {{ p }}
      </button>
      <button
        class="flex h-8 w-8 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
        :disabled="recPage === recTotalPages"
        @click="recPage++"
      >
        <ChevronRight class="h-4 w-4" />
      </button>
    </div>
  </div>

  <!-- ── APPLICANT VIEW ──────────────────────────────────── -->
  <div v-else class="space-y-5">
    <!-- Page header -->
    <div class="flex items-center justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-semibold tracking-[-0.6px]">
          Ажлын байрууд
        </h1>
        <p class="mt-1 text-[13.5px] text-muted-foreground">
          {{
            loading ? "Loading…" : `${filteredJobs.length} ажлын зар олдлоо.`
          }}
        </p>
      </div>
      <!-- Grid / List toggle -->
      <div
        class="flex shrink-0 items-center gap-2.5 rounded-full border border-border bg-card p-1"
      >
        <button
          :class="[
            'flex h-8 w-8 items-center justify-center rounded-full transition',
            view === 'grid'
              ? 'bg-foreground text-background'
              : 'text-muted-foreground hover:text-foreground',
          ]"
          title="Grid view"
          @click="view = 'grid'"
        >
          <LayoutGrid class="h-3.5 w-3.5" />
        </button>
        <button
          :class="[
            'flex h-8 w-8 items-center justify-center rounded-full transition',
            view === 'list'
              ? 'bg-foreground text-background'
              : 'text-muted-foreground hover:text-foreground',
          ]"
          title="List view"
          @click="view = 'list'"
        >
          <LayoutList class="h-3.5 w-3.5" />
        </button>
      </div>
    </div>

    <!-- Filter bar -->
    <div class="flex flex-wrap items-center gap-2.5">
      <!-- Search input -->
      <div class="relative min-w-[260px] flex-1">
        <Search
          class="pointer-events-none absolute left-3.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground"
        />
        <input
          v-model="query"
          class="h-10 w-full rounded-full border border-border bg-card pl-9 pr-4 text-[13.5px] text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/30"
          placeholder="Хайх..."
          type="text"
        />
      </div>
      <!-- Filter pills -->
      <button
        v-for="f in FILTER_OPTIONS"
        :key="f.value"
        :class="[
          'rounded-full border px-4 py-2 text-[13px] font-medium capitalize transition',
          filter === f.value
            ? 'border-foreground bg-foreground text-background'
            : 'border-border bg-card text-foreground hover:bg-muted',
        ]"
        @click="filter = f.value"
      >
        {{ f.name }}
      </button>
      <!-- Filters button (cosmetic) -->
      <button
        class="flex items-center gap-1.5 rounded-full border border-border bg-card px-4 py-2 text-[13px] font-medium text-foreground transition hover:bg-muted"
      >
        <SlidersHorizontal class="h-3.5 w-3.5" />
        Шүүлтүүр
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="h-7 w-7 animate-spin text-primary" />
    </div>

    <!-- Empty (no jobs at all) -->
    <div
      v-else-if="!applicantJobs.length"
      class="flex flex-col items-center py-20 text-center"
    >
      <div
        class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-muted"
      >
        <BriefcaseBusiness class="h-7 w-7 text-muted-foreground" />
      </div>
      <p class="text-[16px] font-semibold">No open positions right now</p>
      <p class="mt-2 text-[13.5px] text-muted-foreground">Check back later.</p>
    </div>

    <!-- No results for current filter -->
    <div
      v-else-if="!filteredJobs.length"
      class="flex flex-col items-center py-16 text-center"
    >
      <div
        class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-muted"
      >
        <Search class="h-5 w-5 text-muted-foreground" />
      </div>
      <p class="text-[15px] font-semibold">No roles match your search</p>
      <p class="mt-1.5 text-[13px] text-muted-foreground">
        Try a different keyword or filter.
      </p>
    </div>

    <template v-else>
      <!-- ─── GRID VIEW ───────────────────────────────────── -->
      <div
        v-if="view === 'grid'"
        class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3"
      >
        <div
          v-for="job in paginatedJobs"
          :key="job.id"
          class="group flex cursor-pointer flex-col gap-3.5 rounded-2xl border border-border bg-card p-5 transition-all hover:border-primary/30 hover:shadow-sm"
          @click="router.push(`/jobs/${job.id}`)"
        >
          <!-- Company mark + title + save -->
          <div class="flex items-start gap-3">
            <img
              v-if="job.company_logo_url"
              :src="logoURL(job.company_logo_url)"
              :alt="job.company_name"
              class="h-11 w-11 shrink-0 rounded-xl object-cover"
            />
            <div
              v-else
              class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl text-[15px] font-bold text-white/90"
              :style="{
                background: companyGradient(job.company_name || job.title),
              }"
            >
              {{ (job.company_name || job.title || "?")[0]?.toUpperCase() }}
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-[12px] text-muted-foreground">
                {{ job.company_name || "Company" }}
              </p>
              <p
                class="mt-0.5 text-[15.5px] font-semibold leading-snug tracking-[-0.2px]"
              >
                {{ job.title }}
              </p>
            </div>
            <button
              class="shrink-0 p-1 transition"
              :class="
                savedJobIds.has(job.id)
                  ? 'text-primary'
                  : 'text-muted-foreground hover:text-foreground'
              "
              @click="toggleSave(job, $event)"
            >
              <Heart
                class="h-4 w-4"
                :class="savedJobIds.has(job.id) ? 'fill-primary' : ''"
              />
            </button>
          </div>

          <!-- Chips -->
          <div class="flex flex-wrap gap-1.5">
            <span
              v-if="formatSalary(job)"
              class="inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground"
              >{{ formatSalary(job) }}</span
            >
            <span
              v-if="job.type"
              class="inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground"
              >{{ job.type }}</span
            >
            <span
              v-if="job.seniority"
              class="inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground"
              >{{ job.seniority }}</span
            >
          </div>

          <!-- Footer: location · date · applicants -->
          <div
            class="flex items-center gap-2 text-[12px] text-muted-foreground"
          >
            <MapPin class="h-3 w-3 shrink-0" />
            <span class="truncate">{{ job.location }}</span>
            <span class="opacity-40">·</span>
            <span class="shrink-0">{{
              formatRelativeDate(job.created_at)
            }}</span>
            <div class="flex-1" />
            <span class="shrink-0 flex items-center gap-1">
              <Users class="h-3 w-3" />{{ job.applicants_count }}
            </span>
          </div>
        </div>
      </div>

      <!-- ─── LIST VIEW ───────────────────────────────────── -->
      <div v-else class="rounded-2xl border border-border bg-card">
        <div
          v-for="(job, i) in paginatedJobs"
          :key="job.id"
          class="flex cursor-pointer items-center gap-3.5 px-5 py-4 transition hover:bg-muted/40"
          :class="i < paginatedJobs.length - 1 ? 'border-b border-border' : ''"
          @click="router.push(`/jobs/${job.id}`)"
        >
          <img
            v-if="job.company_logo_url"
            :src="job.company_logo_url"
            :alt="job.company_name"
            class="h-[42px] w-[42px] shrink-0 rounded-xl object-cover"
          />
          <div
            v-else
            class="flex h-[42px] w-[42px] shrink-0 items-center justify-center rounded-xl text-[14px] font-bold text-white/90"
            :style="{
              background: companyGradient(job.company_name || job.title),
            }"
          >
            {{ (job.company_name || job.title || "?")[0]?.toUpperCase() }}
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-[14.5px] font-semibold">{{ job.title }}</p>
            <p class="mt-0.5 text-[12.5px] text-muted-foreground">
              {{ job.company_name || "Company" }} · {{ job.location }} ·
              {{ formatRelativeDate(job.created_at) }}
            </p>
          </div>
          <span
            v-if="formatSalary(job)"
            class="hidden shrink-0 items-center rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground sm:inline-flex"
            >{{ formatSalary(job) }}</span
          >
          <span
            v-if="job.type"
            class="hidden shrink-0 items-center rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground md:inline-flex"
            >{{ job.type }}</span
          >
          <button
            class="shrink-0 p-1.5 transition"
            :class="
              savedJobIds.has(job.id)
                ? 'text-primary'
                : 'text-muted-foreground hover:text-foreground'
            "
            @click="toggleSave(job, $event)"
          >
            <Heart
              class="h-4 w-4"
              :class="savedJobIds.has(job.id) ? 'fill-primary' : ''"
            />
          </button>
        </div>
      </div>

      <!-- ─── PAGINATION ──────────────────────────────────── -->
      <div class="flex items-center justify-center gap-1 pt-1">
        <button
          class="flex h-8 w-8 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
          :disabled="page === 1"
          @click="page--"
        >
          <ChevronLeft class="h-4 w-4" />
        </button>
        <button
          v-for="p in totalPages"
          :key="p"
          :class="[
            'h-8 w-8 rounded-full text-[13px] font-medium transition',
            p === page
              ? 'bg-primary text-background'
              : 'text-muted-foreground hover:bg-muted',
          ]"
          @click="page = p"
        >
          {{ p }}
        </button>
        <button
          class="flex h-8 w-8 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
          :disabled="page === totalPages"
          @click="page++"
        >
          <ChevronRight class="h-4 w-4" />
        </button>
      </div>
    </template>
  </div>

  <!-- Delete dialog -->
  <AlertDialog v-model:open="deleteDialogOpen">
    <AlertDialogContent class="rounded-2xl">
      <AlertDialogHeader>
        <AlertDialogTitle>Ажлын байр устгах уу?</AlertDialogTitle>
        <AlertDialogDescription>
          <span class="font-medium text-foreground">{{
            selectedJob?.title
          }}</span>
          ажлын байрыг устгана. Энэ үйлдлийг буцаах боломжгүй.
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel :disabled="isDeleting">Болих</AlertDialogCancel>
        <AlertDialogAction :disabled="isDeleting" @click="confirmDelete">
          {{ isDeleting ? "Устгаж байна..." : "Устгах" }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
