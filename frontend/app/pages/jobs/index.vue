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
} from "lucide-vue-next";

const { user } = useAuth();
const jobsAPI = useJobsAPI();
const router = useRouter();

const jobs = ref<Job[]>([]);
const selectedJob = ref<Job | null>(null);
const loading = ref(false);
const deleteDialogOpen = ref(false);
const isDeleting = ref(false);

const isRecruiter = computed(() => user.value?.role === "recruiter");
const recruiterCompanyID = computed(() => user.value?.company_id ?? 0);
const hasRecruiterCompany = computed(() => recruiterCompanyID.value > 0);

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
</script>

<template>
  <!-- ── RECRUITER VIEW ──────────────────────────────────── -->
  <div v-if="isRecruiter" class="space-y-5">
    <!-- Page header -->
    <div class="flex items-center justify-between gap-4">
      <div>
        <h1 class="mt-1 text-[26px] font-semibold tracking-[-0.6px]">
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
          v-for="job in jobs"
          :key="job.id"
          class="flex items-center gap-4 border-b border-border px-5 py-4 last:border-0 cursor-pointer transition hover:bg-muted/30"
          @click="router.push(`/jobs/${job.id}`)"
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
  </div>

  <!-- ── APPLICANT VIEW ──────────────────────────────────── -->
  <div v-else class="space-y-5">
    <!-- Page header -->
    <div>
      <p
        class="text-[12px] font-semibold uppercase tracking-[0.8px] text-muted-foreground"
      >
        Ажлын байр
      </p>
      <h1 class="mt-1 text-[26px] font-semibold tracking-[-0.6px]">
        Нээлттэй ажлын байрууд
      </h1>
      <p class="mt-1 text-[13.5px] text-muted-foreground">
        Идэвхтэй нийтлэгдсэн ажлын байруудыг харж, анкетаа илгээгээрэй.
      </p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="h-7 w-7 animate-spin text-primary" />
    </div>

    <!-- Empty -->
    <div
      v-else-if="!applicantJobs.length"
      class="flex flex-col items-center py-20 text-center"
    >
      <div
        class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-muted"
      >
        <BriefcaseBusiness class="h-7 w-7 text-muted-foreground" />
      </div>
      <p class="text-[16px] font-semibold">Идэвхтэй ажлын байр байхгүй байна</p>
      <p class="mt-2 text-[13.5px] text-muted-foreground">
        Дараа дахин шалгана уу.
      </p>
    </div>

    <!-- Job cards grid -->
    <div v-else class="grid gap-4 lg:grid-cols-2">
      <div
        v-for="job in applicantJobs"
        :key="job.id"
        class="group flex flex-col rounded-2xl border border-border bg-card p-5 transition-all hover:border-primary/30 hover:shadow-sm cursor-pointer"
        @click="router.push(`/jobs/${job.id}`)"
      >
        <!-- Header -->
        <div class="flex items-start justify-between gap-3">
          <div class="flex-1 min-w-0">
            <h3 class="text-[16px] font-semibold leading-snug">
              {{ job.title }}
            </h3>
            <p class="mt-0.5 text-[13px] text-muted-foreground">
              {{ job.company_name || "Компани" }}
            </p>
          </div>
          <span
            class="inline-flex shrink-0 items-center rounded-full px-2.5 py-1 text-[11.5px] font-semibold badge-success"
          >
            <span class="mr-1.5 h-1.5 w-1.5 rounded-full bg-success" />
            Нийтлэгдсэн
          </span>
        </div>

        <!-- Tags -->
        <div class="mt-3 flex flex-wrap gap-1.5">
          <span
            class="inline-flex items-center gap-1 rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground"
          >
            <MapPin class="h-3 w-3" /> {{ job.location }}
          </span>
          <span
            v-if="job.employment_type"
            class="inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground"
          >
            {{ job.employment_type }}
          </span>
          <span
            v-if="job.seniority"
            class="inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground"
          >
            {{ job.seniority }}
          </span>
        </div>

        <!-- Salary -->
        <p
          v-if="job.min_salary || job.max_salary"
          class="mt-3 text-[14px] font-semibold text-foreground"
        >
          {{ job.min_salary?.toLocaleString() }}₮
          <span v-if="job.max_salary">
            – {{ job.max_salary?.toLocaleString() }}₮</span
          >
        </p>

        <!-- Description -->
        <p
          class="mt-2.5 flex-1 text-[13.5px] leading-[1.55] text-muted-foreground line-clamp-3"
        >
          {{ job.description }}
        </p>

        <!-- CTA -->
        <div
          class="mt-4 flex items-center justify-between border-t border-border pt-3.5"
        >
          <p class="text-[12px] text-muted-foreground">
            {{ formatDate(job.created_at) }}
          </p>
          <span
            class="text-[13px] font-semibold text-primary group-hover:underline"
          >
            Дэлгэрэнгүй →
          </span>
        </div>
      </div>
    </div>
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
