<script setup lang="ts">
import type { Application, Task } from "../../../composables/types";
import { toast } from "vue-sonner";
import {
  ChevronLeft,
  Loader2,
  Users,
  CheckCircle2,
  XCircle,
  MinusCircle,
  ChevronRight,
  CalendarDays,
  ClipboardList,
  Sparkles,
  ArrowUpDown,
  BookOpen,
  PenLine,
  MapPin,
  Video,
} from "lucide-vue-next";

definePageMeta({ middleware: "auth" });

const route = useRoute();
const router = useRouter();
const applicationsAPI = useApplicationsAPI();
const tasksAPI = useTasksAPI();
const { user } = useAuth();

const jobID = computed(() => Number(route.params.id));
const isRecruiter = computed(() => user.value?.role === "recruiter");

const applications = ref<Application[]>([]);
const tasks = ref<Task[]>([]);
const loading = ref(true);
const updatingStatus = ref(false);

// ── Detail dialog ──────────────────────────────────────────────────────────
const detailOpen = ref(false);
const selected = ref<Application | null>(null);

function openDetail(app: Application) {
  selected.value = app;
  detailOpen.value = true;
}

// ── Interview dialog ───────────────────────────────────────────────────────
const interviewOpen = ref(false);
const interviewTarget = ref<Application | null>(null);
const interviewDate = ref("");
const interviewTime = ref("");
const interviewLocation = ref("");
const interviewMeetLink = ref("");
const interviewNote = ref("");
const isScheduling = ref(false);
type InterviewType = "onsite" | "google_meet";
const interviewType = ref<InterviewType>("onsite");

function isMeetUrl(val: string | null | undefined): boolean {
  return !!val && val.startsWith("https://meet.google.com");
}

const companyAPI = useCompanyAPI();
const integrationsAPI = useIntegrationsAPI();
const companyLocation = ref("");
const gcalConnected = ref(false);

onMounted(async () => {
  try {
    const [company, calStatus] = await Promise.allSettled([
      companyAPI.get(),
      integrationsAPI.googleCalendarStatus(),
    ]);
    if (company.status === "fulfilled") {
      const parts = [company.value.city, company.value.district].filter(
        Boolean,
      );
      companyLocation.value = parts.join(", ");
    }
    if (calStatus.status === "fulfilled") {
      gcalConnected.value = calStatus.value.connected;
    }
  } catch {
    // ignore
  }
});

function openInterviewDialog(app: Application) {
  interviewTarget.value = app;
  interviewDate.value = "";
  interviewTime.value = "";
  interviewNote.value = app.interview_note || "";
  if (isMeetUrl(app.interview_location)) {
    interviewType.value = "google_meet";
    interviewMeetLink.value = app.interview_location!;
    interviewLocation.value = "";
  } else {
    interviewType.value = "onsite";
    interviewLocation.value = app.interview_location || companyLocation.value;
    interviewMeetLink.value = "";
  }
  interviewOpen.value = true;
}

async function scheduleInterview() {
  if (!interviewDate.value || !interviewTime.value) {
    toast.error("Огноо болон цагаа оруулна уу");
    return;
  }
  isScheduling.value = true;
  try {
    const interviewAt = `${interviewDate.value}T${interviewTime.value}`;
    const isGoogleMeet = interviewType.value === "google_meet";
    const generateMeet = isGoogleMeet && gcalConnected.value;
    const locationValue = isGoogleMeet
      ? gcalConnected.value
        ? ""
        : interviewMeetLink.value
      : interviewLocation.value;
    const updated = await applicationsAPI.scheduleInterview(
      interviewTarget.value!.id,
      interviewAt,
      locationValue,
      interviewNote.value,
      generateMeet,
    );
    const idx = applications.value.findIndex((a) => a.id === updated.id);
    if (idx !== -1) applications.value[idx] = updated;
    if (selected.value?.id === updated.id) selected.value = updated;
    interviewOpen.value = false;
    toast.success(
      `Ярилцлага ${interviewTarget.value?.applicant_name} — ${formatDateTime(updated.interview_at!)} товлогдлоо.`,
    );
  } catch (e: any) {
    toast.error(e?.data?.message || "Ярилцлага товлоход алдаа гарлаа");
  } finally {
    isScheduling.value = false;
  }
}

// ── Task dialog ────────────────────────────────────────────────────────────
const taskOpen = ref(false);
const taskTarget = ref<Application | null>(null);
const taskTitle = ref("");
const taskDescription = ref("");
const taskDueDate = ref("");
const isGeneratingTask = ref(false);
const isSendingTask = ref(false);

