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
  Search,
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

// ── Selected applicant (inline panel) ─────────────────────────────────────
const selected = ref<Application | null>(null);

function selectApplicant(app: Application) {
  selected.value = app;
}

// ── Search & filter ────────────────────────────────────────────────────────
const searchQuery = ref("");
type FilterKey = "all" | "new" | "shortlisted" | "interview";
const activeFilter = ref<FilterKey>("all");

const filterOptions: { key: FilterKey; label: string }[] = [
  { key: "all", label: "Бүгд" },
  { key: "new", label: "Шинэ" },
  { key: "shortlisted", label: "Ажилд авсан" },
  { key: "interview", label: "Ярилцлага" },
];

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
  taskDueDate.value = task.due_date ? (task.due_date.split("T")[0] ?? "") : "";
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

// ── Filter logic ───────────────────────────────────────────────────────────
const filterCounts = computed(() => ({
  all: applications.value.length,
  new: applications.value.filter((a) => derivedStatus(a) === "applied").length,
  shortlisted: applications.value.filter((a) => derivedStatus(a) === "hired")
    .length,
  interview: applications.value.filter(
    (a) => derivedStatus(a) === "interview_scheduled",
  ).length,
}));

const filteredApplications = computed(() => {
  let arr = [...applications.value];

  // status filter
  if (activeFilter.value === "new")
    arr = arr.filter((a) => derivedStatus(a) === "applied");
  else if (activeFilter.value === "shortlisted")
    arr = arr.filter((a) => derivedStatus(a) === "hired");
  else if (activeFilter.value === "interview")
    arr = arr.filter((a) => derivedStatus(a) === "interview_scheduled");

  // search filter
  const q = searchQuery.value.trim().toLowerCase();
  if (q)
    arr = arr.filter(
      (a) =>
        (a.applicant_name || "").toLowerCase().includes(q) ||
        (a.applicant_email || "").toLowerCase().includes(q),
    );

  return arr;
});

// ── Pagination ─────────────────────────────────────────────────────────────
const applicantPage = ref(1);
const APPLICANT_PAGE_SIZE = 15;
const applicantTotalPages = computed(() =>
  Math.max(1, Math.ceil(sortedApplications.value.length / APPLICANT_PAGE_SIZE)),
);
const paginatedApplicants = computed(() => {
  const start = (applicantPage.value - 1) * APPLICANT_PAGE_SIZE;
  return sortedApplications.value.slice(start, start + APPLICANT_PAGE_SIZE);
});

watch([sortBy, searchQuery, activeFilter], () => {
  applicantPage.value = 1;
});

