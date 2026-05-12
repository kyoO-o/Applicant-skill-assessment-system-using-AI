<script setup lang="ts">
import type { Task } from "../composables/types";
import { toast } from "vue-sonner";
import { ClipboardList, Plus, Wand2, Send, Loader2, ChevronRight } from "lucide-vue-next";

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

async function createTask() {
  if (!form.title.trim() || !form.description.trim()) {
    toast.warning("Гарчиг болон тайлбар шаардлагатай");
    return;
  }
  isSubmitting.value = true;
  try {
    await tasksAPI.create({
      job_posting_id: selectedJobID.value!,
      title: form.title,
      description: form.description,
      due_date: form.due_date || undefined,
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

function statusVariant(status: Task["status"]) {
  if (status === "graded") return "default";
  if (status === "completed") return "secondary";
  if (status === "sent") return "outline";
  return "outline";
}

function formatDate(d: string | null | undefined) {
  if (!d) return "";
  return new Intl.DateTimeFormat("mn-MN", { dateStyle: "medium" }).format(new Date(d));
}
</script>

<template>
  <div class="space-y-6">
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-sm font-medium text-muted-foreground">
            {{ isRecruiter ? "Даалгаврын удирдлага" : "Миний даалгаврууд" }}
          </p>
          <h2 class="mt-2 text-3xl font-semibold tracking-tight">Даалгаврууд</h2>
        </div>
        <div v-if="isRecruiter" class="flex gap-2">
          <Button variant="outline" class="rounded-full" :disabled="isGenerating" @click="generateTask">
            <Wand2 class="mr-2 h-4 w-4" />
            {{ isGenerating ? "Үүсгэж байна..." : "AI-аар үүсгэх" }}
          </Button>
          <Button class="rounded-full" :disabled="!selectedJobID" @click="createOpen = true">
            <Plus class="mr-2 h-4 w-4" /> Гараар үүсгэх
          </Button>
        </div>
      </div>

      <!-- Job selector for recruiter -->
      <div v-if="isRecruiter && jobs.length" class="mt-4">
        <select
          v-model="selectedJobID"
          class="rounded-xl border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <option :value="null">-- Ажлын байр сонгох --</option>
          <option v-for="job in jobs" :key="job.id" :value="job.id">{{ job.title }}</option>
        </select>
      </div>
    </section>

    <div v-if="loading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <div v-else-if="!tasks.length"
         class="rounded-3xl border border-dashed border-border px-6 py-20 text-center">
      <ClipboardList class="mx-auto h-10 w-10 text-muted-foreground" />
      <p class="mt-4 text-lg font-medium">Даалгавар байхгүй</p>
    </div>

    <div v-else class="space-y-4">
      <div v-for="task in tasks" :key="task.id"
           class="rounded-3xl border border-border bg-card px-6 py-5 shadow-sm">
        <div class="flex items-start justify-between gap-4">
          <div class="space-y-1 min-w-0">
            <div class="flex items-center gap-2">
              <h3 class="font-semibold">{{ task.title }}</h3>
              <Badge v-if="task.created_by_ai" variant="secondary" class="rounded-full px-2 text-xs">AI</Badge>
            </div>
            <p class="text-sm text-muted-foreground line-clamp-2">{{ task.description }}</p>
            <p v-if="task.due_date" class="text-xs text-muted-foreground">
              Дуусах огноо: {{ formatDate(task.due_date) }}
            </p>
          </div>
          <div class="flex shrink-0 items-center gap-2">
            <Badge :variant="statusVariant(task.status) as any" class="rounded-full px-3">
              {{ statusLabel(task.status) }}
            </Badge>
            <!-- Recruiter: send button -->
            <Button v-if="isRecruiter && task.status === 'draft'" size="sm" class="rounded-full" @click="sendTask(task)">
              <Send class="mr-1 h-3 w-3" /> Илгээх
            </Button>
            <!-- Applicant: submit button -->
            <Button v-if="!isRecruiter && task.status === 'sent'" size="sm" class="rounded-full"
                    @click="selectedTask = task; submitOpen = true">
              <ChevronRight class="mr-1 h-3 w-3" /> Гүйцэтгэл илгээх
            </Button>
          </div>
        </div>

        <!-- Submissions for recruiter -->
        <div v-if="isRecruiter && task.submissions?.length" class="mt-3 border-t border-border pt-3">
          <p class="mb-2 text-xs font-medium text-muted-foreground">Илгээлтүүд ({{ task.submissions.length }})</p>
          <div v-for="sub in task.submissions" :key="sub.id" class="rounded-2xl bg-muted px-3 py-2 text-sm">
            <p class="line-clamp-2">{{ sub.content }}</p>
            <p v-if="sub.grade !== null" class="mt-1 text-xs text-muted-foreground">Оноо: {{ sub.grade }}/100</p>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Create task dialog -->
  <Dialog v-model:open="createOpen">
    <DialogContent class="rounded-3xl sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>Даалгавар үүсгэх</DialogTitle>
      </DialogHeader>
      <div class="space-y-4 py-2">
        <div class="space-y-2">
          <Label>Гарчиг</Label>
          <Input v-model="form.title" placeholder="Даалгаврын гарчиг" />
        </div>
        <div class="space-y-2">
          <Label>Тайлбар</Label>
          <textarea
            v-model="form.description"
            rows="5"
            placeholder="Даалгаврын дэлгэрэнгүй тайлбар..."
            class="w-full resize-none rounded-xl border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
        </div>
        <div class="space-y-2">
          <Label>Дуусах огноо (заавал биш)</Label>
          <Input v-model="form.due_date" type="date" />
        </div>
      </div>
      <DialogFooter>
        <Button variant="outline" @click="createOpen = false" :disabled="isSubmitting">Болих</Button>
        <Button @click="createTask" :disabled="isSubmitting">
          <Loader2 v-if="isSubmitting" class="mr-2 h-4 w-4 animate-spin" />
          Хадгалах
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>

  <!-- Submit task dialog -->
  <Dialog v-model:open="submitOpen">
    <DialogContent class="rounded-3xl sm:max-w-lg">
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
            class="w-full resize-none rounded-xl border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
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
