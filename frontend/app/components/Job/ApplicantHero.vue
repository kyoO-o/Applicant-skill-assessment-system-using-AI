<script setup lang="ts">
import { Upload, Sparkles, FileText, User, CheckCircle2 } from "lucide-vue-next";

const props = defineProps<{
  hasSavedCV: boolean;
}>();

const emit = defineEmits<{
  analyze: [file: File | null, useProfile: boolean];
}>();

type Stage = "idle" | "uploading" | "scoring";
const stage = ref<Stage>("idle");
const uploadProgress = ref(0);
const scoringStep = ref(0);
const fileName = ref("");

const scoringSteps = [
  {
    title: "CV-г PDF-ээс задлаж байна…",
    sub: "Хэлбэрлэлт, хэсэг, огноог уншиж байна",
  },
  {
    title: "Ажлын туршлага, чадваруудыг шинжилж байна…",
    sub: "Туршлагаас мэргэжлийг дүгнэж байна",
  },
  {
    title: "Шаардлагуудтай харьцуулж байна…",
    sub: "CV-г ажлын байрны тодорхойлолттой харьцуулж байна",
  },
  {
    title: "Үүрэг тус бүрийг оноо өгч байна…",
    sub: "Үүрэг тус бүрийн тохирлын дүгнэлт",
  },
  {
    title: "Зөвлөмж бэлтгэж байна…",
    sub: "Өргөдлийг хэрхэн бэхжүүлэх талаар",
  },
];

let progressTimer: ReturnType<typeof setInterval> | null = null;
let stepTimer: ReturnType<typeof setInterval> | null = null;

function startAnimation(name: string) {
  fileName.value = name;
  stage.value = "uploading";
  uploadProgress.value = 0;

  let p = 0;
  progressTimer = setInterval(() => {
    p += 5 + Math.random() * 12;
    if (p >= 100) {
      p = 100;
      clearInterval(progressTimer!);
      progressTimer = null;
      startScoring();
    }
    uploadProgress.value = Math.round(p);
  }, 80);
}

function startScoring() {
  stage.value = "scoring";
  scoringStep.value = 0;
  let i = 0;
  stepTimer = setInterval(() => {
    i++;
    if (i < scoringSteps.length) {
      scoringStep.value = i;
    } else {
      clearInterval(stepTimer!);
      stepTimer = null;
    }
  }, 900);
}

function triggerFile(e: Event) {
  const f = (e.target as HTMLInputElement).files?.[0];
  if (!f) return;
  startAnimation(f.name);
  emit("analyze", f, false);
}

function onDrop(e: DragEvent) {
  e.preventDefault();
  const f = e.dataTransfer?.files?.[0];
  if (!f) return;
  startAnimation(f.name);
  emit("analyze", f, false);
}

function useProfileCV() {
  startAnimation("Хадгалагдсан CV");
  emit("analyze", null, true);
}

onUnmounted(() => {
  if (progressTimer) clearInterval(progressTimer);
  if (stepTimer) clearInterval(stepTimer);
});
</script>