const sortedApplications = computed(() => {
  const arr = [...filteredApplications.value];
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
  if (applications.value.length) selected.value = applications.value[0] ?? null;
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

function dutyBarWidth(s: string): number {
  if (s === "met") return 90;
  if (s === "partial") return 55;
  return 20;
}

function dutyBarColor(s: string) {
  if (s === "met") return "bg-success";
  if (s === "partial") return "bg-warning";
  return "bg-destructive/70";
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

function initials(name: string | null | undefined): string {
  if (!name) return "?";
  const parts = name.trim().split(/\s+/);
  if (parts.length >= 2)
    return ((parts[0]![0] ?? "") + (parts[1]![0] ?? "")).toUpperCase();
  return name.slice(0, 2).toUpperCase();
}

function avatarGradient(id: number): string {
  const gradients = [
    "from-violet-500 to-purple-600",
    "from-blue-500 to-cyan-500",
    "from-emerald-500 to-teal-600",
    "from-orange-500 to-amber-500",
    "from-pink-500 to-rose-500",
    "from-indigo-500 to-blue-600",
  ];
  return gradients[id % gradients.length]!;
}
</script>

<template>
  <!-- Split-panel layout -->
  <div class="flex gap-4 min-h-0" style="height: calc(100vh - 48px)">
    <!-- ── Left panel: list ──────────────────────────────────────────────── -->
    <div
      :class="[
        'flex flex-col gap-3 lg:w-[340px] lg:shrink-0 min-h-0',
        selected ? 'hidden lg:flex' : 'flex w-full',
      ]"
    >
      <!-- Header row -->
      <div class="flex items-center justify-between gap-2 shrink-0">
        <div>
          <p
            class="text-[11px] font-semibold uppercase tracking-[0.8px] text-muted-foreground"
          >
            Горилогчид
          </p>
          <h1 class="text-[20px] font-semibold tracking-[-0.5px]">
            {{ loading ? "…" : `${applications.length} горилогч` }}
          </h1>
        </div>
        <div class="flex items-center gap-1.5">
          <ArrowUpDown class="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
          <Select v-model="sortBy">
            <SelectTrigger class="w-40 rounded-xl h-8 text-[12px]">
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

      <!-- Search -->
      <div class="relative shrink-0">
        <Search
          class="absolute left-3 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-muted-foreground"
        />
        <Input
          v-model="searchQuery"
          placeholder="Горилогч хайх…"
          class="pl-8 rounded-xl h-9 text-[13px]"
        />
      </div>

      <!-- Filter pills -->
      <div class="flex flex-wrap gap-1.5 shrink-0">
        <button
          v-for="f in filterOptions"
          :key="f.key"
          class="inline-flex items-center gap-1 rounded-full px-3 py-1 text-[12px] font-medium transition"
          :class="
            activeFilter === f.key
              ? 'bg-primary text-primary-foreground'
              : 'border border-border text-muted-foreground hover:bg-muted hover:text-foreground'
          "
          @click="activeFilter = f.key"
        >
          {{ f.label }}
          <span
            class="text-[10.5px] font-semibold"
            :class="activeFilter === f.key ? 'opacity-80' : 'opacity-60'"
          >
            {{ filterCounts[f.key] }}
          </span>
        </button>
      </div>

      <!-- Loading state -->
      <div v-if="loading" class="flex justify-center py-12">
        <Loader2 class="h-6 w-6 animate-spin text-primary" />
      </div>

      <!-- Empty list -->
      <div
        v-else-if="!paginatedApplicants.length"
        class="flex flex-col items-center py-12 text-center"
      >
        <div
          class="mb-3 flex h-12 w-12 items-center justify-center rounded-2xl bg-muted"
        >
          <Users class="h-5 w-5 text-muted-foreground" />
        </div>
        <p class="text-[14px] font-semibold">Горилогч байхгүй</p>
        <p class="mt-1 text-[12.5px] text-muted-foreground">
          {{
            searchQuery || activeFilter !== "all"
              ? "Шүүлтийг өөрчилнө үү"
              : "Горилогчид анкет илгээх үед энд харагдана."
          }}
        </p>
      </div>

      <!-- Applicant list -->
      <div v-else class="flex flex-col gap-2 overflow-y-auto min-h-0 flex-1">
        <div
          v-for="app in paginatedApplicants"
          :key="app.id"
          class="group cursor-pointer rounded-2xl border p-3 transition-all"
          :class="
            selected?.id === app.id
              ? 'border-primary/40 bg-primary/5'
              : 'border-border bg-card hover:border-primary/20 hover:bg-muted/30'
          "
          @click="selectApplicant(app)"
        >
          <div class="flex items-center gap-3">
            <!-- Avatar -->
            <div
              class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-gradient-to-br text-white text-[12px] font-bold"
              :class="avatarGradient(app.applicant_id)"
            >
              {{ initials(app.applicant_name) }}
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
              <p class="text-[13.5px] font-semibold truncate">
                {{ app.applicant_name || "Хэрэглэгч" }}
              </p>
              <p class="text-[11.5px] text-muted-foreground">
                {{ formatDate(app.created_at) }}
              </p>
            </div>

            <!-- Score -->
            <div
              class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full border-[2.5px] text-[12px] font-bold"
              :class="scoreBorderBg(app.overall_score)"
            >
              <span :class="scoreColor(app.overall_score)">
                {{ app.overall_score }}
              </span>
            </div>
          </div>

          <!-- Status badge -->
          <div class="mt-2 flex">
            <span
              class="inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-semibold"
              :class="statusConfig[derivedStatus(app)].cls"
            >
              {{ statusConfig[derivedStatus(app)].label }}
            </span>
          </div>
        </div>

        <!-- Pagination -->
        <div
          v-if="applicantTotalPages > 1"
          class="flex items-center justify-center gap-1 pt-1 pb-2 shrink-0"
        >
          <button
            class="flex h-7 w-7 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
            :disabled="applicantPage === 1"
            @click="applicantPage--"
          >
            <ChevronLeft class="h-3.5 w-3.5" />
          </button>
          <button
            v-for="p in applicantTotalPages"
            :key="p"
            :class="[
              'h-7 w-7 rounded-full text-[12px] font-medium transition',
              p === applicantPage
                ? 'bg-primary text-background'
                : 'text-muted-foreground hover:bg-muted',
            ]"
            @click="applicantPage = p"
          >
            {{ p }}
          </button>
          <button
            class="flex h-7 w-7 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
            :disabled="applicantPage === applicantTotalPages"
            @click="applicantPage++"
          >
            <ChevronRight class="h-3.5 w-3.5" />
          </button>
        </div>
      </div>
    </div>

    <!-- ── Right panel: detail ───────────────────────────────────────────── -->
    <div
      :class="[
        'flex-1 min-w-0 min-h-0',
        selected
          ? 'flex flex-col'
          : 'hidden lg:flex items-center justify-center',
      ]"
    >
      <!-- Mobile back button -->
      <button
        class="mb-3 flex items-center gap-1.5 text-[13px] text-muted-foreground transition hover:text-foreground lg:hidden shrink-0"
        @click="selected = null"
      >
        <ChevronLeft class="h-4 w-4" />
        Буцах
      </button>

      <!-- Empty state (desktop) -->
      <div
        v-if="!selected && !loading"
        class="flex flex-col items-center text-center"
      >
        <div
          class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-muted"
        >
          <Users class="h-7 w-7 text-muted-foreground" />
        </div>
        <p class="text-[15px] font-semibold">Горилогч сонгоно уу</p>
        <p class="mt-1.5 text-[13px] text-muted-foreground">
          Зүүн талын жагсаалтаас горилогч сонгоно уу.
        </p>
      </div>

      <!-- Detail panel -->
      <div
        v-else-if="selected"
        class="rounded-2xl border border-border bg-card overflow-y-auto min-h-0"
      >
        <!-- ── Header ── -->
        <div class="flex items-start gap-4 p-5 border-b border-border shrink-0">
          <!-- Avatar large -->
          <div
            class="flex h-14 w-14 shrink-0 items-center justify-center rounded-full bg-gradient-to-br text-white text-[18px] font-bold"
            :class="avatarGradient(selected.applicant_id)"
          >
            {{ initials(selected.applicant_name) }}
          </div>

          <!-- Name + meta -->
          <div class="flex-1 min-w-0">
            <h2
              class="text-[18px] font-semibold tracking-[-0.4px] leading-tight"
            >
              {{ selected.applicant_name || "Горилогч" }}
            </h2>
            <p class="mt-0.5 text-[12.5px] text-muted-foreground">
              {{ selected.applicant_email }}
            </p>
            <div class="mt-2 flex flex-wrap gap-1.5">
              <span
                class="inline-flex items-center rounded-full px-2.5 py-0.5 text-[11.5px] font-semibold"
                :class="statusConfig[derivedStatus(selected)].cls"
              >
                {{ statusConfig[derivedStatus(selected)].label }}
              </span>
              <span
                class="inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-[11.5px] text-muted-foreground"
              >
                Илгээсэн {{ formatDate(selected.created_at) }}
              </span>
              <span
                v-if="selected.interview_at"
                class="inline-flex items-center gap-1 rounded-full bg-primary/10 px-2.5 py-0.5 text-[11.5px] font-medium text-primary"
              >
                <CalendarDays class="h-3 w-3" />
                {{ formatDateTime(selected.interview_at) }}
              </span>
              <span
                v-if="selected.confidence_score"
                class="inline-flex items-center gap-1 rounded-full bg-primary/10 px-2.5 py-0.5 text-[11.5px] font-medium text-primary"
                title="Найдвартай байдал — CV-ийн мэдээллийн бүрэн гүйцэд байдал дээр тулгуурласан"
              >
                Confidence : {{ selected.confidence_score }}
              </span>
              <!-- <span
                v-if="selected.confidence_score"
                class="text-center"
                title="Найдвартай байдал — CV-ийн мэдээллийн бүрэн гүйцэд байдал дээр тулгуурласан"
              >
                <div
                  class="flex h-12 w-12 items-center justify-center rounded-full border-[2px] border-border bg-muted"
                >
                  <p class="text-[14px] font-semibold text-muted-foreground">
                    {{ selected.confidence_score }}
                  </p>
                </div>
                <p class="mt-1 text-[10px] text-muted-foreground">Confidence</p>
              </span> -->
            </div>
          </div>

          <!-- Score gauges -->
          <div class="shrink-0 flex flex-col items-center gap-2">
            <div class="text-center">
              <div
                class="flex h-16 w-16 items-center justify-center rounded-full border-[3.5px]"
                :class="scoreBorderBg(selected.overall_score)"
              >
                <div>
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
              <p class="mt-1 text-[10.5px] text-muted-foreground">AI оноо</p>
            </div>
          </div>
        </div>

        <!-- ── Body ── -->
        <div class="flex flex-col gap-4 p-5">
          <!-- AI Summary: full width -->
          <div class="rounded-xl border p-4 bg-primary/5">
            <div class="mb-2 flex items-center gap-1.5">
              <Sparkles class="h-3.5 w-3.5 text-primary" />
              <p
                class="text-[12px] font-semibold uppercase tracking-[0.5px] text-primary"
              >
                AI дүгнэлт
              </p>
            </div>
            <p class="text-[13px] leading-[1.6] text-muted-foreground">
              {{ selected.summary || "Дүгнэлт байхгүй." }}
            </p>
          </div>

          <!-- Strengths + Gaps: 2 columns -->
          <div class="grid grid-cols-2 gap-4">
            <!-- Strengths -->
            <div class="space-y-1.5">
              <p
                class="text-[10.5px] font-semibold uppercase tracking-[0.6px] text-success-foreground"
              >
                Тохирсон ур чадварууд
              </p>
              <div
                v-if="selected.matched_skills?.length"
                class="flex flex-wrap gap-1.5"
              >
                <span
                  v-for="s in selected.matched_skills"
                  :key="s.skill"
                  class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-[11.5px] font-medium badge-success"
                >
                  <CheckCircle2 class="h-3 w-3" />
                  {{ s.skill }}
                </span>
              </div>
              <p v-else class="text-[12.5px] text-muted-foreground">Байхгүй</p>
            </div>

            <!-- Gaps -->
            <div class="space-y-1.5">
              <p
                class="text-[10.5px] font-semibold uppercase tracking-[0.6px] text-destructive"
              >
                Дутуу ур чадварууд
              </p>
              <div
                v-if="selected.missing_skills?.length"
                class="flex flex-wrap gap-1.5"
              >
                <span
                  v-for="s in selected.missing_skills"
                  :key="s.skill"
                  class="inline-flex items-center gap-1 rounded-full bg-destructive/10 px-2.5 py-0.5 text-[11.5px] font-medium text-destructive"
                >
                  <XCircle class="h-3 w-3" />
                  {{ s.skill }}
                </span>
              </div>
              <p v-else class="text-[12.5px] text-muted-foreground">Байхгүй</p>
            </div>
          </div>
        </div>

        <!-- ── Footer: actions ── -->
        <div
          class="flex flex-wrap items-center gap-2 border-t border-border px-5 py-4 shrink-0"
        >
          <!-- Shortlist -->
          <button
            v-if="selected.status !== 'shortlisted'"
            class="inline-flex items-center gap-2 rounded-full px-4 py-2 text-[13px] font-semibold text-white transition hover:opacity-90 disabled:opacity-50"
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

          <!-- Schedule interview -->
          <button
            class="inline-flex items-center gap-1.5 rounded-full border border-border px-4 py-2 text-[13px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground"
            @click="openInterviewDialog(selected)"
          >
            <CalendarDays class="h-4 w-4" />
            {{
              selected.interview_at
                ? "Ярилцлагын тов өөрчлөх"
                : "Ярилцлага товлох"
            }}
          </button>

          <!-- Send task -->
          <button
            class="inline-flex items-center gap-1.5 rounded-full border border-border px-4 py-2 text-[13px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground disabled:opacity-40"
            :disabled="!!taskForApp(selected.id)"
            :title="
              taskForApp(selected.id)
                ? 'Даалгавар илгээгдсэн'
                : 'Даалгавар илгээх'
            "
            @click="openTaskDialog(selected)"
          >
            <ClipboardList class="h-4 w-4" />
            Даалгавар илгээх
          </button>

          <div class="flex-1" />

          <!-- Reject -->
          <button
            v-if="
              selected.status !== 'rejected' &&
              selected.status !== 'shortlisted'
            "
            class="inline-flex items-center gap-1.5 rounded-full border border-destructive/30 px-4 py-2 text-[13px] font-medium text-destructive transition hover:bg-destructive/10 disabled:opacity-50"
            :disabled="updatingStatus"
            @click="updateStatus(selected.id, 'rejected')"
          >
            Тэнцээгүй
          </button>
        </div>
      </div>
    </div>
  </div>

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

        <!-- Google Meet -->
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
            <Label>Дуусах огноо</Label>
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
