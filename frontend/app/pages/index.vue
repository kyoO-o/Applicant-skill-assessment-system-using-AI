<script setup lang="ts">
import { JobStatus } from "../composables/types";
import { toast } from "vue-sonner";
import {
  BriefcaseBusiness,
  FileText,
  Sparkles,
  Calendar,
  ClipboardList,
  Users,
  ArrowRight,
} from "lucide-vue-next";

const { user } = useAuth();
const jobsAPI = useJobsAPI();
const router = useRouter();

const recruiterJobs = ref<Job[]>([]);
const loading = ref(false);

const recruiterStats = computed(() => {
  const totalJobs = recruiterJobs.value.length;
  const activeJobs = recruiterJobs.value.filter(
    (job) => job.status === JobStatus.Posted,
  ).length;
  const totalApplicants = recruiterJobs.value.reduce(
    (sum, job) => sum + job.applicants_count,
    0,
  );
  return { totalJobs, activeJobs, totalApplicants };
});

const recruiterNeedsCompany = computed(
  () => user.value?.role === "recruiter" && !user.value?.company_id,
);

function openRecruiterJobs() {
  if (recruiterNeedsCompany.value) {
    toast.warning("Ажлын байр удирдахаас өмнө компанийн мэдээллээ бүртгэнэ үү.");
    router.push("/profile");
    return;
  }
  router.push("/jobs");
}

async function loadRecruiterDashboard() {
  if (user.value?.role !== "recruiter") return;
  loading.value = true;
  try {
    recruiterJobs.value = await jobsAPI.list();
  } finally {
    loading.value = false;
  }
}

