<script setup lang="ts">
import type { TaskSubmission } from "../composables/types";
import { toast } from "vue-sonner";
import {
  ClipboardCheck,
  Loader2,
  Sparkles,
  Pencil,
  ChevronRight,
  ChevronLeft,
  CheckCircle2,
  Clock,
} from "lucide-vue-next";

definePageMeta({ middleware: "auth", fullscreen: true });

const tasksAPI = useTasksAPI();
const { user } = useAuth();
const config = useRuntimeConfig();

const isRecruiter = computed(() => user.value?.role === "recruiter");
if (!isRecruiter.value) navigateTo("/");

const submissions = ref<TaskSubmission[]>([]);
const loading = ref(true);
const selected = ref<TaskSubmission | null>(null);

const subPage = ref(1);
const SUB_PAGE_SIZE = 10;
const subTotalPages = computed(() =>
  Math.max(1, Math.ceil(submissions.value.length / SUB_PAGE_SIZE)),
);
const paginatedSubmissions = computed(() => {
  const start = (subPage.value - 1) * SUB_PAGE_SIZE;
  return submissions.value.slice(start, start + SUB_PAGE_SIZE);
});
const pdfBlobUrl = ref<string | null>(null);
const pdfLoading = ref(false);

async function loadPdf(sub: TaskSubmission) {
  if (pdfBlobUrl.value) {
    URL.revokeObjectURL(pdfBlobUrl.value);
    pdfBlobUrl.value = null;
  }
  if (!sub.has_file) return;
  pdfLoading.value = true;
  try {
    const res = await fetch(
      `${config.public.apiBase}/api/tasks/submissions/${sub.id}/file`,
      { credentials: "include" },
    );
    if (!res.ok) throw new Error("fetch failed");
    const blob = await res.blob();
    pdfBlobUrl.value = URL.createObjectURL(blob);
  } catch {
    toast.error("PDF ачааллаж чадсангүй");
  } finally {
    pdfLoading.value = false;
  }
}

onUnmounted(() => {
  if (pdfBlobUrl.value) URL.revokeObjectURL(pdfBlobUrl.value);
});

const gradeOpen = ref(false);
const gradeValue = ref("");
const gradeFeedback = ref("");
const isGrading = ref(false);
const isAIGrading = ref(false);

async function load() {
  loading.value = true;
  try {
    submissions.value = await tasksAPI.listAllSubmissions();
  } catch {
    toast.error("Мэдээлэл ачааллаж чадсангүй");
  } finally {
    loading.value = false;
  }
}

await load();

function select(sub: TaskSubmission) {
  selected.value = sub;
  loadPdf(sub);
}

function openGrade() {
  if (!selected.value) return;
  gradeValue.value =
    selected.value.grade !== null ? String(selected.value.grade) : "";
  gradeFeedback.value = selected.value.feedback ?? "";
  gradeOpen.value = true;
}

async function saveGrade() {
  if (!selected.value) return;
  const g = Number(gradeValue.value);
  if (!gradeValue.value || isNaN(g) || g < 0 || g > 100) {
    toast.warning("0-100 хооронд оноо оруулна уу");
    return;
  }
  isGrading.value = true;
  try {
    const updated = await tasksAPI.grade(
      selected.value.id,
      g,
      gradeFeedback.value,
    );
    gradeOpen.value = false;
    toast.success("Үнэлгээ хадгалагдлаа");
    await load();
    selected.value =
      submissions.value.find((s) => s.id === updated.id) ?? updated;
  } catch {
    toast.error("Алдаа гарлаа");
  } finally {
    isGrading.value = false;
  }
}

async function aiGrade() {
  if (!selected.value) return;
  isAIGrading.value = true;
  try {
    const updated = await tasksAPI.aiGrade(selected.value.id);
    toast.success("AI үнэлгээ хийлээ");
    await load();
    selected.value =
      submissions.value.find((s) => s.id === updated.id) ?? updated;
  } catch (e: any) {
    toast.error(e?.data?.message ?? "AI үнэлэхэд алдаа гарлаа");
  } finally {
    isAIGrading.value = false;
  }
}

function formatDate(d: string) {
  return new Intl.DateTimeFormat("mn-MN", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(d));
}
</script>

