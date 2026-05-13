<script setup lang="ts">
import type {
  Job,
  Application,
  AssessmentResult,
  DutyAssessment,
  RequirementAssessment,
  Company,
} from "../../../composables/types";
import { toast } from "vue-sonner";
import {
  MapPin,
  Upload,
  CheckCircle2,
  XCircle,
  MinusCircle,
  ChevronRight,
  ChevronLeft,
  Loader2,
  Search,
  Pencil,
} from "lucide-vue-next";

definePageMeta({ middleware: "auth" });

const route = useRoute();
const router = useRouter();
const { user } = useAuth();
const jobsAPI = useJobsAPI();
const applicationsAPI = useApplicationsAPI();
const companyAPI = useCompanyAPI();

const job = ref<Job | null>(null);
const loading = ref(true);

// ── Recruiter state ────────────────────────────────────────────────────────
const isEditing = ref(false);
const company = ref<Company | null>(null);
const isSubmitting = ref(false);
const editErrorMessage = ref("");

// ── Applicant state ────────────────────────────────────────────────────────
const application = ref<Application | null>(null);
const analyzeResult = ref<AssessmentResult | null>(null);
const showAnalyzePreview = ref(false);
const applyOpen = ref(false);
const applyCvFile = ref<File | null>(null);
const isApplying = ref(false);
const analyzeOpen = ref(false);
const analyzeCvFile = ref<File | null>(null);
const isAnalyzing = ref(false);
const pollTimer = ref<ReturnType<typeof setInterval> | null>(null);

// ── Computed ───────────────────────────────────────────────────────────────
const jobID = computed(() => Number(route.params.id));
const isRecruiter = computed(() => user.value?.role === "recruiter");
const isApplicant = computed(() => user.value?.role === "user");
const recruiterCompanyID = computed(() => user.value?.company_id ?? 0);
const storageKey = computed(
  () => `analyze_result_job_${jobID.value}_user_${user.value?.id}`,
);

const displayResult = computed<AssessmentResult | null>(() => {
  if (!isApplicant.value) return null;
  // Fresh analyze preview takes priority over everything
  if (showAnalyzePreview.value && analyzeResult.value) return analyzeResult.value;
  // Non-pending application result
  if (application.value && application.value.status !== "pending") {
    return application.value as unknown as AssessmentResult;
  }
  // Pending → show nothing (spinner in header handles this)
  if (application.value?.status === "pending") return null;
  // No application yet: show cached analyze result if available
  return analyzeResult.value;
});

const isPreview = computed(
  () => showAnalyzePreview.value || !application.value || application.value.status === "pending",
);

// ── Recruiter functions ────────────────────────────────────────────────────
async function loadJobForRecruiter() {
  try {
    const [fetchedJob, fetchedCompany] = await Promise.all([
      jobsAPI.get(recruiterCompanyID.value, jobID.value),
      companyAPI.get(),
    ]);
    job.value = fetchedJob;
    company.value = fetchedCompany;
  } catch {
    toast.error("Ажлын байрны мэдээлэл ачааллаж чадсангүй");
    router.replace("/jobs");
  }
}

async function handleEditSubmit(payload: any) {
  editErrorMessage.value = "";
  isSubmitting.value = true;
  try {
    job.value = await jobsAPI.update(
      recruiterCompanyID.value,
      jobID.value,
      payload,
    );
    isEditing.value = false;
    toast.success("Ажлын байр амжилттай шинэчлэгдлээ.");
  } catch (e: any) {
    editErrorMessage.value =
      e?.data?.message || "Ажлын байр хадгалахад алдаа гарлаа.";
  } finally {
    isSubmitting.value = false;
  }
}

// ── Applicant functions ────────────────────────────────────────────────────
async function loadJob() {
  try {
    const jobs = await jobsAPI.list();
    job.value = jobs.find((j) => j.id === jobID.value) ?? null;
  } catch {
    toast.error("Ажлын байрны мэдээлэл ачааллаж чадсангүй");
  }
}

async function loadMyApplication() {
  if (!isApplicant.value) return;
  try {
    const apps = await applicationsAPI.listMine();
    application.value =
      apps.find((a) => a.job_posting_id === jobID.value) ?? null;
  } catch {
    /* not applied yet */
  }
}

