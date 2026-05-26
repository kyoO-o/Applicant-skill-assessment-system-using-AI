<script setup lang="ts">
import type { Task } from "../composables/types";
import { toast } from "vue-sonner";
import {
  ClipboardList,
  Plus,
  Loader2,
  ChevronLeft,
  ChevronRight,
  Sparkles,
  Pencil,
  Trash2,
  FileText,
} from "lucide-vue-next";
import { toTypedSchema } from "@vee-validate/zod";
import { taskSchema } from "~/utils/schemas";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";
import { Textarea } from "~/components/ui/textarea";

definePageMeta({ middleware: "auth" });

const tasksAPI = useTasksAPI();
const jobsAPI = useJobsAPI();
const { user } = useAuth();

const tasks = ref<Task[]>([]);
const loading = ref(true);
const isRecruiter = computed(() => user.value?.role === "recruiter");
const createOpen = ref(false);
const submitOpen = ref(false);
const selectedTask = ref<Task | null>(null);
const isSubmitting = ref(false);
const isGenerating = ref(false);
const editOpen = ref(false);
const editingTask = ref<Task | null>(null);
const deleteConfirmOpen = ref(false);
const deletingTask = ref<Task | null>(null);
const isDeleting = ref(false);
const editForm = reactive({ title: "", description: "", duration_days: "" });
const editKey = ref(0);

const jobs = ref<any[]>([]);
const selectedJobId = ref<string>("all");
const selectedJobIdNum = computed(() =>
  selectedJobId.value === "all" ? null : Number(selectedJobId.value),
);

const form = reactive({
  title: "",
  description: "",
  duration_days: "",
});
const taskFormSchema = toTypedSchema(taskSchema);

const submitContent = ref("");
const submitFile = ref<File | null>(null);
const submitFileInput = ref<HTMLInputElement | null>(null);

async function load() {
  loading.value = true;
  try {
    if (isRecruiter) {
      tasks.value = await tasksAPI.list(selectedJobIdNum.value ?? undefined);
    } else {
      tasks.value = await tasksAPI.list(selectedJobIdNum.value ?? undefined);
    }
  } catch {
    toast.error("Мэдээлэл ачааллаж чадсангүй");
  } finally {
    loading.value = false;
  }
}

async function loadJobs() {
  if (!isRecruiter.value) return;
  jobs.value = await jobsAPI.list().catch(() => []);
}

await loadJobs();
await load();

const taskPage = ref(1);
const TASK_PAGE_SIZE = 8;
const taskTotalPages = computed(() =>
  Math.max(1, Math.ceil(tasks.value.length / TASK_PAGE_SIZE)),
);
const paginatedTasks = computed(() => {
  const start = (taskPage.value - 1) * TASK_PAGE_SIZE;
  return tasks.value.slice(start, start + TASK_PAGE_SIZE);
});

watch(selectedJobId, () => {
  load();
  taskPage.value = 1;
});

async function generateTask() {
  if (!selectedJobIdNum.value) {
    toast.warning("Эхлээд ажлын байр сонгоно уу");
    return;
  }
  isGenerating.value = true;
  try {
    const result = await tasksAPI.generate(selectedJobIdNum.value);
    form.title = result.title;
    form.description = result.description;
    createOpen.value = true;
  } catch {
    toast.error("AI даалгавар үүсгэхэд алдаа гарлаа");
  } finally {
    isGenerating.value = false;
  }
}

async function createTask(values: Record<string, unknown>) {
  isSubmitting.value = true;
  try {
    await tasksAPI.create({
      title: values.title as string,
      description: (values.description as string) ?? "",
      duration_days: values.duration_days
        ? Number(values.duration_days)
        : undefined,
      job_posting_id: selectedJobIdNum.value ?? undefined,
    });
    createOpen.value = false;
    form.title = "";
    form.description = "";
    form.duration_days = "";
    toast.success("Даалгавар үүсгэгдлээ");
    await load();
  } catch {
    toast.error("Алдаа гарлаа");
  } finally {
    isSubmitting.value = false;
  }
}

async function submitTask() {
  if (!submitContent.value.trim() && !submitFile.value) {
    toast.warning("Агуулга эсвэл PDF файл шаардлагатай");
    return;
  }
  isSubmitting.value = true;
  try {
    await tasksAPI.submit(
      selectedTask.value!.id,
      submitContent.value,
      submitFile.value ?? undefined,
    );
    submitOpen.value = false;
    submitContent.value = "";
    submitFile.value = null;
    toast.success("Даалгавар амжилттай илгээгдлээ");
    await load();
  } catch {
    toast.error("Алдаа гарлаа");
  } finally {
    isSubmitting.value = false;
  }
}