// library mode: pick from pre-created tasks
type TaskMode = "library" | "new";
const taskMode = ref<TaskMode>("library");
const libraryTasks = ref<Task[]>([]);
const selectedLibraryTask = ref<Task | null>(null);

async function loadLibraryTasks() {
  libraryTasks.value = await tasksAPI.list(jobID.value).catch(() => []);
}

function openTaskDialog(app: Application) {
  taskTarget.value = app;
  taskTitle.value = "";
  taskDescription.value = "";
  taskDueDate.value = "";
  taskMode.value = "library";
  selectedLibraryTask.value = null;
  taskOpen.value = true;
  loadLibraryTasks();
}

function pickLibraryTask(task: Task) {
  selectedLibraryTask.value = task;
  taskTitle.value = task.title;
  taskDescription.value = task.description;
  taskDueDate.value = task.due_date ? task.due_date.split("T")[0] : "";
}

async function generateTask() {
  isGeneratingTask.value = true;
  try {
    const gen = await tasksAPI.generate(jobID.value);
    taskTitle.value = gen.title;
    taskDescription.value = gen.description;
    toast.success("AI даалгавар үүсгэлээ. Засварлаж болно.");
  } catch (e: any) {
    toast.error(e?.data?.message || "Даалгавар үүсгэхэд алдаа гарлаа");
  } finally {
    isGeneratingTask.value = false;
  }
}

async function sendTask() {
  isSendingTask.value = true;
  try {
    if (taskMode.value === "library") {
      const sent = await tasksAPI.send(
        selectedLibraryTask.value!.id,
        taskTarget.value?.id,
        taskDueDate.value || undefined,
      );
      tasks.value.push(sent);
      taskOpen.value = false;
      toast.success(
        `Даалгавар "${sent.title}" — ${taskTarget.value?.applicant_name}-д амжилттай илгээгдлээ.`,
      );
    } else {
      if (!taskTitle.value.trim() || !taskDescription.value.trim()) {
        toast.error("Гарчиг болон тайлбар шаардлагатай");
        return;
      }
      const task = await tasksAPI.create({
        job_posting_id: jobID.value,
        application_id: taskTarget.value?.id,
        title: taskTitle.value.trim(),
        description: taskDescription.value.trim(),
        due_date: taskDueDate.value || undefined,
      });
      await tasksAPI.send(task.id);
      tasks.value.push(task);
      taskOpen.value = false;
      toast.success(
        `Даалгавар "${task.title}" — ${taskTarget.value?.applicant_name}-д амжилттай илгээгдлээ.`,
      );
    }
  } catch (e: any) {
    toast.error(e?.data?.message || "Даалгавар илгээхэд алдаа гарлаа");
  } finally {
    isSendingTask.value = false;
  }
}

// ── Status actions ─────────────────────────────────────────────────────────
async function updateStatus(id: number, status: string) {
  updatingStatus.value = true;
  try {
    await applicationsAPI.updateStatus(id, status);
    await reload();
    if (selected.value?.id === id) {
      selected.value = applications.value.find((a) => a.id === id) ?? null;
    }
    const label = status === "shortlisted" ? "Ажилд авсан" : "Тэнцээгүй";
    toast.success(`Горилогчийн төлөв "${label}" болж шинэчлэгдлээ.`);
    if (status === "shortlisted" || status === "rejected") {
      detailOpen.value = false;
    }
  } catch {
    toast.error("Төлөв шинэчлэхэд алдаа гарлаа");
  } finally {
    updatingStatus.value = false;
  }
}

// ── Sort ───────────────────────────────────────────────────────────────────
type SortKey = "score_desc" | "score_asc" | "name" | "status" | "date";
const sortBy = ref<SortKey>("score_desc");

const sortOptions: { value: SortKey; label: string }[] = [
  { value: "score_desc", label: "Оноо (Их → Бага)" },
  { value: "score_asc", label: "Оноо (Бага → Их)" },
  { value: "name", label: "Нэрээр" },
  { value: "status", label: "Төлвөөр" },
  { value: "date", label: "Огноогоор" },
];

// ── Derived status ─────────────────────────────────────────────────────────
function taskForApp(appId: number): Task | null {
  return tasks.value.find((t) => t.application_id === appId) ?? null;
}

type DerivedStatus =
  | "applied"
  | "interview_scheduled"
  | "task_sent"
  | "task_completed"
  | "hired"
  | "failed";

