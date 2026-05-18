<script setup lang="ts">
import type { Task } from "../composables/types";
import { toast } from "vue-sonner";
import { ClipboardList, Plus, Wand2, Send, Loader2, ChevronRight, Sparkles } from "lucide-vue-next";
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

const jobs = ref<any[]>([]);
const selectedJobID = ref<number | null>(null);

const form = reactive({
  title: "",
  description: "",
  due_date: "",
});
const taskFormSchema = toTypedSchema(taskSchema);

const submitContent = ref("");

async function load() {
  loading.value = true;
  try {
    if (isRecruiter.value) {
      if (selectedJobID.value) {
        tasks.value = await tasksAPI.list(selectedJobID.value);
      } else {
        tasks.value = [];
      }
    } else {
      tasks.value = await tasksAPI.list();
    }
  } catch {
    toast.error("Мэдээлэл ачааллаж чадсангүй");
  } finally {
    loading.value = false;
  }
}

async function loadJobs() {
  if (!isRecruiter.value) return;
  const all = await jobsAPI.list().catch(() => []);
  jobs.value = all;
}

await loadJobs();
await load();

watch(selectedJobID, load);

async function generateTask() {
  if (!selectedJobID.value) {
    toast.warning("Эхлээд ажлын байр сонгоно уу");
    return;
  }
  isGenerating.value = true;
  try {
    const result = await tasksAPI.generate(selectedJobID.value);
    form.title = result.title;
    form.description = result.description;
    createOpen.value = true;
  } catch {
    toast.error("AI даалгавар үүсгэхэд алдаа гарлаа");
  } finally {
    isGenerating.value = false;
  }
}

async function createTask(values: {
  title: string;
  description?: string;
  due_date?: string;
}) {
  isSubmitting.value = true;
  try {
    await tasksAPI.create({
      job_posting_id: selectedJobID.value!,
      title: values.title,
      description: values.description ?? "",
      due_date: values.due_date || undefined,
    });
    createOpen.value = false;
    form.title = "";
    form.description = "";
    form.due_date = "";
    toast.success("Даалгавар үүсгэгдлээ");
    await load();
  } catch {
    toast.error("Алдаа гарлаа");
  } finally {
    isSubmitting.value = false;
  }
}

async function sendTask(task: Task) {
  try {
    await tasksAPI.send(task.id);
    toast.success("Даалгавар илгээгдлээ");
    await load();
  } catch {
    toast.error("Алдаа гарлаа");
  }
}

async function submitTask() {
  if (!submitContent.value.trim()) {
    toast.warning("Гүйцэтгэлийн агуулга шаардлагатай");
    return;
  }
  isSubmitting.value = true;
  try {
    await tasksAPI.submit(selectedTask.value!.id, submitContent.value);
    submitOpen.value = false;
    submitContent.value = "";
    toast.success("Даалгавар амжилттай илгээгдлээ");
    await load();
  } catch {
    toast.error("Алдаа гарлаа");
  } finally {
    isSubmitting.value = false;
  }
}

function statusLabel(status: Task["status"]) {
  const map: Record<string, string> = {
    draft: "Ноорог",
    sent: "Илгээгдсэн",
    completed: "Гүйцэтгэгдсэн",
    graded: "Үнэлэгдсэн",
  };
  return map[status] ?? status;
}

function statusClass(status: Task["status"]) {
  if (status === "graded") return "badge-success";
  if (status === "completed") return "badge-info";
  if (status === "sent") return "badge-warning";
  return "bg-muted text-muted-foreground";
}

function formatDate(d: string | null | undefined) {
  if (!d) return "";
  return new Intl.DateTimeFormat("mn-MN", { dateStyle: "medium" }).format(new Date(d));
}
</script>