function onFileChange(e: Event) {
  const target = e.target as HTMLInputElement;
  const file = target.files?.[0];
  if (file && file.type !== "application/pdf") {
    toast.warning("Зөвхөн PDF файл зөвшөөрнө");
    target.value = "";
    return;
  }
  submitFile.value = file ?? null;
}

function openDelete(task: Task) {
  deletingTask.value = task;
  deleteConfirmOpen.value = true;
}

async function confirmDelete() {
  if (!deletingTask.value) return;
  isDeleting.value = true;
  try {
    await tasksAPI.delete(deletingTask.value.id);
    deleteConfirmOpen.value = false;
    toast.success("Даалгавар устгагдлаа");
    await load();
  } catch {
    toast.error("Устгахад алдаа гарлаа");
  } finally {
    isDeleting.value = false;
  }
}

function openEdit(task: Task) {
  editingTask.value = task;
  editForm.title = task.title;
  editForm.description = task.description;
  editForm.duration_days = task.duration_days ? String(task.duration_days) : "";
  editKey.value++;
  editOpen.value = true;
}

async function updateTask(values: Record<string, unknown>) {
  if (!editingTask.value) return;
  isSubmitting.value = true;
  try {
    await tasksAPI.update(editingTask.value.id, {
      title: values.title as string,
      description: (values.description as string) ?? "",
      duration_days: values.duration_days
        ? Number(values.duration_days)
        : undefined,
    });
    editOpen.value = false;
    toast.success("Даалгавар шинэчлэгдлээ");
    await load();
  } catch {
    toast.error("Алдаа гарлаа");
  } finally {
    isSubmitting.value = false;
  }
}

function formatDate(d: string | null | undefined) {
  if (!d) return "";
  return new Intl.DateTimeFormat("mn-MN", { dateStyle: "medium" }).format(
    new Date(d),
  );
}
</script>