function derivedStatus(app: Application): DerivedStatus {
  if (app.status === "shortlisted") return "hired";
  if (app.status === "rejected") return "failed";
  const task = taskForApp(app.id);
  if (task) {
    if (task.status === "completed" || task.status === "graded")
      return "task_completed";
    if (task.status === "sent") return "task_sent";
  }
  if (app.interview_at) return "interview_scheduled";
  return "applied";
}

const statusConfig: Record<DerivedStatus, { label: string; cls: string }> = {
  applied: { label: "Анкет илгээсэн", cls: "badge-info" },
  interview_scheduled: {
    label: "Ярилцлага товлосон",
    cls: "bg-primary/10 text-primary",
  },
  task_sent: { label: "Даалгавар илгээсэн", cls: "badge-warning" },
  task_completed: {
    label: "Даалгавар гүйцэтгэсэн",
    cls: "bg-success-bg text-success-foreground border border-success/20",
  },
  hired: { label: "Ажилд авсан", cls: "badge-success" },
  failed: { label: "Тэнцээгүй", cls: "bg-destructive/10 text-destructive" },
};

const statusOrder: Record<DerivedStatus, number> = {
  hired: 0,
  task_completed: 1,
  task_sent: 2,
  interview_scheduled: 3,
  applied: 4,
  failed: 5,
};

const applicantPage = ref(1);
const APPLICANT_PAGE_SIZE = 10;
const applicantTotalPages = computed(() =>
  Math.max(1, Math.ceil(sortedApplications.value.length / APPLICANT_PAGE_SIZE)),
);
const paginatedApplicants = computed(() => {
  const start = (applicantPage.value - 1) * APPLICANT_PAGE_SIZE;
  return sortedApplications.value.slice(start, start + APPLICANT_PAGE_SIZE);
});

watch(sortBy, () => {
  applicantPage.value = 1;
});

const sortedApplications = computed(() => {
  const arr = [...applications.value];
  switch (sortBy.value) {
    case "score_desc":
      return arr.sort((a, b) => b.overall_score - a.overall_score);
    case "score_asc":
      return arr.sort((a, b) => a.overall_score - b.overall_score);
    case "name":
      return arr.sort((a, b) =>
        (a.applicant_name || "").localeCompare(b.applicant_name || ""),
      );
    case "status":
      return arr.sort(
        (a, b) => statusOrder[derivedStatus(a)] - statusOrder[derivedStatus(b)],
      );
    case "date":
      return arr.sort(
        (a, b) =>
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
      );
    default:
      return arr;
  }
});

// ── Load ───────────────────────────────────────────────────────────────────
async function reload() {
  const [apps, jobTasks] = await Promise.all([
    applicationsAPI.listForJob(jobID.value),
    tasksAPI.list(jobID.value),
  ]);
  applications.value = apps;
  tasks.value = jobTasks;
}

if (!isRecruiter.value) router.replace("/jobs");

loading.value = true;
try {
  await reload();
} catch {
  toast.error("Мэдээлэл ачааллаж чадсангүй");
} finally {
  loading.value = false;
}

// ── Helpers ────────────────────────────────────────────────────────────────
function scoreColor(score: number) {
  if (score >= 75) return "text-success-foreground";
  if (score >= 50) return "text-warning-foreground";
  return "text-destructive";
}

function scoreBorderBg(score: number) {
  if (score >= 75) return "border-success bg-success-bg";
  if (score >= 50) return "border-warning bg-warning-bg";
  return "border-destructive/40 bg-destructive/5";
}

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

function formatDate(d: string | null | undefined) {
  if (!d) return "—";
  return new Intl.DateTimeFormat("mn-MN", { dateStyle: "medium" }).format(
    new Date(d),
  );
}

function formatDateTime(d: string | null | undefined) {
  if (!d) return "—";
  return new Intl.DateTimeFormat("mn-MN", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(d));
}
</script>