function loadCachedAnalyze() {
  if (!import.meta.client || !isApplicant.value) return;
  const raw = localStorage.getItem(storageKey.value);
  if (raw) {
    try {
      analyzeResult.value = JSON.parse(raw);
    } catch {
      /* ignore */
    }
  }
}

function saveCachedAnalyze(result: AssessmentResult) {
  if (!import.meta.client) return;
  localStorage.setItem(storageKey.value, JSON.stringify(result));
}

function onApplyFileChange(e: Event) {
  applyCvFile.value =
    (e.target as HTMLInputElement).files?.[0] ?? null;
}

function onAnalyzeFileChange(e: Event) {
  analyzeCvFile.value =
    (e.target as HTMLInputElement).files?.[0] ?? null;
}

async function submitApplication() {
  if (!applyCvFile.value) {
    toast.error("CV файлаа сонгоно уу");
    return;
  }
  isApplying.value = true;
  try {
    await applicationsAPI.applyToJob(jobID.value, applyCvFile.value);
    applyOpen.value = false;
    applyCvFile.value = null;
    showAnalyzePreview.value = false;
    toast.success("Анкет амжилттай илгээгдлээ! AI үнэлгээ хийгдэж байна...");
    await loadMyApplication();
    startPolling();
  } catch (e: any) {
    toast.error(e?.data?.message ?? "Анкет илгээхэд алдаа гарлаа");
  } finally {
    isApplying.value = false;
  }
}

async function runAnalyze() {
  if (!analyzeCvFile.value) {
    toast.error("CV файлаа сонгоно уу");
    return;
  }
  isAnalyzing.value = true;
  try {
    const result = await applicationsAPI.analyzeJob(
      jobID.value,
      analyzeCvFile.value,
    );
    analyzeResult.value = result;
    showAnalyzePreview.value = true;
    saveCachedAnalyze(result);
    applyCvFile.value = analyzeCvFile.value;
    analyzeOpen.value = false;
    toast.success("Үнэлгээ амжилттай хийгдлээ!");
  } catch (e: any) {
    toast.error(e?.data?.message ?? "Үнэлгээ хийхэд алдаа гарлаа");
  } finally {
    isAnalyzing.value = false;
  }
}

function startPolling() {
  stopPolling();
  pollTimer.value = setInterval(async () => {
    if (application.value?.status !== "pending") {
      stopPolling();
      return;
    }
    const apps = await applicationsAPI.listMine().catch(() => []);
    const updated = apps.find((a) => a.job_posting_id === jobID.value);
    if (updated) application.value = updated;
    if (updated?.status !== "pending") stopPolling();
  }, 4000);
}

function stopPolling() {
  if (pollTimer.value) {
    clearInterval(pollTimer.value);
    pollTimer.value = null;
  }
}

onUnmounted(stopPolling);

// ── Init ───────────────────────────────────────────────────────────────────
loading.value = true;
if (isRecruiter.value) {
  await loadJobForRecruiter();
} else {
  await loadJob();
  await loadMyApplication();
  loadCachedAnalyze();
}
loading.value = false;

if (application.value?.status === "pending") startPolling();

// ── Helpers ────────────────────────────────────────────────────────────────
function scoreColor(score: number) {
  if (score >= 75) return "text-green-600";
  if (score >= 50) return "text-yellow-600";
  return "text-red-600";
}

function scoreBg(score: number) {
  if (score >= 75) return "bg-green-50 border-green-200";
  if (score >= 50) return "bg-yellow-50 border-yellow-200";
  return "bg-red-50 border-red-200";
}

function statusIcon(status: string) {
  if (status === "met") return CheckCircle2;
  if (status === "partial") return MinusCircle;
  return XCircle;
}

function statusColor(status: string) {
  if (status === "met") return "text-green-600";
  if (status === "partial") return "text-yellow-600";
  return "text-red-500";
}

function statusBg(status: string) {
  if (status === "met") return "bg-green-50 border-green-100";
  if (status === "partial") return "bg-yellow-50 border-yellow-100";
  return "bg-red-50 border-red-100";
}

function dutyAssessFor(dutyText: string): DutyAssessment | null {
  return (
    displayResult.value?.duty_assessments?.find(
      (d) => d.duty.trim() === dutyText.trim(),
    ) ?? null
  );
}

function reqAssessFor(reqText: string): RequirementAssessment | null {
  return (
    displayResult.value?.requirement_assessments?.find(
      (r) => r.requirement.trim() === reqText.trim(),
    ) ?? null
  );
}