<template>
  <!-- ── Idle: Upload zone ── -->
  <div
    v-if="stage === 'idle'"
    class="relative overflow-hidden rounded-2xl border border-primary/20 bg-card p-6 sm:p-7"
    style="
      background-image: radial-gradient(
        130% 80% at 50% 0%,
        rgba(107, 81, 145, 0.1),
        transparent 60%
      );
    "
  >
    <div class="relative flex flex-col gap-6 sm:flex-row sm:items-center">
      <!-- Left: text -->
      <div class="flex-1 min-w-0">
        <div
          class="inline-flex items-center gap-1.5 rounded-full border border-primary/20 bg-primary/10 px-2.5 py-1 text-[10.5px] font-semibold tracking-[0.4px] text-primary"
        >
          <Sparkles class="h-2.5 w-2.5" />
          AI үнэлгээ
        </div>
        <h2
          class="mt-3 text-[20px] font-semibold tracking-[-0.4px] leading-[1.25]"
        >
          CV-гээ оруулж, ажлын байранд хэр тохирохоо шалга
        </h2>
        <p
          class="mt-2 text-[13.5px] leading-[1.6] text-muted-foreground max-w-md"
        >
          Ажлын байрны шаардлага, чадвар, үүрэгтэй харьцуулж үнэлнэ — үр дүн
          ~15 секундэд гарна.
        </p>
        <button
          v-if="hasSavedCV"
          class="mt-4 inline-flex items-center gap-2 rounded-full border border-border bg-background px-3.5 py-2 text-[13px] font-medium transition hover:bg-muted"
          @click="useProfileCV"
        >
          <User class="h-3.5 w-3.5 text-muted-foreground" />
          Хадгалагдсан CV ашиглах
        </button>
      </div>

      <!-- Right: dropzone -->
      <label
        class="relative flex min-w-[220px] cursor-pointer flex-col items-center gap-3 rounded-2xl border-2 border-dashed border-primary/35 bg-background/60 px-7 py-8 text-center transition hover:border-primary/60 hover:bg-background"
        @dragover.prevent
        @drop="onDrop"
      >
        <input type="file" accept=".pdf" class="sr-only" @change="triggerFile" />
        <div
          class="flex h-12 w-12 items-center justify-center rounded-full text-white shadow-sm"
          style="
            background: linear-gradient(
              135deg,
              var(--primary),
              oklch(0.348 0.106 295)
            );
          "
        >
          <Upload class="h-5 w-5" />
        </div>
        <div>
          <p class="text-[14px] font-semibold">CV энд оруулах</p>
          <p class="mt-0.5 text-[11.5px] text-muted-foreground">
            эсвэл дарж сонгох · PDF, max 5MB
          </p>
        </div>
      </label>
    </div>
  </div>

  <!-- ── Uploading ── -->
  <div
    v-else-if="stage === 'uploading'"
    class="relative overflow-hidden rounded-2xl border border-primary/20 bg-card p-6"
    style="
      background-image: radial-gradient(
        130% 80% at 50% 0%,
        rgba(107, 81, 145, 0.1),
        transparent 60%
      );
    "
  >
    <div class="relative flex items-center gap-5">
      <div
        class="flex h-14 w-14 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary"
      >
        <FileText class="h-6 w-6" />
      </div>
      <div class="flex-1">
        <p class="text-[15px] font-semibold">{{ fileName }}</p>
        <p class="mt-1 text-[12.5px] text-muted-foreground">
          Ачаалж байна… {{ uploadProgress }}%
        </p>
        <div class="mt-2.5 h-1.5 overflow-hidden rounded-full bg-muted">
          <div
            class="h-full rounded-full"
            style="
              background: linear-gradient(
                90deg,
                var(--primary),
                oklch(0.348 0.106 295)
              );
              transition: width 0.1s linear;
            "
            :style="`width: ${uploadProgress}%`"
          />
        </div>
      </div>
    </div>
  </div>

  <!-- ── Scoring ── -->
  <div
    v-else
    class="relative overflow-hidden rounded-2xl border border-primary/20 bg-card p-6"
    style="
      background-image: radial-gradient(
        130% 80% at 50% 0%,
        rgba(107, 81, 145, 0.1),
        transparent 60%
      );
    "
  >
    <div class="relative">
      <!-- Header row -->
      <div class="mb-5 flex items-center gap-4">
        <div class="relative h-12 w-12 shrink-0">
          <div
            class="absolute inset-0 animate-spin rounded-full"
            style="
              background: conic-gradient(
                var(--primary),
                oklch(0.348 0.106 295),
                var(--primary)
              );
              animation-duration: 1.6s;
            "
          />
          <div
            class="absolute inset-[3px] flex items-center justify-center rounded-full bg-card"
          >
            <Sparkles class="h-4 w-4 text-primary" />
          </div>
        </div>
        <div>
          <p class="text-[16px] font-semibold tracking-[-0.3px]">
            Claude таны CV-г шалгаж байна…
          </p>
          <p class="mt-0.5 text-[12.5px] text-muted-foreground">
            Энэ нь ихэвчлэн 8–15 секунд болдог
          </p>
        </div>
      </div>

      <!-- Steps -->
      <div class="space-y-1">
        <div
          v-for="(s, i) in scoringSteps"
          :key="i"
          class="flex items-center gap-3 rounded-xl px-3 py-2.5 transition-all duration-300"
          :class="[i === scoringStep ? 'bg-primary/10' : '', i > scoringStep ? 'opacity-40' : '']"
        >
          <div
            class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg transition-colors"
            :class="
              i < scoringStep
                ? 'bg-success-bg text-success-foreground'
                : i === scoringStep
                  ? 'bg-primary/10 text-primary'
                  : 'bg-muted text-muted-foreground'
            "
          >
            <CheckCircle2 v-if="i < scoringStep" class="h-3.5 w-3.5" />
            <Sparkles v-else class="h-3.5 w-3.5" />
          </div>
          <div class="flex-1">
            <p
              class="text-[13.5px]"
              :class="i === scoringStep ? 'font-semibold' : 'font-medium'"
            >
              {{ s.title }}
            </p>
            <p class="text-[11.5px] text-muted-foreground">{{ s.sub }}</p>
          </div>
          <!-- Bounce dots for active step -->
          <div v-if="i === scoringStep" class="flex gap-1">
            <span
              v-for="d in 3"
              :key="d"
              class="h-1.5 w-1.5 animate-bounce rounded-full bg-primary"
              :style="`animation-delay: ${(d - 1) * 0.15}s`"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
