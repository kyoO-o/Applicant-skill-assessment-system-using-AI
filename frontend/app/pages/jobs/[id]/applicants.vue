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
const interviewNote = ref("");
const isScheduling = ref(false);

const companyAPI = useCompanyAPI();
const companyLocation = ref("");

onMounted(async () => {
  try {
    const company = await companyAPI.get();
    const parts = [company.city, company.district].filter(Boolean);
    companyLocation.value = parts.join(", ");
  } catch {
    // no company yet — that's fine
  }
});

function openInterviewDialog(app: Application) {
  interviewTarget.value = app;
  interviewDate.value = "";
  interviewTime.value = "";
  interviewLocation.value = app.interview_location || companyLocation.value;
  interviewNote.value = app.interview_note || "";
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
    const updated = await applicationsAPI.scheduleInterview(
      interviewTarget.value!.id,
      interviewAt,
      interviewLocation.value,
      interviewNote.value,
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

function openTaskDialog(app: Application) {
  taskTarget.value = app;
  taskTitle.value = "";
  taskDescription.value = "";
  taskDueDate.value = "";
  taskOpen.value = true;
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
  if (!taskTitle.value.trim() || !taskDescription.value.trim()) {
    toast.error("Гарчиг болон тайлбар шаардлагатай");
    return;
  }
  isSendingTask.value = true;
  try {
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
    if (task.status === "completed" || task.status === "graded") return "task_completed";
    if (task.status === "sent") return "task_sent";
  }
  if (app.interview_at) return "interview_scheduled";
  return "applied";
}

const statusConfig: Record<DerivedStatus, { label: string; class: string }> = {
  applied: {
    label: "Анкет илгээсэн",
    class: "bg-blue-50 text-blue-700 border-blue-200",
  },
  interview_scheduled: {
    label: "Ярилцлага товлосон",
    class: "bg-purple-50 text-purple-700 border-purple-200",
  },
  task_sent: {
    label: "Даалгавар илгээсэн",
    class: "bg-amber-50 text-amber-700 border-amber-200",
  },
  task_completed: {
    label: "Даалгавар гүйцэтгэсэн",
    class: "bg-teal-50 text-teal-700 border-teal-200",
  },
  hired: {
    label: "Ажилд авсан",
    class: "bg-green-50 text-green-700 border-green-200",
  },
  failed: {
    label: "Тэнцээгүй",
    class: "bg-red-50 text-red-700 border-red-200",
  },
};

const statusOrder: Record<DerivedStatus, number> = {
  hired: 0,
  task_completed: 1,
  task_sent: 2,
  interview_scheduled: 3,
  applied: 4,
  failed: 5,
};

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
        (a, b) =>
          statusOrder[derivedStatus(a)] - statusOrder[derivedStatus(b)],
      );
    case "date":
      return arr.sort(
        (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
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
  if (score >= 75) return "text-green-600";
  if (score >= 50) return "text-yellow-600";
  return "text-red-600";
}

function scoreBadge(score: number) {
  if (score >= 75) return "bg-green-50 text-green-700";
  if (score >= 50) return "bg-yellow-50 text-yellow-700";
  return "bg-red-50 text-red-700";
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

function assessIcon(s: string) {
  if (s === "met") return CheckCircle2;
  if (s === "partial") return MinusCircle;
  return XCircle;
}

function assessColor(s: string) {
  if (s === "met") return "text-green-600";
  if (s === "partial") return "text-yellow-600";
  return "text-red-500";
}

function assessBg(s: string) {
  if (s === "met") return "bg-green-50 border-green-100";
  if (s === "partial") return "bg-yellow-50 border-yellow-100";
  return "bg-red-50 border-red-100";
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <Button variant="ghost" size="icon" class="rounded-xl" @click="router.back()">
            <ChevronLeft class="h-4 w-4" />
          </Button>
          <div>
            <p class="text-sm font-medium text-muted-foreground">Ажлын горилогчид</p>
            <h2 class="text-2xl font-semibold tracking-tight">
              {{ applications.length }} горилогч
            </h2>
          </div>
        </div>

        <!-- Sort -->
        <div class="flex items-center gap-2">
          <ArrowUpDown class="h-4 w-4 shrink-0 text-muted-foreground" />
          <Select v-model="sortBy">
            <SelectTrigger class="w-52 rounded-2xl">
              <SelectValue placeholder="Эрэмбэлэх" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="opt in sortOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>
    </section>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <!-- Empty -->
    <div
      v-else-if="!applications.length"
      class="rounded-3xl border border-dashed border-border px-6 py-20 text-center"
    >
      <Users class="mx-auto h-10 w-10 text-muted-foreground" />
      <p class="mt-4 text-lg font-medium">Горилогч байхгүй байна</p>
      <p class="mt-2 text-sm text-muted-foreground">
        Горилогчид анкет илгээх үед энд харагдана.
      </p>
    </div>

    <!-- Table -->
    <Card v-else class="rounded-3xl border-border shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead class="min-w-[180px]">Горилогч</TableHead>
              <TableHead class="min-w-[90px] text-center">Үнэлгээ</TableHead>
              <TableHead class="min-w-[160px]">Төлөв</TableHead>
              <TableHead class="min-w-[120px]">Илгээсэн огноо</TableHead>
              <TableHead class="min-w-[140px]">Ярилцлагийн тов</TableHead>
              <TableHead class="min-w-[140px]">Даалгавар илгээсэн</TableHead>
              <TableHead class="min-w-[150px]">Даалгавар гүйцэтгэсэн</TableHead>
              <TableHead class="min-w-[180px] text-right">Үйлдэл</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="app in sortedApplications"
              :key="app.id"
              class="cursor-pointer hover:bg-muted/30 transition-colors"
              @click="openDetail(app)"
            >
              <!-- Name -->
              <TableCell>
                <div class="min-w-0">
                  <p class="truncate font-medium text-sm">{{ app.applicant_name || "Хэрэглэгч" }}</p>
                  <p class="truncate text-xs text-muted-foreground">{{ app.applicant_email }}</p>
                </div>
              </TableCell>

              <!-- Score -->
              <TableCell class="text-center">
                <span
                  :class="[
                    'inline-block rounded-xl px-2.5 py-1 text-sm font-bold',
                    scoreBadge(app.overall_score),
                  ]"
                >
                  {{ app.overall_score }}%
                </span>
              </TableCell>

              <!-- Status -->
              <TableCell>
                <span
                  :class="[
                    'inline-flex items-center rounded-full border px-3 py-1 text-xs font-medium',
                    statusConfig[derivedStatus(app)].class,
                  ]"
                >
                  {{ statusConfig[derivedStatus(app)].label }}
                </span>
              </TableCell>

              <!-- Applied date -->
              <TableCell class="text-sm text-muted-foreground">
                {{ formatDate(app.created_at) }}
              </TableCell>

              <!-- Interview date -->
              <TableCell class="text-sm text-muted-foreground">
                {{ app.interview_at ? formatDateTime(app.interview_at) : "—" }}
              </TableCell>

              <!-- Task sent date -->
              <TableCell class="text-sm text-muted-foreground">
                <template v-if="taskForApp(app.id)">
                  {{ formatDate(taskForApp(app.id)!.created_at) }}
                </template>
                <template v-else>—</template>
              </TableCell>

              <!-- Task completed date -->
              <TableCell class="text-sm text-muted-foreground">
                <template
                  v-if="taskForApp(app.id) && (taskForApp(app.id)!.status === 'completed' || taskForApp(app.id)!.status === 'graded')"
                >
                  {{ formatDate(taskForApp(app.id)!.updated_at) }}
                </template>
                <template v-else>—</template>
              </TableCell>

              <!-- Actions -->
              <TableCell class="text-right" @click.stop>
                <div class="flex justify-end gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    class="rounded-xl gap-1.5"
                    :title="app.interview_at ? 'Ярилцлага товлогдсон' : 'Ярилцлага товлох'"
                    @click="openInterviewDialog(app)"
                  >
                    <CalendarDays class="h-3.5 w-3.5" />
                    <span class="hidden sm:inline">Ярилцлага</span>
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    class="rounded-xl gap-1.5"
                    :disabled="!!taskForApp(app.id)"
                    :title="taskForApp(app.id) ? 'Даалгавар илгээгдсэн' : 'Даалгавар илгээх'"
                    @click="openTaskDialog(app)"
                  >
                    <ClipboardList class="h-3.5 w-3.5" />
                    <span class="hidden sm:inline">Даалгавар</span>
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>
    </Card>
  </div>

  <!-- ── Detail dialog ──────────────────────────────────────────────────── -->
  <Dialog v-model:open="detailOpen">
    <DialogContent class="rounded-3xl sm:max-w-2xl max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>{{ selected?.applicant_name || "Горилогч" }}</DialogTitle>
        <DialogDescription>{{ selected?.applicant_email }}</DialogDescription>
      </DialogHeader>

      <div v-if="selected" class="space-y-5 py-2">
        <!-- Score + summary -->
        <div class="flex items-start gap-4 rounded-2xl border border-border bg-muted/20 p-4">
          <div class="shrink-0 text-center">
            <p :class="['text-4xl font-bold', scoreColor(selected.overall_score)]">
              {{ selected.overall_score }}%
            </p>
            <p class="text-xs text-muted-foreground mt-1">үнэлгээ</p>
          </div>
          <div class="flex-1">
            <p class="text-sm leading-6 text-muted-foreground">{{ selected.summary }}</p>
          </div>
        </div>

        <!-- Status timeline -->
        <div class="grid gap-2 sm:grid-cols-2">
          <div class="rounded-2xl border border-border bg-muted/20 px-4 py-3">
            <p class="text-xs text-muted-foreground">Одоогийн төлөв</p>
            <span
              :class="[
                'mt-1 inline-flex items-center rounded-full border px-3 py-1 text-xs font-medium',
                statusConfig[derivedStatus(selected)].class,
              ]"
            >
              {{ statusConfig[derivedStatus(selected)].label }}
            </span>
          </div>
          <div class="rounded-2xl border border-border bg-muted/20 px-4 py-3">
            <p class="text-xs text-muted-foreground">Ярилцлагийн тов</p>
            <p class="mt-1 text-sm font-medium">
              {{ selected.interview_at ? formatDateTime(selected.interview_at) : "Товлоогүй" }}
            </p>
          </div>
        </div>

        <!-- Hire / Reject + Interview / Task -->
        <div class="flex flex-wrap gap-2">
          <Button
            v-if="selected.status !== 'shortlisted'"
            class="rounded-full"
            :disabled="updatingStatus"
            @click="updateStatus(selected.id, 'shortlisted')"
          >
            Ажилд авах
          </Button>
          <Button
            v-if="selected.status !== 'rejected'"
            variant="outline"
            class="rounded-full"
            :disabled="updatingStatus"
            @click="updateStatus(selected.id, 'rejected')"
          >
            Тэнцээгүй
          </Button>
          <Button
            variant="outline"
            class="rounded-full"
            @click="() => { detailOpen = false; openInterviewDialog(selected!); }"
          >
            <CalendarDays class="mr-2 h-4 w-4" />
            Ярилцлага товлох
          </Button>
          <Button
            variant="outline"
            class="rounded-full"
            :disabled="!!taskForApp(selected.id)"
            @click="() => { detailOpen = false; openTaskDialog(selected!); }"
          >
            <ClipboardList class="mr-2 h-4 w-4" />
            Даалгавар илгээх
          </Button>
        </div>

        <!-- Skills -->
        <div class="grid gap-4 sm:grid-cols-2">
          <div class="space-y-2">
            <p class="text-xs font-medium uppercase tracking-wide text-green-700">
              Тохирсон ур чадварууд
            </p>
            <div
              v-for="s in selected.matched_skills"
              :key="s.skill"
              class="rounded-xl bg-green-50 px-3 py-2"
            >
              <p class="text-xs font-medium text-green-800">✓ {{ s.skill }}</p>
              <p class="mt-0.5 text-xs text-green-600">{{ s.explanation }}</p>
            </div>
            <p v-if="!selected.matched_skills?.length" class="text-xs text-muted-foreground">
              Байхгүй
            </p>
          </div>
          <div class="space-y-2">
            <p class="text-xs font-medium uppercase tracking-wide text-red-700">
              Дутуу ур чадварууд
            </p>
            <div
              v-for="s in selected.missing_skills"
              :key="s.skill"
              class="rounded-xl bg-red-50 px-3 py-2"
            >
              <p class="text-xs font-medium text-red-800">✗ {{ s.skill }}</p>
              <p class="mt-0.5 text-xs text-red-600">{{ s.explanation }}</p>
            </div>
            <p v-if="!selected.missing_skills?.length" class="text-xs text-muted-foreground">
              Байхгүй
            </p>
          </div>
        </div>

        <!-- Duty assessments -->
        <div v-if="selected.duty_assessments?.length" class="space-y-2">
          <p class="text-sm font-medium">Үүрэг хариуцлагын үнэлгээ</p>
          <div
            v-for="d in selected.duty_assessments"
            :key="d.duty"
            :class="['rounded-2xl border p-3', assessBg(d.status)]"
          >
            <div class="flex items-start gap-2">
              <component
                :is="assessIcon(d.status)"
                class="mt-0.5 h-4 w-4 shrink-0"
                :class="assessColor(d.status)"
              />
              <div>
                <p class="text-sm font-medium">{{ d.duty }}</p>
                <p class="mt-0.5 text-xs" :class="assessColor(d.status)">
                  {{ d.explanation }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Requirement assessments -->
        <div v-if="selected.requirement_assessments?.length" class="space-y-2">
          <p class="text-sm font-medium">Шаардлагын үнэлгээ</p>
          <div
            v-for="r in selected.requirement_assessments"
            :key="r.requirement"
            :class="['rounded-2xl border p-3', assessBg(r.status)]"
          >
            <div class="flex items-start gap-2">
              <component
                :is="assessIcon(r.status)"
                class="mt-0.5 h-4 w-4 shrink-0"
                :class="assessColor(r.status)"
              />
              <div>
                <p class="text-sm font-medium">{{ r.requirement }}</p>
                <p class="mt-0.5 text-xs" :class="assessColor(r.status)">
                  {{ r.explanation }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Recommendations -->
        <div v-if="selected.recommendations?.length" class="space-y-2">
          <p class="text-sm font-medium">Зөвлөмж</p>
          <ul class="space-y-1">
            <li
              v-for="r in selected.recommendations"
              :key="r"
              class="flex items-start gap-2 text-sm text-muted-foreground"
            >
              <ChevronRight class="mt-0.5 h-4 w-4 shrink-0 text-primary" /> {{ r }}
            </li>
          </ul>
        </div>
      </div>
    </DialogContent>
  </Dialog>

  <!-- ── Interview dialog ───────────────────────────────────────────────── -->
  <Dialog v-model:open="interviewOpen">
    <DialogContent class="rounded-3xl sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Ярилцлага товлох</DialogTitle>
        <DialogDescription>
          {{ interviewTarget?.applicant_name }}-д ярилцлагын дэлгэрэнгүй мэдээллийг оруулна уу.
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-2">
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
        <div class="space-y-2">
          <Label>Байршил</Label>
          <Input v-model="interviewLocation" placeholder="Уулзах газар (жишээ: Конкорд Тауэр, 5 давхар)" />
        </div>
        <div class="space-y-2">
          <Label>Нэмэлт тэмдэглэл <span class="text-muted-foreground text-xs">(заавал биш)</span></Label>
          <Textarea v-model="interviewNote" rows="3" placeholder="Горилогчид дамжуулах нэмэлт мэдээлэл..." />
        </div>
      </div>
      <DialogFooter>
        <Button variant="outline" :disabled="isScheduling" @click="interviewOpen = false">
          Болих
        </Button>
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
    <DialogContent class="rounded-3xl sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>Даалгавар илгээх</DialogTitle>
        <DialogDescription>
          {{ taskTarget?.applicant_name }}-д практик даалгавар илгээнэ үү.
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-2">
        <Button
          type="button"
          variant="outline"
          class="w-full rounded-2xl"
          :disabled="isGeneratingTask"
          @click="generateTask"
        >
          <Loader2 v-if="isGeneratingTask" class="mr-2 h-4 w-4 animate-spin" />
          <Sparkles v-else class="mr-2 h-4 w-4" />
          {{ isGeneratingTask ? "AI үүсгэж байна..." : "AI-аар үүсгэх" }}
        </Button>
        <div class="space-y-2">
          <Label>Гарчиг</Label>
          <Input v-model="taskTitle" placeholder="Даалгаврын гарчиг" />
        </div>
        <div class="space-y-2">
          <Label>Тайлбар</Label>
          <Textarea v-model="taskDescription" rows="5" placeholder="Даалгаврын дэлгэрэнгүй тайлбар..." />
        </div>
        <div class="space-y-2">
          <Label>Дуусах огноо (заавал биш)</Label>
          <Input v-model="taskDueDate" type="date" />
        </div>
      </div>
      <DialogFooter>
        <Button variant="outline" :disabled="isSendingTask" @click="taskOpen = false">
          Болих
        </Button>
        <Button :disabled="isSendingTask || !taskTitle || !taskDescription" @click="sendTask">
          <Loader2 v-if="isSendingTask" class="mr-2 h-4 w-4 animate-spin" />
          <ClipboardList v-else class="mr-2 h-4 w-4" />
          {{ isSendingTask ? "Илгээж байна..." : "Илгээх" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