<template>
  <div class="space-y-5 min-w-0 overflow-x-hidden">
    <!-- Page header -->
    <div class="flex items-start justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-semibold tracking-[-0.6px]">
          {{ isRecruiter ? "Даалгаврын сан" : "Миний даалгаврууд" }}
        </h1>
        <p class="mt-1 text-[13.5px] text-muted-foreground">
          {{
            isRecruiter
              ? "Горилогчдод илгээх даалгаврын загваруудаа энд үүсгэнэ үү."
              : "Та бүртгүүлсэн ажлын байрнаасаа хүлээн авсан даалгаврууд."
          }}
        </p>
      </div>
      <!-- Create buttons — visible only when a specific job is selected -->
      <div
        v-if="isRecruiter && selectedJobIdNum"
        class="flex shrink-0 items-center gap-2"
      >
        <button
          class="inline-flex items-center gap-2 rounded-full border border-border px-4 py-2 text-[13.5px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground disabled:opacity-50"
          :disabled="isGenerating"
          @click="generateTask"
        >
          <Loader2 v-if="isGenerating" class="h-4 w-4 animate-spin" />
          <Sparkles v-else class="h-4 w-4 text-primary" />
          {{ isGenerating ? "Үүсгэж байна..." : "AI-аар үүсгэх" }}
        </button>
        <button
          class="inline-flex items-center gap-2 rounded-full px-4 py-2 text-[13.5px] font-semibold text-white transition hover:opacity-90"
          style="
            background: linear-gradient(
              135deg,
              var(--primary),
              oklch(0.348 0.106 295)
            );
          "
          @click="createOpen = true"
        >
          <Plus class="h-4 w-4" /> Гараар үүсгэх
        </button>
      </div>
    </div>

    <!-- Job selector for recruiters -->
    <div v-if="isRecruiter && jobs.length">
      <Select v-model="selectedJobId">
        <SelectTrigger class="w-64">
          <SelectValue placeholder="Ажлын байр сонгох" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">Бүх ажлын байр</SelectItem>
          <SelectItem v-for="job in jobs" :key="job.id" :value="String(job.id)">
            {{ job.title }}
          </SelectItem>
        </SelectContent>
      </Select>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="h-7 w-7 animate-spin text-primary" />
    </div>

    <!-- Empty -->
    <div
      v-else-if="!tasks.length"
      class="flex flex-col items-center py-20 text-center"
    >
      <div
        class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-muted"
      >
        <ClipboardList class="h-7 w-7 text-muted-foreground" />
      </div>
      <p class="text-[16px] font-semibold">Даалгавар байхгүй</p>
      <p v-if="isRecruiter" class="mt-2 text-[13.5px] text-muted-foreground">
        Дээрх товчлуураар даалгавар үүсгэж санд хадгалаарай.
      </p>
    </div>

    <!-- Task cards -->
    <div v-else class="w-full space-y-3.5">
      <div
        v-for="task in paginatedTasks"
        :key="task.id"
        class="rounded-2xl border border-border bg-card p-5 transition hover:border-primary/15"
      >
        <div class="flex items-start justify-between gap-4">
          <div class="flex-1 min-w-0 space-y-1">
            <!-- Job label (applicant only) -->
            <p
              v-if="!isRecruiter && task.job_title"
              class="text-[11.5px] font-medium text-primary/80 uppercase tracking-[0.5px]"
            >
              {{ task.job_title }}
            </p>
            <div class="flex items-center gap-2 flex-wrap">
              <h3 class="text-[15px] font-semibold">{{ task.title }}</h3>
              <span
                v-if="task.created_by_ai"
                class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[11px] font-semibold bg-primary/10 text-primary"
              >
                <Sparkles class="h-2.5 w-2.5" /> AI
              </span>
              <!-- Status badge (applicant only) -->
              <span
                v-if="!isRecruiter"
                :class="[
                  'inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-semibold',
                  task.status === 'sent'
                    ? 'bg-amber-500/10 text-amber-600 dark:text-amber-400'
                    : task.status === 'completed'
                      ? 'bg-blue-500/10 text-blue-600 dark:text-blue-400'
                      : 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400',
                ]"
              >
                {{
                  task.status === "sent"
                    ? "Хүлээгдэж байна"
                    : task.status === "completed"
                      ? "Илгээсэн"
                      : "Үнэлэгдсэн"
                }}
              </span>
            </div>
            <p class="text-[13.5px] text-muted-foreground line-clamp-2 w-200">
              {{ task.description }}
            </p>
            <div class="flex flex-wrap items-center gap-3">
              <p
                v-if="task.duration_days"
                class="text-[12px] text-muted-foreground"
              >
                Хугацаа: {{ task.duration_days }} өдөр
              </p>
              <p v-if="task.due_date" class="text-[12px] text-muted-foreground">
                Дуусах огноо: {{ formatDate(task.due_date) }}
              </p>
            </div>
          </div>
          <div class="flex shrink-0 items-center gap-2">
            <!-- Recruiter: edit + delete -->
            <template v-if="isRecruiter">
              <button
                class="inline-flex items-center gap-1.5 rounded-full border border-border px-3 py-1.5 text-[12.5px] font-medium text-muted-foreground transition hover:bg-muted hover:text-foreground"
                @click="openEdit(task)"
              >
                <Pencil class="h-3.5 w-3.5" /> Засах
              </button>
              <button
                v-if="task.status === 'draft'"
                class="inline-flex items-center gap-1.5 rounded-full border border-destructive/30 px-3 py-1.5 text-[12.5px] font-medium text-destructive transition hover:bg-destructive/10"
                @click="openDelete(task)"
              >
                <Trash2 class="h-3.5 w-3.5" /> Устгах
              </button>
            </template>
            <!-- Applicant: submit -->
            <button
              v-if="!isRecruiter && task.status === 'sent'"
              class="inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-[12.5px] font-semibold text-white transition hover:opacity-90"
              style="
                background: linear-gradient(
                  135deg,
                  var(--primary),
                  oklch(0.348 0.106 295)
                );
              "
              @click="
                selectedTask = task;
                submitOpen = true;
              "
            >
              <ChevronRight class="h-3.5 w-3.5" /> Гүйцэтгэл илгээх
            </button>
          </div>
        </div>

        <!-- Applicant: own submission result -->
        <div
          v-if="!isRecruiter && task.submissions?.length"
          class="mt-4 border-t border-border pt-4 space-y-2"
        >
          <div
            v-for="sub in task.submissions"
            :key="sub.id"
            class="rounded-xl bg-muted/40 px-4 py-3 space-y-2"
          >
            <!-- <p class="text-[13px] text-foreground line-clamp-3">
              {{ sub.content }}
            </p> -->
            <div v-if="sub.grade !== null" class="space-y-1">
              <p
                class="flex items-center gap-1.5 text-[12.5px] font-semibold text-emerald-600 dark:text-emerald-400"
              >
                <span class="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                Оноо: {{ sub.grade }}/100
              </p>
              <p
                v-if="sub.feedback"
                class="text-[12.5px] text-muted-foreground"
              >
                {{ sub.feedback }}
              </p>
            </div>
            <p v-else class="text-[12px] text-muted-foreground italic">
              Үнэлгээ хүлээгдэж байна...
            </p>
          </div>
        </div>

        <!-- Recruiter: submissions -->
        <!-- <div
          v-if="isRecruiter && task.submissions?.length"
          class="mt-4 space-y-2 border-t border-border pt-4"
        >
          <p
            class="text-[12px] font-semibold uppercase tracking-[0.5px] text-muted-foreground"
          >
            Илгээлтүүд ({{ task.submissions.length }})
          </p>
          <div
            v-for="sub in task.submissions"
            :key="sub.id"
            class="rounded-xl bg-muted/40 px-4 py-3 text-[13.5px]"
          >
            <p class="line-clamp-2 text-foreground">{{ sub.content }}</p>
            <p
              v-if="sub.grade !== null"
              class="mt-1.5 flex items-center gap-1.5 text-[12px] font-semibold text-success-foreground"
            >
              <span class="h-1.5 w-1.5 rounded-full bg-success" />
              Оноо: {{ sub.grade }}/100
            </p>
          </div>
        </div> -->
      </div>

      <!-- Task pagination -->
      <div class="flex items-center justify-center gap-1 pt-2">
        <button
          class="flex h-8 w-8 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
          :disabled="taskPage === 1"
          @click="taskPage--"
        >
          <ChevronLeft class="h-4 w-4" />
        </button>
        <button
          v-for="p in taskTotalPages"
          :key="p"
          :class="[
            'h-8 w-8 rounded-full text-[13px] font-medium transition',
            p === taskPage
              ? 'bg-primary text-background'
              : 'text-muted-foreground hover:bg-muted',
          ]"
          @click="taskPage = p"
        >
          {{ p }}
        </button>
        <button
          class="flex h-8 w-8 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
          :disabled="taskPage === taskTotalPages"
          @click="taskPage++"
        >
          <ChevronRight class="h-4 w-4" />
        </button>
      </div>
    </div>
  </div>

  <!-- Create task dialog -->
  <Dialog v-model:open="createOpen">
    <DialogContent
      class="rounded-2xl sm:max-w-2xl max-h-[90vh] overflow-y-auto"
    >
      <DialogHeader>
        <DialogTitle>Даалгавар үүсгэх</DialogTitle>
        <DialogDescription>
          Санд хадгалагдах загвар даалгавар үүсгэнэ үү.
        </DialogDescription>
      </DialogHeader>
      <Form
        :validation-schema="taskFormSchema"
        :initial-values="form"
        class="space-y-4 py-2"
        @submit="createTask"
      >
        <FormField v-slot="{ componentField }" name="title">
          <FormItem>
            <FormLabel>Гарчиг</FormLabel>
            <FormControl>
              <Input v-bind="componentField" placeholder="Даалгаврын гарчиг" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="description">
          <FormItem>
            <FormLabel>Тайлбар</FormLabel>
            <FormControl>
              <textarea
                v-bind="componentField"
                rows="9"
                placeholder="Даалгаврын дэлгэрэнгүй тайлбар..."
                class="w-full resize-y rounded-xl border border-input bg-background px-3 py-2 text-[13.5px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="duration_days">
          <FormItem>
            <FormLabel> Хугацаа (өдөр) </FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="number"
                min="1"
                placeholder="Жишээ: 7"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <DialogFooter>
          <Button
            variant="outline"
            type="button"
            :disabled="isSubmitting"
            @click="createOpen = false"
            >Болих</Button
          >
          <Button type="submit" :disabled="isSubmitting">
            <Loader2 v-if="isSubmitting" class="mr-2 h-4 w-4 animate-spin" />
            Хадгалах
          </Button>
        </DialogFooter>
      </Form>
    </DialogContent>
  </Dialog>

  <!-- Edit task dialog -->
  <Dialog v-model:open="editOpen">
    <DialogContent
      class="rounded-2xl sm:max-w-2xl max-h-[90vh] overflow-y-auto"
    >
      <DialogHeader>
        <DialogTitle>Даалгавар засах</DialogTitle>
        <DialogDescription
          >Даалгаврын мэдээллийг шинэчилнэ үү.</DialogDescription
        >
      </DialogHeader>
      <Form
        :key="editKey"
        :validation-schema="taskFormSchema"
        :initial-values="editForm"
        class="space-y-4 py-2"
        @submit="updateTask"
      >
        <FormField v-slot="{ componentField }" name="title">
          <FormItem>
            <FormLabel>Гарчиг</FormLabel>
            <FormControl>
              <Input v-bind="componentField" placeholder="Даалгаврын гарчиг" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="description">
          <FormItem>
            <FormLabel>Тайлбар</FormLabel>
            <FormControl>
              <Textarea
                v-bind="componentField"
                placeholder="Даалгаврын дэлгэрэнгүй тайлбар..."
                class="min-h-[200px] max-h-[280px] text-[13.5px]"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="duration_days">
          <FormItem>
            <FormLabel>
              Хугацаа (өдөр)
              <span class="text-[11.5px] text-muted-foreground"
                >(заавал биш)</span
              >
            </FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="number"
                min="1"
                placeholder="Жишээ: 7"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <DialogFooter>
          <Button
            variant="outline"
            type="button"
            :disabled="isSubmitting"
            @click="editOpen = false"
            >Болих</Button
          >
          <Button type="submit" :disabled="isSubmitting">
            <Loader2 v-if="isSubmitting" class="mr-2 h-4 w-4 animate-spin" />
            Хадгалах
          </Button>
        </DialogFooter>
      </Form>
    </DialogContent>
  </Dialog>

  <!-- Delete task confirmation -->
  <AlertDialog v-model:open="deleteConfirmOpen">
    <AlertDialogContent class="rounded-2xl">
      <AlertDialogHeader>
        <AlertDialogTitle>Даалгавар устгах уу?</AlertDialogTitle>
        <AlertDialogDescription>
          <span class="font-medium text-foreground">{{
            deletingTask?.title
          }}</span>
          даалгаврыг устгах гэж байна. Энэ үйлдлийг буцаах боломжгүй.
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel :disabled="isDeleting">Болих</AlertDialogCancel>
        <AlertDialogAction
          class="bg-destructive text-destructive-foreground hover:bg-destructive/90"
          :disabled="isDeleting"
          @click.prevent="confirmDelete"
        >
          <Loader2 v-if="isDeleting" class="mr-2 h-4 w-4 animate-spin" />
          Устгах
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>

  <!-- Submit task dialog -->
  <Dialog v-model:open="submitOpen">
    <DialogContent class="rounded-2xl sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>{{ selectedTask?.title }}</DialogTitle>
        <DialogDescription>{{ selectedTask?.description }}</DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-2">
        <div class="space-y-2">
          <Label
            >Гүйцэтгэлийн агуулга
            <span class="text-[11.5px] text-muted-foreground"
              >(заавал биш)</span
            ></Label
          >
          <textarea
            v-model="submitContent"
            rows="5"
            placeholder="Гүйцэтгэсэн ажлынхаа тайлбар, хариулт эсвэл холбоосыг энд бичнэ үү..."
            class="w-full resize-none rounded-xl border border-input bg-background px-3 py-2 text-[13.5px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
        </div>
        <div class="space-y-2">
          <Label
            >PDF файл
            <span class="text-[11.5px] text-muted-foreground"
              >(заавал биш)</span
            ></Label
          >
          <div
            class="flex items-center gap-3 rounded-xl border border-dashed border-input bg-muted/30 px-4 py-3 cursor-pointer hover:bg-muted/50 transition"
            @click="submitFileInput?.click()"
          >
            <input
              ref="submitFileInput"
              type="file"
              accept=".pdf"
              class="hidden"
              @change="onFileChange"
            />
            <FileText class="h-4 w-4 text-muted-foreground flex-shrink-0" />
            <span class="text-[13px] text-muted-foreground flex-1 truncate">
              {{ submitFile ? submitFile.name : "PDF файл сонгох..." }}
            </span>
            <button
              v-if="submitFile"
              class="text-[11px] text-destructive hover:underline"
              @click.stop="submitFile = null"
            >
              Хасах
            </button>
          </div>
        </div>
      </div>
      <DialogFooter>
        <Button
          variant="outline"
          @click="submitOpen = false"
          :disabled="isSubmitting"
          >Болих</Button
        >
        <Button @click="submitTask" :disabled="isSubmitting">
          <Loader2 v-if="isSubmitting" class="mr-2 h-4 w-4 animate-spin" />
          Илгээх
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