<template>
  <div class="flex h-[calc(100vh-64px)] gap-0 overflow-hidden">
    <!-- ── Left list panel ─────────────────────────────────────── -->
    <div
      class="flex w-[300px] shrink-0 flex-col border-r border-border bg-card"
    >
      <!-- <div class="border-b border-border px-4 py-4">
        <p class="mt-0.5 text-[12px] text-muted-foreground">
          {{ submissions.length }} илгээлт
        </p>
      </div> -->

      <div v-if="loading" class="flex flex-1 items-center justify-center">
        <Loader2 class="h-5 w-5 animate-spin text-muted-foreground" />
      </div>

      <div
        v-else-if="!submissions.length"
        class="flex flex-1 flex-col items-center justify-center gap-2 px-6 text-center"
      >
        <ClipboardCheck class="h-8 w-8 text-muted-foreground/40" />
        <p class="text-[12.5px] text-muted-foreground">Илгээлт байхгүй</p>
      </div>

      <div v-else class="flex flex-col flex-1 overflow-hidden">
        <div class="flex-1 overflow-y-auto">
          <button
            v-for="sub in paginatedSubmissions"
            :key="sub.id"
            class="w-full border-b border-border/60 px-4 py-3 text-left transition hover:bg-muted/50"
            :class="selected?.id === sub.id ? 'bg-primary/5' : ''"
            @click="select(sub)"
          >
            <div class="flex items-start justify-between gap-2">
              <div class="min-w-0 flex-1">
                <p class="truncate text-[12px] font-medium text-primary/70">
                  {{ sub.task_title }}
                </p>
                <p class="truncate text-[13.5px] font-semibold text-foreground">
                  {{ sub.applicant_name }}
                </p>
                <p class="mt-0.5 text-[11.5px] text-muted-foreground">
                  {{ formatDate(sub.created_at) }}
                </p>
              </div>
              <div class="flex shrink-0 flex-col items-end gap-1 pt-0.5">
                <span
                  :class="[
                    'inline-flex items-center gap-1 rounded-full px-1.5 py-0.5 text-[10.5px] font-semibold',
                    sub.status === 'graded'
                      ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                      : 'bg-amber-500/10 text-amber-600 dark:text-amber-400',
                  ]"
                >
                  <CheckCircle2
                    v-if="sub.status === 'graded'"
                    class="h-2.5 w-2.5"
                  />
                  <Clock v-else class="h-2.5 w-2.5" />
                  {{
                    sub.status === "graded" ? "Үнэлэгдсэн" : "Хүлээгдэж байна"
                  }}
                </span>
                <span
                  v-if="sub.grade !== null"
                  class="text-[11px] font-semibold text-emerald-600 dark:text-emerald-400"
                >
                  {{ sub.grade }}/100
                </span>
              </div>
            </div>
          </button>
        </div>
        <!-- Sidebar pagination -->
        <div
          class="flex items-center justify-center gap-1 border-t border-border py-2"
        >
          <button
            class="flex h-7 w-7 items-center justify-center rounded-full text-muted-foreground transition hover:bg-muted disabled:opacity-40"
            :disabled="subPage === 1"
            @click="subPage--"
          >
            <ChevronLeft class="h-3.5 w-3.5" />
          </button>
          <span class="text-[12px] text-muted-foreground"
            >{{ subPage }} / {{ subTotalPages }}</span
          >
          <button
            class="flex h-7 w-7 items-center justify-center rounded-full text-muted-foreground transition hover:bg-muted disabled:opacity-40"
            :disabled="subPage === subTotalPages"
            @click="subPage++"
          >
            <ChevronRight class="h-3.5 w-3.5" />
          </button>
        </div>
      </div>
    </div>

    <!-- ── Right detail panel ──────────────────────────────────── -->
    <div class="flex min-w-0 flex-1 flex-col bg-background">
      <!-- Empty state -->
      <div
        v-if="!selected"
        class="flex flex-1 flex-col items-center justify-center gap-3 text-center"
      >
        <ChevronRight class="h-10 w-10 text-muted-foreground/30" />
        <p class="text-[14px] font-medium text-muted-foreground">
          Жагсаалтаас илгээлт сонгоно уу
        </p>
      </div>

      <template v-else>
        <!-- Two-column content -->
        <div class="flex min-h-0 flex-1 overflow-hidden">
          <!-- Left col: PDF + answers -->
          <div
            class="flex flex-1 min-w-0 flex-col overflow-y-auto border-r border-border p-5"
          >
            <!-- PDF embed -->
            <div v-if="selected.has_file" class="mb-5">
              <p
                class="mb-2 text-[11.5px] font-semibold uppercase tracking-[0.5px] text-muted-foreground"
              >
                Хавсаргасан файл
              </p>
              <div
                v-if="pdfLoading"
                class="flex h-[480px] w-full items-center justify-center rounded-xl border border-border bg-muted"
              >
                <Loader2 class="h-5 w-5 animate-spin text-muted-foreground" />
              </div>
              <iframe
                v-else-if="pdfBlobUrl"
                :src="pdfBlobUrl"
                class="h-[600px] w-full rounded-xl border border-border bg-muted"
              />
            </div>

            <!-- Text answer -->
            <div v-if="selected.content">
              <p
                class="mb-2 text-[11.5px] font-semibold uppercase tracking-[0.5px] text-muted-foreground"
              >
                Хариулт
              </p>
              <div
                class="rounded-xl border border-border bg-card p-4 text-[13.5px] leading-relaxed whitespace-pre-wrap"
              >
                {{ selected.content }}
              </div>
            </div>

            <div
              v-if="!selected.has_file && !selected.content"
              class="flex flex-1 items-center justify-center text-[13px] text-muted-foreground"
            >
              Агуулга байхгүй
            </div>
          </div>

          <!-- Right col: Task details -->
          <div class="w-[350px] shrink-0 overflow-y-auto p-5">
            <p
              class="mb-3 text-[11.5px] font-semibold uppercase tracking-[0.5px] text-muted-foreground"
            >
              Даалгаврын дэлгэрэнгүй
            </p>
            <!-- Action bar -->

            <div class="flex shrink-0 items-center gap-2">
              <!-- <button
                class="inline-flex items-center gap-1.5 rounded-lg border border-primary/40 bg-primary/5 px-3 py-1.5 text-[12.5px] font-medium text-primary transition hover:bg-primary/10 disabled:opacity-50"
                :disabled="isAIGrading"
                @click="aiGrade"
              >
                <Loader2 v-if="isAIGrading" class="h-3.5 w-3.5 animate-spin" />
                <Sparkles v-else class="h-3.5 w-3.5" />
                AI-аар үнэлэх
              </button> -->
            </div>

            <div class="space-y-4">
              <button
                class="inline-flex items-center gap-1.5 rounded-lg border border-border bg-background px-3 py-1.5 text-[12.5px] font-medium text-foreground transition hover:bg-muted"
                @click="openGrade"
              >
                <Pencil class="h-3.5 w-3.5" />
                Үнэлэх
              </button>
              <div>
                <p class="text-[11px] text-muted-foreground">Илгээсэн огноо</p>
                <p class="mt-0.5 text-[13px]">
                  {{ formatDate(selected.created_at) }}
                </p>
              </div>

              <div>
                <p class="text-[11px] text-muted-foreground">Горилогч</p>
                <p class="mt-0.5 text-[13px] font-medium">
                  {{ selected.applicant_name }}
                </p>
              </div>

              <!-- Grade result -->
              <div
                v-if="selected.grade !== null"
                class="rounded-xl border border-emerald-500/20 bg-emerald-500/5 p-3"
              >
                <p
                  class="text-[11px] font-semibold text-emerald-700 dark:text-emerald-400"
                >
                  Үнэлгээний дүн
                </p>
                <p
                  class="mt-1 text-[22px] font-bold text-emerald-600 dark:text-emerald-400"
                >
                  {{ selected.grade }}
                  <span
                    class="text-[13px] font-normal text-emerald-700/60 dark:text-emerald-500"
                    >/100</span
                  >
                </p>
                <p
                  v-if="selected.feedback"
                  class="mt-2 text-[12px] text-emerald-700/70 dark:text-emerald-400/70 whitespace-pre-wrap"
                >
                  {{ selected.feedback }}
                </p>
              </div>

              <div
                v-else
                class="rounded-xl border border-amber-500/20 bg-amber-500/5 p-3"
              >
                <p
                  class="text-[11px] font-semibold text-amber-700 dark:text-amber-400"
                >
                  Статус
                </p>
                <p
                  class="mt-0.5 text-[13px] text-amber-600 dark:text-amber-400"
                >
                  Үнэлэгдэх хүлээгдэж байна
                </p>
              </div>
              <div>
                <p class="text-[11px] text-muted-foreground">Гарчиг</p>
                <p class="mt-0.5 text-[13.5px] font-semibold">
                  {{ selected.task_title }}
                </p>
              </div>

              <div v-if="selected.task_description">
                <p class="text-[11px] text-muted-foreground">Тайлбар</p>
                <p
                  class="mt-0.5 text-[13px] leading-relaxed text-foreground/80 whitespace-pre-wrap"
                >
                  {{ selected.task_description }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>

  <!-- Manual grade dialog -->
  <Dialog v-model:open="gradeOpen">
    <DialogContent class="rounded-2xl sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Үнэлгээ өгөх</DialogTitle>
        <DialogDescription>
          <span class="font-medium text-foreground">{{
            selected?.applicant_name
          }}</span>
          — {{ selected?.task_title }}
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-2">
        <div class="space-y-2">
          <Label>Оноо (0–100)</Label>
          <Input
            v-model="gradeValue"
            type="number"
            min="0"
            max="100"
            placeholder="Жишээ: 85"
          />
        </div>
        <div class="space-y-2">
          <Label>
            Санал хүсэлт
            <span class="text-[11.5px] text-muted-foreground"
              >(заавал биш)</span
            >
          </Label>
          <textarea
            v-model="gradeFeedback"
            rows="4"
            placeholder="Горилогчид өгөх санал хүсэлт..."
            class="w-full resize-none rounded-xl border border-input bg-background px-3 py-2 text-[13.5px] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
        </div>
      </div>
      <DialogFooter>
        <Button
          variant="outline"
          :disabled="isGrading"
          @click="gradeOpen = false"
        >
          Болих
        </Button>
        <Button :disabled="isGrading" @click="saveGrade">
          <Loader2 v-if="isGrading" class="mr-2 h-4 w-4 animate-spin" />
          Хадгалах
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