<template>
  <div class="space-y-5">
    <!-- Page header -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <p class="text-[12px] font-semibold uppercase tracking-[0.8px] text-muted-foreground">
          {{ isRecruiter ? "Удирдлага" : "Миний даалгаврууд" }}
        </p>
        <h1 class="mt-1 text-[26px] font-semibold tracking-[-0.6px]">Даалгаврууд</h1>
      </div>
      <div v-if="isRecruiter" class="flex gap-2">
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
          class="inline-flex items-center gap-2 rounded-full px-4 py-2 text-[13.5px] font-semibold text-white transition hover:opacity-90 disabled:opacity-50"
          style="background: linear-gradient(135deg, var(--primary), oklch(0.348 0.106 295))"
          :disabled="!selectedJobID"
          @click="createOpen = true"
        >
          <Plus class="h-4 w-4" /> Гараар үүсгэх
        </button>
      </div>
    </div>

    <!-- Job selector -->
    <div v-if="isRecruiter && jobs.length" class="flex items-center gap-2">
      <label class="text-[13px] font-medium text-muted-foreground">Ажлын байр:</label>
      <select
        v-model="selectedJobID"
        class="rounded-xl border border-input bg-background px-3 py-2 text-[13.5px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        <option :value="null">-- Сонгох --</option>
        <option v-for="job in jobs" :key="job.id" :value="job.id">{{ job.title }}</option>
      </select>
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
      <div class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-muted">
        <ClipboardList class="h-7 w-7 text-muted-foreground" />
      </div>
      <p class="text-[16px] font-semibold">Даалгавар байхгүй</p>
      <p v-if="isRecruiter && !selectedJobID" class="mt-2 text-[13.5px] text-muted-foreground">
        Даалгавар харахын тулд ажлын байр сонгоно уу.
      </p>
    </div>

    <!-- Task cards -->
    <div v-else class="space-y-3.5">
      <div
        v-for="task in tasks"
        :key="task.id"
        class="rounded-2xl border border-border bg-card p-5 transition hover:border-primary/15"
      >
        <div class="flex items-start justify-between gap-4">
          <div class="flex-1 min-w-0 space-y-1">
            <div class="flex items-center gap-2 flex-wrap">
              <h3 class="text-[15px] font-semibold">{{ task.title }}</h3>
              <span
                v-if="task.created_by_ai"
                class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[11px] font-semibold bg-primary/10 text-primary"
              >
                <Sparkles class="h-2.5 w-2.5" /> AI
              </span>
            </div>
            <p class="text-[13.5px] text-muted-foreground line-clamp-2">{{ task.description }}</p>
            <p v-if="task.due_date" class="text-[12px] text-muted-foreground">
              Дуусах огноо: {{ formatDate(task.due_date) }}
            </p>
          </div>
          <div class="flex shrink-0 items-center gap-2">
            <span
              class="inline-flex items-center rounded-full px-3 py-1 text-[12px] font-semibold"
              :class="statusClass(task.status)"
            >
              {{ statusLabel(task.status) }}
            </span>
            <!-- Recruiter: send -->
            <button
              v-if="isRecruiter && task.status === 'draft'"
              class="inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-[12.5px] font-semibold text-white transition hover:opacity-90"
              style="background: linear-gradient(135deg, var(--primary), oklch(0.348 0.106 295))"
              @click="sendTask(task)"
            >
              <Send class="h-3.5 w-3.5" /> Илгээх
            </button>
            <!-- Applicant: submit -->
            <button
              v-if="!isRecruiter && task.status === 'sent'"
              class="inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-[12.5px] font-semibold text-white transition hover:opacity-90"
              style="background: linear-gradient(135deg, var(--primary), oklch(0.348 0.106 295))"
              @click="selectedTask = task; submitOpen = true"
            >
              <ChevronRight class="h-3.5 w-3.5" /> Гүйцэтгэл илгээх
            </button>
          </div>
        </div>

        <!-- Recruiter: submissions -->
        <div v-if="isRecruiter && task.submissions?.length" class="mt-4 space-y-2 border-t border-border pt-4">
          <p class="text-[12px] font-semibold uppercase tracking-[0.5px] text-muted-foreground">
            Илгээлтүүд ({{ task.submissions.length }})
          </p>
          <div
            v-for="sub in task.submissions"
            :key="sub.id"
            class="rounded-xl bg-muted/40 px-4 py-3 text-[13.5px]"
          >
            <p class="line-clamp-2 text-foreground">{{ sub.content }}</p>
            <p v-if="sub.grade !== null" class="mt-1.5 flex items-center gap-1.5 text-[12px] font-semibold text-success-foreground">
              <span class="h-1.5 w-1.5 rounded-full bg-success" />
              Оноо: {{ sub.grade }}/100
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Create task dialog -->
  <Dialog v-model:open="createOpen">
    <DialogContent class="rounded-2xl sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>Даалгавар үүсгэх</DialogTitle>
      </DialogHeader>
      <Form
        :validation-schema="taskFormSchema"
        :initial-values="form"
        class="space-y-4 py-2"
        @submit="createTask"
      >
        <FormField v-slot="{ componentField }" name="title">
          <FormItem class="space-y-2">
            <FormLabel>Гарчиг</FormLabel>
            <FormControl>
              <Input v-bind="componentField" placeholder="Даалгаврын гарчиг" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="description">
          <FormItem class="space-y-2">
            <FormLabel>Тайлбар</FormLabel>
            <FormControl>
              <textarea
                v-bind="componentField"
                rows="5"
                placeholder="Даалгаврын дэлгэрэнгүй тайлбар..."
                class="w-full resize-none rounded-xl border border-input bg-background px-3 py-2 text-[13.5px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="due_date">
          <FormItem class="space-y-2">
            <FormLabel>
              Дуусах огноо
              <span class="text-[11.5px] text-muted-foreground">(заавал биш)</span>
            </FormLabel>
            <FormControl>
              <Input v-bind="componentField" type="date" />
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

  <!-- Submit task dialog -->
  <Dialog v-model:open="submitOpen">
    <DialogContent class="rounded-2xl sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>{{ selectedTask?.title }}</DialogTitle>
        <DialogDescription>{{ selectedTask?.description }}</DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-2">
        <div class="space-y-2">
          <Label>Гүйцэтгэлийн агуулга</Label>
          <textarea
            v-model="submitContent"
            rows="6"
            placeholder="Гүйцэтгэсэн ажлынхаа тайлбар, хариулт эсвэл холбоосыг энд бичнэ үү..."
            class="w-full resize-none rounded-xl border border-input bg-background px-3 py-2 text-[13.5px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
        </div>
      </div>
      <DialogFooter>
        <Button variant="outline" @click="submitOpen = false" :disabled="isSubmitting">Болих</Button>
        <Button @click="submitTask" :disabled="isSubmitting">
          <Loader2 v-if="isSubmitting" class="mr-2 h-4 w-4 animate-spin" />
          Илгээх
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