<template>
  <div class="space-y-5">
    <!-- Page header -->
    <div class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <div>
          <p
            class="text-[12px] font-semibold uppercase tracking-[0.8px] text-muted-foreground"
          >
            Горилогчид
          </p>
          <h1 class="text-[22px] font-semibold tracking-[-0.5px]">
            {{ loading ? "…" : `${applications.length} горилогч` }}
          </h1>
        </div>
      </div>

      <!-- Sort -->
      <div class="flex items-center gap-2">
        <ArrowUpDown class="h-4 w-4 shrink-0 text-muted-foreground" />
        <Select v-model="sortBy">
          <SelectTrigger class="w-48 rounded-xl">
            <SelectValue placeholder="Эрэмбэлэх" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem
              v-for="opt in sortOptions"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="h-7 w-7 animate-spin text-primary" />
    </div>

    <!-- Empty -->
    <div
      v-else-if="!applications.length"
      class="flex flex-col items-center py-20 text-center"
    >
      <div
        class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-muted"
      >
        <Users class="h-7 w-7 text-muted-foreground" />
      </div>
      <p class="text-[16px] font-semibold">Горилогч байхгүй байна</p>
      <p class="mt-2 text-[13.5px] text-muted-foreground">
        Горилогчид анкет илгээх үед энд харагдана.
      </p>
    </div>

    <!-- Applicant cards -->
    <div v-else class="space-y-3">
      <div
        v-for="app in paginatedApplicants"
        :key="app.id"
        class="group cursor-pointer rounded-2xl border border-border bg-card p-5 transition-all hover:border-primary/30 hover:shadow-sm"
        @click="openDetail(app)"
      >
        <div class="flex items-center gap-4">
          <!-- Score circle -->
          <div
            class="flex h-14 w-14 shrink-0 items-center justify-center rounded-full border-[3px]"
            :class="scoreBorderBg(app.overall_score)"
          >
            <span
              class="text-[13px] font-bold leading-none"
              :class="scoreColor(app.overall_score)"
            >
              {{ app.overall_score }}<span class="text-[9px]">%</span>
            </span>
          </div>

          <!-- Applicant info -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <p class="text-[15px] font-semibold">
                {{ app.applicant_name || "Хэрэглэгч" }}
              </p>
              <span
                class="inline-flex items-center rounded-full px-2.5 py-0.5 text-[11.5px] font-semibold"
                :class="statusConfig[derivedStatus(app)].cls"
              >
                {{ statusConfig[derivedStatus(app)].label }}
              </span>
            </div>
            <p class="mt-0.5 text-[12.5px] text-muted-foreground">
              {{ app.applicant_email }}
            </p>

            <!-- Skill chips preview -->
            <!-- <div
              v-if="app.matched_skills?.length"
              class="mt-2 flex flex-wrap gap-1"
            >
              <span
                v-for="s in app.matched_skills.slice(0, 4)"
                :key="s.skill"
                class="inline-flex items-center gap-0.5 rounded-full px-2 py-0.5 text-[11px] font-medium badge-success"
              >
                ✓ {{ s.skill }}
              </span>
              <span
                v-for="s in app.missing_skills?.slice(0, 2)"
                :key="'m-' + s.skill"
                class="inline-flex items-center gap-0.5 rounded-full px-2 py-0.5 text-[11px] font-medium bg-destructive/10 text-destructive"
              >
                ✗ {{ s.skill }}
              </span>
            </div> -->
          </div>

          <!-- Date + actions -->
          <div class="shrink-0 flex flex-col items-end gap-2" @click.stop>
            <p class="text-[12px] text-muted-foreground">
              {{ formatDate(app.created_at) }}
            </p>
            <div class="flex items-center gap-1.5">
              <button
                class="flex h-8 items-center gap-1.5 rounded-lg border border-border px-2.5 text-[12px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground"
                :title="
                  app.interview_at ? 'Ярилцлага товлогдсон' : 'Ярилцлага товлох'
                "
                @click="openInterviewDialog(app)"
              >
                <CalendarDays class="h-3.5 w-3.5" />
                <span class="hidden sm:inline">Ярилцлага</span>
              </button>
              <button
                class="flex h-8 items-center gap-1.5 rounded-lg border border-border px-2.5 text-[12px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground disabled:opacity-40"
                :disabled="!!taskForApp(app.id)"
                :title="
                  taskForApp(app.id)
                    ? 'Даалгавар илгээгдсэн'
                    : 'Даалгавар илгээх'
                "
                @click="openTaskDialog(app)"
              >
                <ClipboardList class="h-3.5 w-3.5" />
                <span class="hidden sm:inline">Даалгавар</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Applicants pagination -->
      <div class="flex items-center justify-center gap-1 pt-2">
        <button
          class="flex h-8 w-8 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
          :disabled="applicantPage === 1"
          @click="applicantPage--"
        >
          <ChevronLeft class="h-4 w-4" />
        </button>
        <button
          v-for="p in applicantTotalPages"
          :key="p"
          :class="[
            'h-8 w-8 rounded-full text-[13px] font-medium transition',
            p === applicantPage
              ? 'bg-primary text-background'
              : 'text-muted-foreground hover:bg-muted',
          ]"
          @click="applicantPage = p"
        >
          {{ p }}
        </button>
        <button
          class="flex h-8 w-8 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
          :disabled="applicantPage === applicantTotalPages"
          @click="applicantPage++"
        >
          <ChevronRight class="h-4 w-4" />
        </button>
      </div>
    </div>
  </div>

  <!-- ── Detail dialog ──────────────────────────────────────────────────── -->
  <Dialog v-model:open="detailOpen">
    <DialogContent
      class="rounded-2xl sm:max-w-2xl max-h-[90vh] overflow-y-auto"
    >
      <DialogHeader>
        <DialogTitle class="text-[18px]">{{
          selected?.applicant_name || "Горилогч"
        }}</DialogTitle>
        <DialogDescription>{{ selected?.applicant_email }}</DialogDescription>
      </DialogHeader>

      <div v-if="selected" class="space-y-5 py-2">
        <!-- Score + summary -->
        <div
          class="flex items-start gap-4 rounded-xl border p-4"
          :class="
            selected.overall_score >= 75
              ? 'bg-success-bg border-success/20'
              : selected.overall_score >= 50
                ? 'bg-warning-bg border-warning/20'
                : 'bg-destructive/5 border-destructive/15'
          "
        >
          <div
            class="flex h-16 w-16 shrink-0 items-center justify-center rounded-full border-[3px]"
            :class="scoreBorderBg(selected.overall_score)"
          >
            <div class="text-center">
              <p
                class="text-[18px] font-bold leading-none"
                :class="scoreColor(selected.overall_score)"
              >
                {{ selected.overall_score }}
              </p>
              <p
                class="text-[10px]"
                :class="scoreColor(selected.overall_score)"
              >
                /100
              </p>
            </div>
          </div>
          <div class="flex-1">
            <div class="mb-1.5 flex items-center gap-2">
              <Sparkles class="h-3.5 w-3.5 text-primary" />
              <p class="text-[13px] font-semibold">AI үнэлгээний дүгнэлт</p>
            </div>
            <p class="text-[13px] leading-[1.6] text-muted-foreground">
              {{ selected.summary }}
            </p>
          </div>
        </div>

        <!-- Status + interview -->
        <div class="grid gap-2.5 sm:grid-cols-2">
          <div class="rounded-xl border border-border bg-muted/20 px-4 py-3">
            <p
              class="text-[11px] font-semibold uppercase tracking-[0.5px] text-muted-foreground mb-1.5"
            >
              Одоогийн төлөв
            </p>
            <span
              class="inline-flex items-center rounded-full px-3 py-1 text-[12px] font-semibold"
              :class="statusConfig[derivedStatus(selected)].cls"
            >
              {{ statusConfig[derivedStatus(selected)].label }}
            </span>
          </div>
          <div class="rounded-xl border border-border bg-muted/20 px-4 py-3">
            <p
              class="text-[11px] font-semibold uppercase tracking-[0.5px] text-muted-foreground mb-1.5"
            >
              Ярилцлагийн тов
            </p>
            <p class="text-[13.5px] font-medium">
              {{
                selected.interview_at
                  ? formatDateTime(selected.interview_at)
                  : "Товлоогүй"
              }}
            </p>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap gap-2">
          <button
            v-if="selected.status !== 'shortlisted'"
            class="inline-flex items-center gap-2 rounded-full px-4 py-2 text-[13.5px] font-semibold text-white transition hover:opacity-90"
            style="
              background: linear-gradient(
                135deg,
                var(--primary),
                oklch(0.348 0.106 295)
              );
            "
            :disabled="updatingStatus"
            @click="updateStatus(selected.id, 'shortlisted')"
          >
            Ажилд авах
          </button>
          <button
            v-if="selected.status !== 'rejected'"
            class="inline-flex items-center gap-2 rounded-full border border-border px-4 py-2 text-[13.5px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground"
            :disabled="updatingStatus"
            @click="updateStatus(selected.id, 'rejected')"
          >
            Тэнцээгүй
          </button>
          <button
            class="inline-flex items-center gap-2 rounded-full border border-border px-4 py-2 text-[13.5px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground"
            @click="
              () => {
                detailOpen = false;
                openInterviewDialog(selected!);
              }
            "
          >
            <CalendarDays class="h-4 w-4" />
            Ярилцлага товлох
          </button>
          <button
            class="inline-flex items-center gap-2 rounded-full border border-border px-4 py-2 text-[13.5px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground disabled:opacity-40"
            :disabled="!!taskForApp(selected.id)"
            @click="
              () => {
                detailOpen = false;
                openTaskDialog(selected!);
              }
            "
          >
            <ClipboardList class="h-4 w-4" />
            Даалгавар илгээх
          </button>
        </div>

        <!-- Skills -->
        <div class="grid gap-4 sm:grid-cols-2">
          <div class="space-y-2">
            <p
              class="text-[11px] font-semibold uppercase tracking-[0.6px] text-success-foreground"
            >
              Тохирсон ур чадварууд
            </p>
            <div
              v-for="s in selected.matched_skills"
              :key="s.skill"
              class="rounded-xl border border-success/20 bg-success-bg p-3"
            >
              <p
                class="flex items-center gap-1.5 text-[12.5px] font-semibold text-success-foreground"
              >
                <CheckCircle2 class="h-3.5 w-3.5 shrink-0" /> {{ s.skill }}
              </p>
              <p class="mt-0.5 text-[12px] text-muted-foreground">
                {{ s.explanation }}
              </p>
            </div>
            <p
              v-if="!selected.matched_skills?.length"
              class="text-[12.5px] text-muted-foreground"
            >
              Байхгүй
            </p>
          </div>
          <div class="space-y-2">
            <p
              class="text-[11px] font-semibold uppercase tracking-[0.6px] text-destructive"
            >
              Дутуу ур чадварууд
            </p>
            <div
              v-for="s in selected.missing_skills"
              :key="s.skill"
              class="rounded-xl border border-destructive/15 bg-destructive/5 p-3"
            >
              <p
                class="flex items-center gap-1.5 text-[12.5px] font-semibold text-destructive"
              >
                <XCircle class="h-3.5 w-3.5 shrink-0" /> {{ s.skill }}
              </p>
              <p class="mt-0.5 text-[12px] text-muted-foreground">
                {{ s.explanation }}
              </p>
            </div>
            <p
              v-if="!selected.missing_skills?.length"
              class="text-[12.5px] text-muted-foreground"
            >
              Байхгүй
            </p>
          </div>
        </div>

        <!-- Duty assessments -->
        <div v-if="selected.duty_assessments?.length" class="space-y-2">
          <p class="text-[13.5px] font-semibold">Үүрэг хариуцлагын үнэлгээ</p>
          <div
            v-for="d in selected.duty_assessments"
            :key="d.duty"
            class="flex items-start gap-3 rounded-xl border p-3"
            :class="assessRowBg(d.status)"
          >
            <component
              :is="assessIcon(d.status)"
              class="mt-0.5 h-4 w-4 shrink-0"
              :class="assessColor(d.status)"
            />
            <div>
              <p class="text-[13px] font-medium">{{ d.duty }}</p>
              <p class="mt-0.5 text-[12px]" :class="assessColor(d.status)">
                {{ d.explanation }}
              </p>
            </div>
          </div>
        </div>

        <!-- Requirement assessments -->
        <div v-if="selected.requirement_assessments?.length" class="space-y-2">
          <p class="text-[13.5px] font-semibold">Шаардлагын үнэлгээ</p>
          <div
            v-for="r in selected.requirement_assessments"
            :key="r.requirement"
            class="flex items-start gap-3 rounded-xl border p-3"
            :class="assessRowBg(r.status)"
          >
            <component
              :is="assessIcon(r.status)"
              class="mt-0.5 h-4 w-4 shrink-0"
              :class="assessColor(r.status)"
            />
            <div>
              <p class="text-[13px] font-medium">{{ r.requirement }}</p>
              <p class="mt-0.5 text-[12px]" :class="assessColor(r.status)">
                {{ r.explanation }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </DialogContent>
  </Dialog>

  <!-- ── Interview dialog ───────────────────────────────────────────────── -->
  <Dialog v-model:open="interviewOpen">
    <DialogContent class="rounded-2xl sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Ярилцлага товлох</DialogTitle>
        <DialogDescription>
          {{ interviewTarget?.applicant_name }}-д ярилцлагын мэдээллийг оруулна
          уу.
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-2">
        <!-- Meeting type tabs -->
        <div
          class="flex rounded-xl border border-border overflow-hidden text-[13px] font-medium"
        >
          <button
            class="flex flex-1 items-center justify-center gap-1.5 py-2 transition"
            :class="
              interviewType === 'onsite'
                ? 'bg-primary text-white'
                : 'text-muted-foreground hover:bg-muted'
            "
            @click="interviewType = 'onsite'"
          >
            <MapPin class="h-3.5 w-3.5" /> Биечлэн уулзах
          </button>
          <button
            class="flex flex-1 items-center justify-center gap-1.5 py-2 transition"
            :class="
              interviewType === 'google_meet'
                ? 'bg-primary text-white'
                : 'text-muted-foreground hover:bg-muted'
            "
            @click="interviewType = 'google_meet'"
          >
            <Video class="h-3.5 w-3.5" /> Google Meet
          </button>
        </div>

        <!-- Date & Time -->
        <div class="grid gap-4 sm:grid-cols-2">
          <div class="space-y-2">
            <Label>Огноо</Label>
            <Input v-model="interviewDate" type="date" />
          </div>
          <div class="space-y-2">
            <Label>Цаг</Label>
            <Input v-model="interviewTime" type="time" />
          </div>
        </div>

        <!-- Onsite: location -->
        <div v-if="interviewType === 'onsite'" class="space-y-2">
          <Label>Байршил</Label>
          <Input v-model="interviewLocation" placeholder="Уулзах газар" />
        </div>

        <!-- Google Meet: auto-generate (if gcal connected) or manual link -->
        <div v-else class="space-y-2">
          <template v-if="gcalConnected">
            <div
              class="flex items-start gap-2.5 rounded-xl border border-primary/20 bg-primary/5 px-3.5 py-3 text-[12.5px]"
            >
              <Video class="mt-0.5 h-3.5 w-3.5 shrink-0 text-primary" />
              <span class="text-muted-foreground"
                >Google Calendar-тай холбогдсон тул Meet холбоосыг
                <span class="font-semibold text-primary"
                  >автоматаар үүсгэнэ</span
                >.
              </span>
            </div>
          </template>
          <template v-else>
            <Label>Google Meet холбоос</Label>
            <Input
              v-model="interviewMeetLink"
              placeholder="https://meet.google.com/xxx-xxxx-xxx"
              type="url"
            />
            <p class="text-[11.5px] text-muted-foreground">
              Google Calendar холбогдоогүй тул холбоосыг гараар оруулна уу.
            </p>
          </template>
        </div>

        <div class="space-y-2">
          <Label
            >Нэмэлт тэмдэглэл
            <span class="text-muted-foreground text-[11.5px]"
              >(заавал биш)</span
            ></Label
          >
          <Textarea
            v-model="interviewNote"
            rows="3"
            placeholder="Горилогчид дамжуулах нэмэлт мэдээлэл..."
          />
        </div>
      </div>
      <DialogFooter>
        <Button
          variant="outline"
          :disabled="isScheduling"
          @click="interviewOpen = false"
          >Болих</Button
        >
        <Button :disabled="isScheduling" @click="scheduleInterview">
          <Loader2 v-if="isScheduling" class="mr-2 h-4 w-4 animate-spin" />
          <CalendarDays v-else class="mr-2 h-4 w-4" />
          {{ isScheduling ? "Товлож байна..." : "Товлох" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>

  <!-- ── Task dialog ────────────────────────────────────────────────────── -->
  <Dialog v-model:open="taskOpen">
    <DialogContent class="rounded-2xl sm:max-w-lg flex flex-col max-h-[90vh]">
      <DialogHeader class="shrink-0">
        <DialogTitle>Даалгавар илгээх</DialogTitle>
      </DialogHeader>

      <!-- Mode toggle -->
      <div
        class="shrink-0 flex rounded-xl border border-border overflow-hidden text-[13px] font-medium"
      >
        <button
          class="flex flex-1 items-center justify-center gap-1.5 py-2 transition"
          :class="
            taskMode === 'library'
              ? 'bg-primary text-white'
              : 'text-muted-foreground hover:bg-muted'
          "
          @click="taskMode = 'library'"
        >
          <BookOpen class="h-3.5 w-3.5" /> Сангаас сонгох
        </button>
        <button
          class="flex flex-1 items-center justify-center gap-1.5 py-2 transition"
          :class="
            taskMode === 'new'
              ? 'bg-primary text-white'
              : 'text-muted-foreground hover:bg-muted'
          "
          @click="taskMode = 'new'"
        >
          <PenLine class="h-3.5 w-3.5" /> Шинэ үүсгэх
        </button>
      </div>

      <!-- Scrollable content area -->
      <div class="flex-1 overflow-y-auto min-h-0 px-0.5">
        <!-- Library mode -->
        <div v-if="taskMode === 'library'" class="space-y-2 py-1">
          <div
            v-if="!libraryTasks.length"
            class="flex flex-col items-center py-8 text-center"
          >
            <ClipboardList class="h-8 w-8 text-muted-foreground mb-2" />
            <p class="text-[13.5px] text-muted-foreground">
              Санд даалгавар байхгүй байна.
            </p>
            <button
              class="mt-2 text-[13px] text-primary underline-offset-2 hover:underline"
              @click="taskMode = 'new'"
            >
              Шинэ даалгавар үүсгэх
            </button>
          </div>
          <div
            v-for="lt in libraryTasks"
            :key="lt.id"
            class="cursor-pointer rounded-xl border p-3.5 transition"
            :class="
              selectedLibraryTask?.id === lt.id
                ? 'border-primary bg-primary/5'
                : 'border-border hover:border-primary/40 hover:bg-muted/30'
            "
            @click="pickLibraryTask(lt)"
          >
            <div class="flex items-start justify-between gap-2">
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-1.5">
                  <p class="text-[13.5px] font-semibold truncate">
                    {{ lt.title }}
                  </p>
                  <span
                    v-if="lt.created_by_ai"
                    class="inline-flex shrink-0 items-center gap-0.5 rounded-full px-1.5 py-0.5 text-[10px] font-semibold bg-primary/10 text-primary"
                  >
                    <Sparkles class="h-2 w-2" /> AI
                  </span>
                </div>
                <p
                  class="mt-0.5 text-[12.5px] text-muted-foreground line-clamp-2"
                >
                  {{ lt.description }}
                </p>
                <p
                  v-if="lt.duration_days"
                  class="mt-1 text-[11.5px] text-muted-foreground"
                >
                  Хугацаа: {{ lt.duration_days }} өдөр
                </p>
              </div>
              <ChevronRight
                class="h-4 w-4 shrink-0 text-muted-foreground mt-0.5"
              />
            </div>
          </div>

          <!-- Due date override when a task is selected -->
          <div v-if="selectedLibraryTask" class="space-y-1.5 pt-2">
            <Label class="text-[12.5px]">
              Дуусах огноо
              <span class="text-muted-foreground text-[11.5px]"
                >(заавал биш)</span
              >
            </Label>
            <Input v-model="taskDueDate" type="date" />
          </div>
        </div>

        <!-- New / edit mode -->
        <div v-else class="space-y-4 py-1">
          <button
            type="button"
            class="flex w-full items-center justify-center gap-2 rounded-xl border border-border py-2.5 text-[13.5px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground disabled:opacity-50"
            :disabled="isGeneratingTask"
            @click="generateTask"
          >
            <Loader2 v-if="isGeneratingTask" class="h-4 w-4 animate-spin" />
            <Sparkles v-else class="h-4 w-4 text-primary" />
            {{ isGeneratingTask ? "AI үүсгэж байна..." : "AI-аар үүсгэх" }}
          </button>
          <div class="space-y-2">
            <Label>Гарчиг</Label>
            <Input v-model="taskTitle" placeholder="Даалгаврын гарчиг" />
          </div>
          <div class="space-y-2">
            <Label>Тайлбар</Label>
            <Textarea
              v-model="taskDescription"
              rows="5"
              placeholder="Даалгаврын дэлгэрэнгүй тайлбар..."
            />
          </div>
          <div class="space-y-2">
            <Label> Дуусах огноо </Label>
            <Input v-model="taskDueDate" type="date" />
          </div>
        </div>
      </div>

      <DialogFooter class="shrink-0">
        <Button
          variant="outline"
          :disabled="isSendingTask"
          @click="taskOpen = false"
          >Болих</Button
        >
        <Button
          v-if="taskMode === 'library'"
          :disabled="isSendingTask || !selectedLibraryTask"
          @click="sendTask"
        >
          <Loader2 v-if="isSendingTask" class="mr-2 h-4 w-4 animate-spin" />
          <ClipboardList v-else class="mr-2 h-4 w-4" />
          {{ isSendingTask ? "Илгээж байна..." : "Илгээх" }}
        </Button>
        <Button
          v-else
          :disabled="isSendingTask || !taskTitle || !taskDescription"
          @click="sendTask"
        >
          <Loader2 v-if="isSendingTask" class="mr-2 h-4 w-4 animate-spin" />
          <ClipboardList v-else class="mr-2 h-4 w-4" />
          {{ isSendingTask ? "Илгээж байна..." : "Илгээх" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