await loadRecruiterDashboard();
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
        style="background: linear-gradient(135deg, var(--primary), oklch(0.348 0.106 295))"
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
          <div class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted text-foreground">
            <BriefcaseBusiness class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Нийт ажлын байр</p>
            <p class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]">
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
          <div class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted text-foreground">
            <Users class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Нийт горилогч</p>
            <p class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]">
              {{ loading ? "…" : recruiterStats.totalApplicants }}
            </p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-success-foreground">Шинэ өргөдлүүд</p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl"
            style="background: oklch(0.408 0.124 295 / 8%); color: var(--primary)"
          >
            <Sparkles class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">AI дундаж оноо</p>
            <p class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]">—</p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-primary">AI үнэлгээ</p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted text-foreground">
            <Calendar class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Ярилцлага</p>
            <p class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]">—</p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-muted-foreground">Энэ долоо хоног</p>
      </div>
    </div>

    <!-- Jobs list + Activity -->
    <div class="grid gap-4 lg:grid-cols-[1.3fr_0.9fr]">
      <!-- Jobs -->
      <div class="rounded-2xl border border-border bg-card">
        <div class="flex items-center justify-between border-b border-border px-5 py-4">
          <div>
            <p class="text-[15px] font-semibold tracking-[-0.2px]">Ажлын байрны хөдөлгөөн</p>
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
              <p class="truncate text-[14.5px] font-semibold">{{ job.title }}</p>
              <p class="mt-0.5 text-[12.5px] text-muted-foreground">
                {{ job.location }} · {{ job.employment_type || "Тогтоогдоогүй" }}
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
              <span class="mr-1.5 h-1.5 w-1.5 rounded-full"
                :class="job.status === JobStatus.Posted ? 'bg-success' : 'bg-muted-foreground'"
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

      <!-- Activity -->
      <div class="rounded-2xl border border-border bg-card">
        <div class="border-b border-border px-5 py-4">
          <p class="text-[15px] font-semibold tracking-[-0.2px]">Сүүлийн үйл ажиллагаа</p>
          <p class="mt-0.5 text-[12px] text-muted-foreground">Ажлын байрны хөдөлгөөний товч тойм</p>
        </div>
        <div class="p-3">
          <div
            v-for="job in recruiterJobs.slice(0, 3)"
            :key="'act-' + job.id"
            class="mb-1.5 rounded-xl border border-border bg-muted/30 p-3.5"
          >
            <p class="text-[13px] font-semibold">{{ job.title }}</p>
            <p class="mt-1 text-[12px] text-muted-foreground">
              {{ job.company_name || "Таны компани" }} · {{ job.location }}
            </p>
            <span
              class="mt-2.5 inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-semibold capitalize"
              :class="
                job.status === JobStatus.Posted
                  ? 'badge-success'
                  : 'bg-muted text-muted-foreground'
              "
            >
              {{ job.status }}
            </span>
          </div>
          <div
            v-if="!recruiterJobs.length && !loading"
            class="px-4 py-10 text-center text-[13px] text-muted-foreground"
          >
            Ажлын байр үүсгэсний дараа энд харагдана.
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
        style="background: linear-gradient(135deg, var(--primary), oklch(0.348 0.106 295))"
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
          <div class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted">
            <BriefcaseBusiness class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Идэвхтэй анкет</p>
            <p class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]">—</p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-success-foreground">Илгээсэн</p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl"
            style="background: oklch(0.408 0.124 295 / 8%); color: var(--primary)"
          >
            <Sparkles class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Дундаж AI оноо</p>
            <p class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]">—</p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-primary">AI үнэлгээ</p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted">
            <Calendar class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Ярилцлага</p>
            <p class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]">—</p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-muted-foreground">Хуваарьт</p>
      </div>

      <div class="rounded-2xl border border-border bg-card p-[18px]">
        <div class="flex items-start gap-3">
          <div class="flex h-9 w-9 items-center justify-center rounded-xl bg-muted">
            <ClipboardList class="h-4 w-4" />
          </div>
          <div class="flex-1">
            <p class="text-[12px] text-muted-foreground">Даалгавар</p>
            <p class="mt-0.5 text-[28px] font-semibold leading-none tracking-[-0.8px]">—</p>
          </div>
        </div>
        <p class="mt-2.5 text-[11.5px] font-medium text-warning-foreground">Хүлээгдэж байгаа</p>
      </div>
    </div>

    <!-- Feature cards + AI suggestions -->
    <div class="grid gap-4 lg:grid-cols-3">
      <div
        class="flex flex-col gap-3 rounded-2xl border border-border bg-card p-5 cursor-pointer transition hover:border-primary/30"
        @click="router.push('/jobs')"
      >
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-muted">
          <BriefcaseBusiness class="h-5 w-5 text-muted-foreground" />
        </div>
        <div>
          <p class="text-[15px] font-semibold">Тохирох ажлыг олоорой</p>
          <p class="mt-1.5 text-[13px] leading-[1.5] text-muted-foreground">
            Ажил олгогчдийн зарласан нээлттэй байр дундаас өөрт тохирохыг нь олоорой.
          </p>
        </div>
        <div class="mt-auto flex items-center gap-1 text-[13px] font-medium text-primary">
          Ажил хайх <ArrowRight class="h-3.5 w-3.5" />
        </div>
      </div>

      <div
        class="flex flex-col gap-3 rounded-2xl border border-border bg-card p-5 cursor-pointer transition hover:border-primary/30"
        @click="router.push('/applications')"
      >
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-muted">
          <FileText class="h-5 w-5 text-muted-foreground" />
        </div>
        <div>
          <p class="text-[15px] font-semibold">Бэлэн байдлаа хянаарай</p>
          <p class="mt-1.5 text-[13px] leading-[1.5] text-muted-foreground">
            Үнэлгээ, профайл мэдээлэл болон ажлын нээлтийг нэг дор хадгалаарай.
          </p>
        </div>
        <div class="mt-auto flex items-center gap-1 text-[13px] font-medium text-primary">
          Анкет харах <ArrowRight class="h-3.5 w-3.5" />
        </div>
      </div>

      <!-- AI suggestions card -->
      <div
        class="flex flex-col gap-3 rounded-2xl border p-5 cursor-pointer transition"
        style="background: oklch(0.408 0.124 295 / 4.5%); border-color: oklch(0.408 0.124 295 / 13%);"
        @click="router.push('/chat')"
      >
        <div class="flex items-center gap-3">
          <div
            class="flex h-10 w-10 items-center justify-center rounded-xl"
            style="background: linear-gradient(135deg, var(--primary), oklch(0.348 0.106 295)); color: #fff;"
          >
            <Sparkles class="h-5 w-5" />
          </div>
          <div>
            <p class="text-[15px] font-semibold">AI туслагч</p>
            <p class="text-[12px] text-muted-foreground">Таны CV-д тулгуурласан</p>
          </div>
        </div>
        <p class="text-[13px] leading-[1.5] text-muted-foreground">
          Ярилцлагад бэлэн байгаарай. Анкетаас эхлээд дүгнэлт хүртэлх бүх шатыг энэ платформоор хянаарай.
        </p>
        <div class="mt-auto flex items-center gap-1 text-[13px] font-medium text-primary">
          Асуух <ArrowRight class="h-3.5 w-3.5" />
        </div>
      </div>
    </div>
  </div>
</template>
