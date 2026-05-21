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
  Users,
  Search,
  Pencil,
  Sparkles,
} from "lucide-vue-next";
import { success } from "zod/v4";

definePageMeta({ middleware: "auth" });

const route = useRoute();
const router = useRouter();
const { user } = useAuth();
const jobsAPI = useJobsAPI();
const applicationsAPI = useApplicationsAPI();
const companyAPI = useCompanyAPI();
const cvAPI = useCVAPI();

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
const applyUseProfile = ref(false);
const isApplying = ref(false);
const analyzeOpen = ref(false);
const analyzeCvFile = ref<File | null>(null);
const analyzeUseProfile = ref(false);
const isAnalyzing = ref(false);
const pollTimer = ref<ReturnType<typeof setInterval> | null>(null);
const showJobMap = ref(false);
const hasSavedCV = ref(false);

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
  if (showAnalyzePreview.value && analyzeResult.value)
    return analyzeResult.value;
  if (application.value && application.value.status !== "pending") {
    return application.value as unknown as AssessmentResult;
  }
  if (application.value?.status === "pending") return null;
  return analyzeResult.value;
});

const isPreview = computed(
  () =>
    showAnalyzePreview.value ||
    !application.value ||
    application.value.status === "pending",
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
    const msg = e?.data?.message || "Ажлын байр хадгалахад алдаа гарлаа.";
    editErrorMessage.value = msg;
    toast.error(msg);
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

async function checkSavedCV() {
  if (!isApplicant.value) return;
  try {
    const cv = await cvAPI.get();
    hasSavedCV.value = !!(
      cv.first_name ||
      cv.last_name ||
      cv.about ||
      cv.work_experiences?.length
    );
  } catch {
    hasSavedCV.value = false;
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

function clearCachedAnalyze() {
  if (!import.meta.client) return;
  localStorage.removeItem(storageKey.value);
}

function onApplyFileChange(e: Event) {
  applyCvFile.value = (e.target as HTMLInputElement).files?.[0] ?? null;
}

function onAnalyzeFileChange(e: Event) {
  analyzeCvFile.value = (e.target as HTMLInputElement).files?.[0] ?? null;
}

async function submitApplication() {
  if (!applyUseProfile.value && !applyCvFile.value) {
    toast.error("CV файлаа сонгоно уу эсвэл хадгалагдсан CV ашиглана уу");
    return;
  }
  isApplying.value = true;
  try {
    if (applyUseProfile.value) {
      await applicationsAPI.applyFromProfile(jobID.value);
    } else {
      await applicationsAPI.applyToJob(jobID.value, applyCvFile.value!);
    }
    applyOpen.value = false;
    applyCvFile.value = null;
    applyUseProfile.value = false;
    showAnalyzePreview.value = false;
    analyzeResult.value = null;
    application.value = null;
    clearCachedAnalyze();
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
  if (!analyzeUseProfile.value && !analyzeCvFile.value) {
    toast.error("CV файлаа сонгоно уу эсвэл хадгалагдсан CV ашиглана уу");
    return;
  }
  isAnalyzing.value = true;
  try {
    let result: AssessmentResult;
    if (analyzeUseProfile.value) {
      result = await applicationsAPI.analyzeFromProfile(jobID.value);
    } else {
      result = await applicationsAPI.analyzeJob(
        jobID.value,
        analyzeCvFile.value!,
      );
    }
    analyzeResult.value = result;
    showAnalyzePreview.value = true;
    saveCachedAnalyze(result);
    if (!analyzeUseProfile.value) applyCvFile.value = analyzeCvFile.value;
    analyzeOpen.value = false;
    analyzeUseProfile.value = false;
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
  await Promise.all([loadMyApplication(), checkSavedCV()]);
  loadCachedAnalyze();
}
loading.value = false;

if (application.value?.status === "pending") startPolling();

// ── Helpers ────────────────────────────────────────────────────────────────
function assessIcon(s: string) {
  if (s === "met") return CheckCircle2;
  if (s === "partial") return MinusCircle;
  return XCircle;
}

function assessColor(s: string) {
  if (s === "met") return "text-success-foreground";
  if (s === "partial") return "text-warning-foreground";
  return "text-destructive";
}

function assessRowBg(s: string) {
  if (s === "met") return "bg-success-bg border-success/20";
  if (s === "partial") return "bg-warning-bg border-warning/20";
  return "bg-destructive/5 border-destructive/15";
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
    <Loader2 class="h-8 w-8 animate-spin text-primary" />
  </div>

  <!-- Not found -->
  <div v-else-if="!job" class="flex flex-col items-center py-20 text-center">
    <p class="text-[16px] font-semibold">Ажлын байр олдсонгүй</p>
  </div>

  <!-- ── Recruiter: Edit mode ───────────────────────────────────────────── -->
  <template v-else-if="isRecruiter && isEditing">
    <div class="space-y-5">
      <div class="flex items-center gap-3">
        <!-- <Button
          variant="outline"
          size="icon"
          class="h-9 w-9 rounded-xl text-muted-foreground hover:text-foreground"
          @click="isEditing = false"
        >
          <ChevronLeft class="h-4 w-4" />
        </Button> -->
        <div>
          <p
            class="text-[12px] font-semibold uppercase tracking-[0.8px] text-muted-foreground"
          >
            Засах
          </p>
          <h1 class="text-[22px] font-semibold tracking-[-0.5px]">
            {{ job.title }}
          </h1>
        </div>
      </div>

      <div class="rounded-2xl border border-border bg-card p-6">
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
            <Button
              type="button"
              variant="outline"
              class="rounded-full text-[13.5px]"
              :disabled="isSubmitting"
              @click="isEditing = false"
            >
              Болих
            </Button>
            <Button
              type="submit"
              class="rounded-full text-[13.5px] text-white hover:opacity-90"
              style="
                background: linear-gradient(
                  135deg,
                  var(--primary),
                  oklch(0.348 0.106 295)
                );
              "
              :disabled="isSubmitting"
            >
              {{ isSubmitting ? "Хадгалж байна..." : "Шинэчлэх" }}
            </Button>
          </template>
        </JobForm>
      </div>
    </div>
  </template>

  <!-- ── View mode ──────────────────────────────────────────────────────── -->
  <template v-else>
    <div class="space-y-5">
      <!-- ── Recruiter header ── -->
      <div v-if="isRecruiter">
        <!-- <div class="mb-3 flex items-center gap-2 text-[13px] text-muted-foreground">
          <button class="flex items-center gap-1 hover:text-foreground transition" @click="router.push('/jobs')">
            <ChevronLeft class="h-4 w-4" /> Ажлын байрууд
          </button>
        </div> -->
        <div class="flex items-start justify-between gap-4">
          <div class="flex-1 min-w-0">
            <h1 class="text-[26px] font-semibold tracking-[-0.6px]">
              {{ job.title }}
            </h1>
            <div class="mt-2 flex flex-wrap items-center gap-2">
              <span
                class="flex items-center gap-1 text-[13px] text-muted-foreground"
              >
                <MapPin class="h-3.5 w-3.5 shrink-0" /> {{ job.location }}
              </span>
              <span
                class="inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground"
              >
                {{ job.type || job.employment_type }}
              </span>
              <span
                class="inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground"
              >
                {{ job.level || job.seniority }}
              </span>
              <span
                class="inline-flex items-center rounded-full px-2.5 py-0.5 text-[12px] font-semibold"
                :class="
                  job.status === 'posted'
                    ? 'badge-success'
                    : 'bg-muted text-muted-foreground'
                "
              >
                {{ jobStatusLabel(job.status) }}
              </span>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <Button
              variant="outline"
              class="h-full rounded-full"
              title="Горилогчид"
              @click="router.push(`/jobs/${job.id}/applicants`)"
            >
              <Users class="h-3.5 w-3.5" />
            </Button>
            <Button
              variant="outline"
              class="rounded-full"
              @click="isEditing = true"
            >
              <Pencil class="h-3.5 w-3.5" /> Засах
            </Button>
          </div>
        </div>
      </div>

      <!-- ── Applicant header ── -->
      <div
        v-else
        class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between"
      >
        <div class="flex-1 min-w-0">
          <h1 class="text-[26px] font-semibold tracking-[-0.6px]">
            {{ job.title }}
          </h1>
          <div class="mt-2 flex flex-wrap items-center gap-2">
            <span
              class="flex items-center gap-1 text-[13px] text-muted-foreground"
            >
              <MapPin class="h-3.5 w-3.5 shrink-0" /> {{ job.location }}
            </span>
            <span
              v-if="job.type || job.employment_type"
              class="inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground"
            >
              {{ job.type || job.employment_type }}
            </span>
            <span
              v-if="job.level || job.seniority"
              class="inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-[12px] text-muted-foreground"
            >
              {{ job.level || job.seniority }}
            </span>
          </div>
          <p
            v-if="job.min_salary || job.max_salary"
            class="mt-2 text-[15px] font-semibold"
          >
            {{ job.min_salary?.toLocaleString() }}₮ –
            {{ job.max_salary?.toLocaleString() }}₮
          </p>
        </div>

        <div class="flex flex-col gap-2 sm:flex-row sm:items-start">
          <template v-if="!application">
            <Button
              variant="outline"
              class="rounded-full"
              @click="analyzeOpen = true"
            >
              <Search class="h-4 w-4" /> Нийцэл шалгах
            </Button>
            <Button
              class="rounded-full text-white hover:opacity-90"
              style="
                background: linear-gradient(
                  135deg,
                  var(--primary),
                  oklch(0.348 0.106 295)
                );
              "
              @click="applyOpen = true"
            >
              <Upload class="h-4 w-4" /> CV илгээх
            </Button>
          </template>
          <template v-else-if="application.status === 'pending'">
            <div
              class="flex items-center gap-2 rounded-full border border-border px-4 py-2 text-[13.5px] text-muted-foreground"
            >
              <Loader2 class="h-4 w-4 animate-spin text-primary" /> Үнэлгээ
              хийгдэж байна...
            </div>
          </template>
          <template v-else-if="application?.status == 'assessed'">
            <div class="flex mt-1 flex-col gap-1 sm:items-end">
              <span
                class="inline-flex items-center rounded-full px-3 py-1.5 text-[12.5px] font-semibold badge-success"
              >
                Анкет илгээсэн
              </span>
            </div>
          </template>
        </div>
      </div>

      <!-- ── AI Evaluation Result ──────────────────────────────────────── -->
      <JobEvaluationResult
        v-if="displayResult"
        :result="displayResult"
        :mode="isPreview ? 'preview' : 'submitted'"
        :application-status="application?.status"
        @apply="applyOpen = true"
        @reload-cv="analyzeOpen = true"
      />

      <!-- ── Body: 2-col layout ─────────────────────────────────────────── -->
      <div class="grid gap-5 lg:grid-cols-3">
        <div class="space-y-5 lg:col-span-2">
          <!-- Description -->
          <div class="rounded-2xl border border-border bg-card p-5">
            <p class="mb-3 text-[15px] font-semibold">Тайлбар</p>
            <p
              class="whitespace-pre-line text-[13.5px] leading-[1.7] text-muted-foreground"
            >
              {{ job.additional_info || job.description }}
            </p>
          </div>

          <!-- Duties -->
          <div
            v-if="job.duties?.length"
            class="rounded-2xl border border-border bg-card p-5"
          >
            <p class="mb-3 text-[15px] font-semibold">Үүрэг хариуцлага</p>
            <div class="space-y-2">
              <div
                v-for="duty in job.duties"
                :key="duty"
                class="flex items-start gap-3 rounded-xl border p-3.5 transition"
                :class="
                  dutyAssessFor(duty)
                    ? assessRowBg(dutyAssessFor(duty)!.status)
                    : 'border-border bg-muted/20'
                "
              >
                <component
                  :is="
                    dutyAssessFor(duty)
                      ? assessIcon(dutyAssessFor(duty)!.status)
                      : ChevronRight
                  "
                  class="mt-0.5 h-4 w-4 shrink-0"
                  :class="
                    dutyAssessFor(duty)
                      ? assessColor(dutyAssessFor(duty)!.status)
                      : 'text-primary'
                  "
                />
                <div class="flex-1">
                  <p
                    class="text-[13.5px]"
                    :class="
                      dutyAssessFor(duty)
                        ? 'font-medium'
                        : 'text-muted-foreground'
                    "
                  >
                    {{ duty }}
                  </p>
                  <p
                    v-if="dutyAssessFor(duty)"
                    class="mt-0.5 text-[12px]"
                    :class="assessColor(dutyAssessFor(duty)!.status)"
                  >
                    {{ dutyAssessFor(duty)!.explanation }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Requirements -->
          <div
            v-if="job.requirements?.length"
            class="rounded-2xl border border-border bg-card p-5"
          >
            <p class="mb-3 text-[15px] font-semibold">Шаардлагууд</p>
            <div class="space-y-2">
              <div
                v-for="req in job.requirements"
                :key="req"
                class="flex items-start gap-3 rounded-xl border p-3.5 transition"
                :class="
                  reqAssessFor(req)
                    ? assessRowBg(reqAssessFor(req)!.status)
                    : 'border-border bg-muted/20'
                "
              >
                <component
                  :is="
                    reqAssessFor(req)
                      ? assessIcon(reqAssessFor(req)!.status)
                      : ChevronRight
                  "
                  class="mt-0.5 h-4 w-4 shrink-0"
                  :class="
                    reqAssessFor(req)
                      ? assessColor(reqAssessFor(req)!.status)
                      : 'text-primary'
                  "
                />
                <div class="flex-1">
                  <p
                    class="text-[13.5px]"
                    :class="
                      reqAssessFor(req)
                        ? 'font-medium'
                        : 'text-muted-foreground'
                    "
                  >
                    {{ req }}
                  </p>
                  <p
                    v-if="reqAssessFor(req)"
                    class="mt-0.5 text-[12px]"
                    :class="assessColor(reqAssessFor(req)!.status)"
                  >
                    {{ reqAssessFor(req)!.explanation }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- ── Right column ───────────────────────────────────────────── -->
        <div class="space-y-5">
          <!-- Recommendations -->
          <div
            v-if="displayResult?.recommendations?.length"
            class="rounded-2xl border p-5 ai-surface"
          >
            <div class="mb-3 flex items-center gap-2">
              <Sparkles class="h-4 w-4 text-primary" />
              <p class="text-[15px] font-semibold">Сайжруулах зөвлөмж</p>
            </div>
            <ul class="space-y-2">
              <li
                v-for="rec in displayResult.recommendations"
                :key="rec"
                class="flex items-start gap-2 text-[13.5px] text-muted-foreground"
              >
                <ChevronRight class="mt-0.5 h-4 w-4 shrink-0 text-primary" />
                {{ rec }}
              </li>
            </ul>
          </div>
          <!-- Bonuses -->
          <div
            v-if="job.bonuses?.length"
            class="rounded-2xl border border-border bg-card p-5"
          >
            <p class="mb-3 text-[15px] font-semibold">Нэмэлт давуу тал</p>
            <ul class="space-y-1.5">
              <li
                v-for="bonus in job.bonuses"
                :key="bonus"
                class="flex items-start gap-2 text-[13.5px] text-muted-foreground"
              >
                <span class="mt-0.5 text-success-foreground font-bold">✓</span>
                {{ bonus }}
              </li>
            </ul>
          </div>

          <!-- Contact -->
          <div
            v-if="job.contact_info"
            class="rounded-2xl border border-border bg-card p-5"
          >
            <p class="mb-2 text-[15px] font-semibold">Холбоо барих</p>
            <p class="text-[13.5px] text-muted-foreground">
              {{ job.contact_info }}
            </p>
          </div>

          <!-- Company profile URL -->
          <div
            v-if="job.company_profile_url"
            class="rounded-2xl border border-border bg-card p-5"
          >
            <p class="mb-2 text-[15px] font-semibold">Компанийн вэбсайт</p>
            <a
              :href="job.company_profile_url"
              target="_blank"
              rel="noopener noreferrer"
              class="text-[13.5px] text-primary underline-offset-2 hover:underline break-all"
            >
              {{ job.company_profile_url }}
            </a>
          </div>

          <!-- Location -->
          <div
            v-if="job.location"
            class="rounded-2xl border border-border bg-card p-5"
          >
            <div class="mb-3 flex items-center justify-between gap-2">
              <p class="text-[15px] font-semibold">Байршил</p>
              <label
                v-if="job.location_x && job.location_y"
                class="flex cursor-pointer items-center gap-1.5 text-[12px] text-muted-foreground"
              >
                <Checkbox v-model:checked="showJobMap" />
                Газрын зураг
              </label>
            </div>
            <div
              class="flex items-start gap-1.5 text-[13.5px] text-muted-foreground"
            >
              <MapPin class="mt-0.5 h-4 w-4 shrink-0" />
              <span>{{ job.location }}</span>
            </div>
            <LocationMap
              v-if="showJobMap && job.location_x && job.location_y"
              class="mt-3"
              :lat="job.location_x"
              :lng="job.location_y"
              :label="job.location"
            />
          </div>
        </div>
      </div>
    </div>
  </template>

  <!-- ── Apply dialog ────────────────────────────────────────────────────── -->
  <Dialog v-model:open="applyOpen">
    <DialogContent class="rounded-2xl sm:max-w-md">
      <DialogHeader>
        <DialogTitle>CV илгээх</DialogTitle>
        <DialogDescription>
          AI таны CV-г ажлын байрны шаардлагатай харьцуулж үнэлнэ.
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-2">
        <!-- Use saved profile option -->
        <div v-if="hasSavedCV" class="space-y-3">
          <Button
            type="button"
            variant="ghost"
            class="h-auto w-full justify-start gap-3 rounded-xl border p-3.5 text-left"
            :class="
              applyUseProfile
                ? 'border-primary bg-primary/5'
                : 'border-border hover:bg-muted/40'
            "
            @click="
              applyUseProfile = true;
              applyCvFile = null;
            "
          >
            <div
              class="flex h-4 w-4 shrink-0 items-center justify-center rounded-full border-2 transition"
              :class="
                applyUseProfile
                  ? 'border-primary bg-primary'
                  : 'border-muted-foreground'
              "
            >
              <div
                v-if="applyUseProfile"
                class="h-1.5 w-1.5 rounded-full bg-white"
              />
            </div>
            <div>
              <p class="text-[13.5px] font-medium">Хадгалагдсан CV ашиглах</p>
              <p class="text-[12px] text-muted-foreground">
                CV бүрдүүлэгч хэсэгт оруулсан мэдээлэл
              </p>
            </div>
          </Button>
          <Button
            type="button"
            variant="ghost"
            class="h-auto w-full justify-start gap-3 rounded-xl border p-3.5 text-left"
            :class="
              !applyUseProfile
                ? 'border-primary bg-primary/5'
                : 'border-border hover:bg-muted/40'
            "
            @click="applyUseProfile = false"
          >
            <div
              class="flex h-4 w-4 shrink-0 items-center justify-center rounded-full border-2 transition"
              :class="
                !applyUseProfile
                  ? 'border-primary bg-primary'
                  : 'border-muted-foreground'
              "
            >
              <div
                v-if="!applyUseProfile"
                class="h-1.5 w-1.5 rounded-full bg-white"
              />
            </div>
            <div>
              <p class="text-[13.5px] font-medium">PDF файл оруулах</p>
              <p class="text-[12px] text-muted-foreground">
                Өөрийн PDF CV файлыг upload хийх
              </p>
            </div>
          </Button>
        </div>

        <!-- PDF upload (shown when not using profile, or no saved profile) -->
        <div v-if="!applyUseProfile" class="space-y-2">
          <Label>CV файл (PDF)</Label>
          <div
            v-if="applyCvFile"
            class="flex items-center justify-between gap-3 rounded-xl border border-border bg-muted/30 px-3 py-2"
          >
            <p class="truncate text-[12.5px] text-muted-foreground">
              <span class="font-medium text-foreground">{{
                applyCvFile.name
              }}</span>
              <span class="ml-2 text-success-foreground"
                >— нийтлэл шалгасан CV</span
              >
            </p>
            <Button
              type="button"
              variant="link"
              class="h-auto shrink-0 p-0 text-[12px]"
              @click="applyCvFile = null"
            >
              Өөр файл
            </Button>
          </div>
          <Input v-else type="file" accept=".pdf" @change="onApplyFileChange" />
        </div>

        <p
          v-if="application"
          class="rounded-xl bg-muted/30 px-3 py-2 text-[12.5px] text-muted-foreground"
        >
          Өмнөх анкет шинэ CV-гээр солигдоно.
        </p>
      </div>
      <DialogFooter>
        <Button
          variant="outline"
          :disabled="isApplying"
          @click="applyOpen = false"
          >Болих</Button
        >
        <Button
          :disabled="isApplying || (!applyUseProfile && !applyCvFile)"
          @click="submitApplication"
        >
          <Loader2 v-if="isApplying" class="mr-2 h-4 w-4 animate-spin" />
          {{ isApplying ? "Илгээж байна..." : "Илгээх" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>

  <!-- ── Analyze dialog ──────────────────────────────────────────────────── -->
  <Dialog v-model:open="analyzeOpen">
    <DialogContent class="rounded-2xl sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Нийцэл шалгах</DialogTitle>
        <DialogDescription>
          Ажлын байранд хэр тохирохыг урьдчилан үнэлүүлнэ үү. Анкет илгээхгүй —
          үр дүн хуудсанд харагдана.
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-2">
        <!-- Use saved profile option -->
        <div v-if="hasSavedCV" class="space-y-3">
          <Button
            type="button"
            variant="ghost"
            class="h-auto w-full justify-start gap-3 rounded-xl border p-3.5 text-left"
            :class="
              analyzeUseProfile
                ? 'border-primary bg-primary/5'
                : 'border-border hover:bg-muted/40'
            "
            @click="
              analyzeUseProfile = true;
              analyzeCvFile = null;
            "
          >
            <div
              class="flex h-4 w-4 shrink-0 items-center justify-center rounded-full border-2 transition"
              :class="
                analyzeUseProfile
                  ? 'border-primary bg-primary'
                  : 'border-muted-foreground'
              "
            >
              <div
                v-if="analyzeUseProfile"
                class="h-1.5 w-1.5 rounded-full bg-white"
              />
            </div>
            <div>
              <p class="text-[13.5px] font-medium">Хадгалагдсан CV ашиглах</p>
              <p class="text-[12px] text-muted-foreground">
                CV бүрдүүлэгч хэсэгт оруулсан мэдээлэл
              </p>
            </div>
          </Button>
          <Button
            type="button"
            variant="ghost"
            class="h-auto w-full justify-start gap-3 rounded-xl border p-3.5 text-left"
            :class="
              !analyzeUseProfile
                ? 'border-primary bg-primary/5'
                : 'border-border hover:bg-muted/40'
            "
            @click="analyzeUseProfile = false"
          >
            <div
              class="flex h-4 w-4 shrink-0 items-center justify-center rounded-full border-2 transition"
              :class="
                !analyzeUseProfile
                  ? 'border-primary bg-primary'
                  : 'border-muted-foreground'
              "
            >
              <div
                v-if="!analyzeUseProfile"
                class="h-1.5 w-1.5 rounded-full bg-white"
              />
            </div>
            <div>
              <p class="text-[13.5px] font-medium">PDF файл оруулах</p>
              <p class="text-[12px] text-muted-foreground">
                Өөрийн PDF CV файлыг upload хийх
              </p>
            </div>
          </Button>
        </div>

        <!-- PDF upload -->
        <div v-if="!analyzeUseProfile" class="space-y-2">
          <Label>CV файл (PDF)</Label>
          <Input type="file" accept=".pdf" @change="onAnalyzeFileChange" />
          <p v-if="analyzeCvFile" class="text-[12.5px] text-muted-foreground">
            Сонгосон: {{ analyzeCvFile.name }}
          </p>
        </div>

        <div
          v-if="isAnalyzing"
          class="flex items-center gap-3 rounded-xl bg-muted/30 px-4 py-3 text-[13.5px] text-muted-foreground"
        >
          <Loader2 class="h-4 w-4 animate-spin text-primary" />
          AI үнэлгээ хийгдэж байна... Хэдэн секунд хүлээнэ үү.
        </div>
      </div>
      <DialogFooter>
        <Button
          variant="outline"
          :disabled="isAnalyzing"
          @click="analyzeOpen = false"
          >Болих</Button
        >
        <Button
          :disabled="isAnalyzing || (!analyzeUseProfile && !analyzeCvFile)"
          @click="runAnalyze"
        >
          <Loader2 v-if="isAnalyzing" class="mr-2 h-4 w-4 animate-spin" />
          {{ isAnalyzing ? "Үнэлж байна..." : "Үнэлгээ хийх" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