function jobStatusLabel(status: string) {
  if (status === "posted") return "Нийтлэгдсэн";
  if (status === "closed") return "Хаагдсан";
  return "Ноорог";
}
</script>

<template>
  <!-- Loading -->
  <div v-if="loading" class="flex justify-center py-20">
    <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
  </div>

  <!-- Not found -->
  <div
    v-else-if="!job"
    class="rounded-3xl border border-dashed border-border px-6 py-20 text-center text-muted-foreground"
  >
    Ажлын байр олдсонгүй
  </div>

  <!-- ── Recruiter: Edit mode ───────────────────────────────────────────── -->
  <template v-else-if="isRecruiter && isEditing">
    <div class="space-y-6">
      <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
        <div class="flex items-center gap-3">
          <Button variant="ghost" size="icon" class="rounded-xl" @click="isEditing = false">
            <ChevronLeft class="h-4 w-4" />
          </Button>
          <div>
            <p class="text-sm font-medium text-muted-foreground">Ажлын байр засах</p>
            <h1 class="text-2xl font-semibold tracking-tight">{{ job.title }}</h1>
          </div>
        </div>
      </section>

      <Card class="rounded-3xl border-border shadow-sm">
        <CardContent class="pt-6">
          <JobForm
            :company="company"
            :is-submitting="isSubmitting"
            :error-message="editErrorMessage"
            :initial-title="job.title"
            :initial-contact-info="job.contact_info"
            :initial-type="job.type || job.employment_type"
            :initial-level="job.level || job.seniority"
            :initial-status="job.status"
            :initial-city="job.city"
            :initial-district="job.district"
            :initial-location-x="job.location_x ? String(job.location_x) : ''"
            :initial-location-y="job.location_y ? String(job.location_y) : ''"
            :initial-min-salary="job.min_salary ? String(job.min_salary) : ''"
            :initial-max-salary="job.max_salary ? String(job.max_salary) : ''"
            :initial-additional-info="job.additional_info || job.description"
            :initial-duties="job.duties"
            :initial-requirements="job.requirements"
            :initial-skills="job.skills"
            :initial-bonuses="job.bonuses"
            @submit="handleEditSubmit"
          >
            <template #actions>
              <Button type="button" variant="outline" :disabled="isSubmitting" @click="isEditing = false">
                Болих
              </Button>
              <Button type="submit" :disabled="isSubmitting">
                {{ isSubmitting ? "Хадгалж байна..." : "Шинэчлэх" }}
              </Button>
            </template>
          </JobForm>
        </CardContent>
      </Card>
    </div>
  </template>

  <!-- ── View mode (recruiter read / applicant) ─────────────────────────── -->
  <template v-else>
    <div class="space-y-6">

      <!-- Header: Recruiter -->
      <section v-if="isRecruiter" class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
        <div class="flex items-center gap-2 mb-4">
          <Button variant="ghost" size="icon" class="rounded-xl" @click="router.push('/jobs')">
            <ChevronLeft class="h-4 w-4" />
          </Button>
          <span class="text-sm text-muted-foreground">Ажлын байрууд</span>
        </div>
        <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div class="space-y-2">
            <h1 class="text-3xl font-semibold tracking-tight">{{ job.title }}</h1>
            <div class="flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
              <span class="flex items-center gap-1">
                <MapPin class="h-4 w-4" /> {{ job.location }}
              </span>
              <Badge variant="outline" class="rounded-full px-3">{{ job.type || job.employment_type }}</Badge>
              <Badge variant="outline" class="rounded-full px-3">{{ job.level || job.seniority }}</Badge>
              <Badge class="rounded-full px-3">{{ jobStatusLabel(job.status) }}</Badge>
            </div>
          </div>
          <Button class="rounded-full px-5 shrink-0" @click="isEditing = true">
            <Pencil class="mr-2 h-4 w-4" /> Засах
          </Button>
        </div>
      </section>

      <!-- Header: Applicant -->
      <section v-else class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div class="space-y-2">
            <h1 class="text-3xl font-semibold tracking-tight">{{ job.title }}</h1>
            <div class="flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
              <span class="flex items-center gap-1">
                <MapPin class="h-4 w-4" /> {{ job.location }}
              </span>
              <Badge variant="outline" class="rounded-full px-3">{{ job.type || job.employment_type }}</Badge>
              <Badge variant="outline" class="rounded-full px-3">{{ job.level || job.seniority }}</Badge>
            </div>
            <p v-if="job.min_salary || job.max_salary" class="text-sm font-medium">
              {{ job.min_salary?.toLocaleString() }}₮ – {{ job.max_salary?.toLocaleString() }}₮
            </p>
          </div>

          <div class="flex flex-col gap-2 sm:flex-row sm:items-start">
            <Button variant="outline" class="rounded-full px-5" @click="analyzeOpen = true">
              <Search class="mr-2 h-4 w-4" /> Нийтлэл шалгах
            </Button>

            <template v-if="!application">
              <Button class="rounded-full px-5" @click="applyOpen = true">
                <Upload class="mr-2 h-4 w-4" /> CV илгээх
              </Button>
            </template>
            <template v-else-if="application.status === 'pending'">
              <div class="flex items-center gap-2 rounded-full border border-border px-4 py-2 text-sm text-muted-foreground">
                <Loader2 class="h-4 w-4 animate-spin" /> Үнэлгээ хийгдэж байна...
              </div>
            </template>
            <template v-else>
              <div class="flex flex-col gap-1.5 sm:items-end">
                <Badge variant="secondary" class="rounded-full px-4 py-1.5">Анкет илгээсэн</Badge>
                <button
                  class="text-xs text-muted-foreground hover:text-foreground underline underline-offset-2"
                  @click="applyOpen = true"
                >
                  Дахин илгээх
                </button>
              </div>
            </template>
          </div>
        </div>
      </section>

      <!-- Assessment score banner (applicants only) -->
      <div
        v-if="displayResult"
        :class="['rounded-3xl border px-6 py-5 shadow-sm', scoreBg(displayResult.overall_score)]"
      >
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:gap-6">
          <div class="shrink-0 text-center sm:text-left">
            <p class="text-xs font-medium uppercase tracking-widest text-muted-foreground">
              {{ isPreview ? 'Урьдчилсан үнэлгээ' : 'Таны үнэлгээ' }}
            </p>
            <p :class="['text-5xl font-bold leading-none mt-1', scoreColor(displayResult.overall_score)]">
              {{ displayResult.overall_score }}<span class="text-2xl">%</span>
            </p>
          </div>
          <div class="flex-1 space-y-1">
            <p class="text-sm font-medium">
              {{ isPreview ? 'Нийтлэлтэй тохирол' : 'Ажлын байрны шаардлагатай тохирол' }}
            </p>
            <p class="text-sm leading-6 text-muted-foreground">{{ displayResult.summary }}</p>
          </div>
          <div v-if="!isPreview" class="shrink-0">
            <Badge :class="[
              'rounded-full px-4 py-1.5 text-sm',
              (application?.status === 'shortlisted') ? 'bg-green-100 text-green-800' :
              (application?.status === 'rejected') ? 'bg-red-100 text-red-800' :
              'bg-blue-100 text-blue-800',
            ]">
              {{ application?.status === 'shortlisted' ? 'Сонгогдсон' :
                 application?.status === 'rejected' ? 'Татгалзсан' : 'Үнэлэгдсэн' }}
            </Badge>
          </div>
        </div>
      </div>

      <!-- Job body -->
      <div class="grid gap-6 lg:grid-cols-3">
        <div class="space-y-6 lg:col-span-2">
          <!-- Description -->
          <Card class="rounded-3xl border-border shadow-sm">
            <CardHeader><CardTitle>Тайлбар</CardTitle></CardHeader>
            <CardContent>
              <p class="whitespace-pre-line text-sm leading-7 text-muted-foreground">
                {{ job.additional_info || job.description }}
              </p>
            </CardContent>
          </Card>

          <!-- Duties -->
          <Card v-if="job.duties?.length" class="rounded-3xl border-border shadow-sm">
            <CardHeader><CardTitle>Үүрэг хариуцлага</CardTitle></CardHeader>
            <CardContent class="space-y-2">
              <div
                v-for="duty in job.duties"
                :key="duty"
                class="rounded-2xl border border-border p-3 transition"
                :class="dutyAssessFor(duty) ? statusBg(dutyAssessFor(duty)!.status) : 'bg-muted/20'"
              >
                <div class="flex items-start gap-2">
                  <component
                    :is="dutyAssessFor(duty) ? statusIcon(dutyAssessFor(duty)!.status) : ChevronRight"
                    class="mt-0.5 h-4 w-4 shrink-0"
                    :class="dutyAssessFor(duty) ? statusColor(dutyAssessFor(duty)!.status) : 'text-primary'"
                  />
                  <div class="flex-1">
                    <p class="text-sm" :class="dutyAssessFor(duty) ? 'font-medium' : 'text-muted-foreground'">
                      {{ duty }}
                    </p>
                    <p
                      v-if="dutyAssessFor(duty)"
                      class="mt-0.5 text-xs"
                      :class="statusColor(dutyAssessFor(duty)!.status)"
                    >
                      {{ dutyAssessFor(duty)!.explanation }}
                    </p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- Requirements -->
          <Card v-if="job.requirements?.length" class="rounded-3xl border-border shadow-sm">
            <CardHeader><CardTitle>Шаардлагууд</CardTitle></CardHeader>
            <CardContent class="space-y-2">
              <div
                v-for="req in job.requirements"
                :key="req"
                class="rounded-2xl border border-border p-3 transition"
                :class="reqAssessFor(req) ? statusBg(reqAssessFor(req)!.status) : 'bg-muted/20'"
              >
                <div class="flex items-start gap-2">
                  <component
                    :is="reqAssessFor(req) ? statusIcon(reqAssessFor(req)!.status) : ChevronRight"
                    class="mt-0.5 h-4 w-4 shrink-0"
                    :class="reqAssessFor(req) ? statusColor(reqAssessFor(req)!.status) : 'text-primary'"
                  />
                  <div class="flex-1">
                    <p class="text-sm" :class="reqAssessFor(req) ? 'font-medium' : 'text-muted-foreground'">
                      {{ req }}
                    </p>
                    <p
                      v-if="reqAssessFor(req)"
                      class="mt-0.5 text-xs"
                      :class="statusColor(reqAssessFor(req)!.status)"
                    >
                      {{ reqAssessFor(req)!.explanation }}
                    </p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- Recommendations (applicant preview result) -->
          <Card
            v-if="displayResult?.recommendations?.length"
            class="rounded-3xl border-border shadow-sm"
          >
            <CardHeader><CardTitle>Сайжруулах зөвлөмж</CardTitle></CardHeader>
            <CardContent>
              <ul class="space-y-2">
                <li
                  v-for="rec in displayResult.recommendations"
                  :key="rec"
                  class="flex items-start gap-2 text-sm text-muted-foreground"
                >
                  <ChevronRight class="mt-0.5 h-4 w-4 shrink-0 text-primary" /> {{ rec }}
                </li>
              </ul>
            </CardContent>
          </Card>
        </div>

        <!-- Right column -->
        <div class="space-y-6">
          <Card v-if="job.skills?.length" class="rounded-3xl border-border shadow-sm">
            <CardHeader><CardTitle>Шаардлагатай ур чадвар</CardTitle></CardHeader>
            <CardContent>
              <div v-if="displayResult" class="space-y-3">
                <div v-if="displayResult.matched_skills?.length">
                  <p class="mb-2 text-xs font-medium uppercase tracking-wide text-green-600">
                    Тохирсон ур чадварууд
                  </p>
                  <div class="space-y-2">
                    <div
                      v-for="s in displayResult.matched_skills"
                      :key="s.skill"
                      class="rounded-xl bg-green-50 border border-green-100 px-3 py-2"
                    >
                      <p class="flex items-center gap-1.5 text-xs font-medium text-green-800">
                        <CheckCircle2 class="h-3.5 w-3.5 shrink-0" /> {{ s.skill }}
                      </p>
                      <p class="mt-0.5 text-xs text-green-600">{{ s.explanation }}</p>
                    </div>
                  </div>
                </div>
                <div v-if="displayResult.missing_skills?.length">
                  <p class="mb-2 text-xs font-medium uppercase tracking-wide text-red-500">
                    Дутуу ур чадварууд
                  </p>
                  <div class="space-y-2">
                    <div
                      v-for="s in displayResult.missing_skills"
                      :key="s.skill"
                      class="rounded-xl bg-red-50 border border-red-100 px-3 py-2"
                    >
                      <p class="flex items-center gap-1.5 text-xs font-medium text-red-800">
                        <XCircle class="h-3.5 w-3.5 shrink-0" /> {{ s.skill }}
                      </p>
                      <p class="mt-0.5 text-xs text-red-600">{{ s.explanation }}</p>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="flex flex-wrap gap-2">
                <Badge
                  v-for="skill in job.skills"
                  :key="skill"
                  variant="secondary"
                  class="rounded-full px-3"
                >
                  {{ skill }}
                </Badge>
              </div>
            </CardContent>
          </Card>

          <Card v-if="job.bonuses?.length" class="rounded-3xl border-border shadow-sm">
            <CardHeader><CardTitle>Нэмэлт давуу тал</CardTitle></CardHeader>
            <CardContent>
              <ul class="space-y-1">
                <li v-for="bonus in job.bonuses" :key="bonus" class="text-sm text-muted-foreground">
                  ✓ {{ bonus }}
                </li>
              </ul>
            </CardContent>
          </Card>

          <Card v-if="job.contact_info" class="rounded-3xl border-border shadow-sm">
            <CardHeader><CardTitle>Холбоо барих</CardTitle></CardHeader>
            <CardContent>
              <p class="text-sm text-muted-foreground">{{ job.contact_info }}</p>
            </CardContent>
          </Card>
        </div>
      </div>

    </div>
  </template>

  <!-- ── Applicant dialogs ──────────────────────────────────────────────── -->
  <Dialog v-model:open="applyOpen">
    <DialogContent class="rounded-3xl sm:max-w-md">
      <DialogHeader>
        <DialogTitle>CV илгээх</DialogTitle>
        <DialogDescription>
          CV файлаа PDF хэлбэрээр оруулна уу. AI таны CV-г ажлын байрны шаардлагатай харьцуулж үнэлнэ.
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-2">
        <div class="space-y-2">
          <Label>CV файл (PDF)</Label>
          <div
            v-if="applyCvFile"
            class="flex items-center justify-between gap-3 rounded-2xl border border-border bg-muted/30 px-3 py-2"
          >
            <p class="truncate text-xs text-muted-foreground">
              <span class="font-medium text-foreground">{{ applyCvFile.name }}</span>
              <span class="ml-2 text-green-600">— нийтлэл шалгасан CV</span>
            </p>
            <button
              type="button"
              class="shrink-0 text-xs text-primary underline underline-offset-2"
              @click="applyCvFile = null"
            >
              Өөр файл
            </button>
          </div>
          <Input v-else type="file" accept=".pdf" @change="onApplyFileChange" />
        </div>
        <p v-if="application" class="rounded-2xl bg-muted/30 px-3 py-2 text-xs text-muted-foreground">
          Өмнөх анкет шинэ CV-гээр солигдоно.
        </p>
      </div>
      <DialogFooter>
        <Button variant="outline" :disabled="isApplying" @click="applyOpen = false">Болих</Button>
        <Button :disabled="isApplying || !applyCvFile" @click="submitApplication">
          <Loader2 v-if="isApplying" class="mr-2 h-4 w-4 animate-spin" />
          {{ isApplying ? "Илгээж байна..." : "Илгээх" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>

  <Dialog v-model:open="analyzeOpen">
    <DialogContent class="rounded-3xl sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Нийтлэл шалгах</DialogTitle>
        <DialogDescription>
          CV файлаа оруулж ажлын байранд хэр тохирохыг урьдчилан үнэлүүлнэ үү. Анкет илгээхгүй — үр дүн хуудсанд харагдана.
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-2">
        <div class="space-y-2">
          <Label>CV файл (PDF)</Label>
          <Input type="file" accept=".pdf" @change="onAnalyzeFileChange" />
          <p v-if="analyzeCvFile" class="text-xs text-muted-foreground">Сонгосон: {{ analyzeCvFile.name }}</p>
        </div>
        <div
          v-if="isAnalyzing"
          class="flex items-center gap-3 rounded-2xl bg-muted/30 px-4 py-3 text-sm text-muted-foreground"
        >
          <Loader2 class="h-4 w-4 animate-spin" />
          AI үнэлгээ хийгдэж байна... Хэдэн секунд хүлээнэ үү.
        </div>
      </div>
      <DialogFooter>
        <Button variant="outline" :disabled="isAnalyzing" @click="analyzeOpen = false">Болих</Button>
        <Button :disabled="isAnalyzing || !analyzeCvFile" @click="runAnalyze">
          <Loader2 v-if="isAnalyzing" class="mr-2 h-4 w-4 animate-spin" />
          {{ isAnalyzing ? "Үнэлж байна..." : "Үнэлгээ хийх" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
