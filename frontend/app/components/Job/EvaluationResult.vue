<script setup lang="ts">
import type { AssessmentResult } from "../../composables/types";
import { computed } from "vue";
import {
  CheckCircle2,
  XCircle,
  AlertCircle,
  Sparkles,
  RefreshCw,
  FileText,
  ArrowRight,
  Pencil,
} from "lucide-vue-next";

const props = defineProps<{
  result: AssessmentResult;
  mode?: "preview" | "submitted";
  applicationStatus?: string;
}>();

const emit = defineEmits<{
  apply: [];
  reloadCv: [];
}>();

const isPass = computed(() => props.result.overall_score >= 75);

// Gauge geometry
const GAUGE_SIZE = 144;
const GAUGE_STROKE = 12;
const gaugeR = (GAUGE_SIZE - GAUGE_STROKE) / 2;
const gaugeC = 2 * Math.PI * gaugeR;
const gaugeDash = computed(() => (props.result.overall_score / 100) * gaugeC);
const gaugeId = `evr-${Math.random().toString(36).slice(2, 8)}`;

// KPI stats
const matchedCount = computed(() => props.result.matched_skills?.length ?? 0);
const totalSkills = computed(
  () => matchedCount.value + (props.result.missing_skills?.length ?? 0),
);
const metDuties = computed(
  () =>
    props.result.duty_assessments?.filter((d) => d.status === "met").length ??
    0,
);
const totalDuties = computed(() => props.result.duty_assessments?.length ?? 0);
const metReqs = computed(
  () =>
    props.result.requirement_assessments?.filter((r) => r.status === "met")
      .length ?? 0,
);
const totalReqs = computed(
  () => props.result.requirement_assessments?.length ?? 0,
);

const showApplyBtn = computed(() => props.mode === "preview" && isPass.value);
</script>

<template>
  <!-- ── Main evaluation card ──────────────────────────────────── -->
  <div
    class="rounded-2xl border overflow-hidden"
    :class="isPass ? 'border-success/20' : 'border-border'"
    :style="{
      backgroundImage: isPass
        ? 'radial-gradient(120% 70% at 50% 0%, rgba(16,185,129,0.07), transparent 58%)'
        : 'radial-gradient(120% 70% at 50% 0%, rgba(107,81,145,0.09), transparent 58%)',
    }"
  >
    <div class="bg-card/85 p-5 sm:p-6">
      <!-- Top row: gauge + text -->
      <div class="flex flex-col sm:flex-row sm:items-center gap-5 sm:gap-7">
        <!-- SVG Gauge (desktop only) -->
        <div
          class="hidden sm:block flex-shrink-0"
          :style="`width:${GAUGE_SIZE}px;height:${GAUGE_SIZE}px;position:relative`"
        >
          <svg
            :width="GAUGE_SIZE"
            :height="GAUGE_SIZE"
            :viewBox="`0 0 ${GAUGE_SIZE} ${GAUGE_SIZE}`"
            style="transform: rotate(-90deg)"
            aria-hidden="true"
          >
            <defs>
              <linearGradient :id="gaugeId" x1="0" y1="0" x2="1" y2="1">
                <stop
                  offset="0%"
                  :stop-color="isPass ? '#10b981' : '#a8a29e'"
                />
                <stop
                  offset="100%"
                  :stop-color="isPass ? '#059669' : '#78716c'"
                />
              </linearGradient>
            </defs>
            <!-- Track -->
            <circle
              :cx="GAUGE_SIZE / 2"
              :cy="GAUGE_SIZE / 2"
              :r="gaugeR"
              fill="none"
              stroke="var(--secondary)"
              :stroke-width="GAUGE_STROKE"
            />
            <!-- Fill -->
            <circle
              :cx="GAUGE_SIZE / 2"
              :cy="GAUGE_SIZE / 2"
              :r="gaugeR"
              fill="none"
              :stroke="`url(#${gaugeId})`"
              :stroke-width="GAUGE_STROKE"
              :stroke-dasharray="`${gaugeDash} ${gaugeC - gaugeDash}`"
              stroke-linecap="round"
              style="
                transition: stroke-dasharray 0.8s cubic-bezier(0.2, 0.7, 0.3, 1);
              "
            />
          </svg>
          <!-- Score label -->
          <div
            class="absolute inset-0 grid place-items-center text-center pointer-events-none"
          >
            <div>
              <div
                class="text-[41px] font-semibold leading-none tracking-[-2px]"
                :class="
                  isPass ? 'text-success-foreground' : 'text-muted-foreground'
                "
              >
                {{ result.overall_score }}
              </div>
              <div
                class="text-[10px] font-medium uppercase tracking-[0.3px] text-muted-foreground mt-1 font-mono"
              >
                100-аас
              </div>
            </div>
          </div>
        </div>

        <!-- Text content -->
        <div class="flex-1 min-w-0">
          <!-- AI label + verdict pill -->
          <div class="flex flex-wrap items-center gap-2 mb-3">
            <div
              class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[10.5px] font-semibold tracking-[0.4px] border"
              style="
                background: linear-gradient(
                  135deg,
                  rgba(107, 81, 145, 0.12),
                  rgba(87, 64, 119, 0.12)
                );
                color: var(--primary);
                border-color: rgba(107, 81, 145, 0.2);
              "
            >
              <Sparkles class="h-2.5 w-2.5" />
              AI үнэлгээ
            </div>
          </div>

          <!-- Mobile: score numerals -->
          <div class="sm:hidden flex items-baseline gap-2 mb-3">
            <span
              class="text-[38px] font-semibold leading-none tracking-[-1.5px]"
              :class="
                isPass ? 'text-success-foreground' : 'text-muted-foreground'
              "
            >
              {{ result.overall_score }}
            </span>
            <span class="text-[16px] text-muted-foreground font-medium"
              >/ 100</span
            >
          </div>

          <!-- Headline -->
          <h2 class="text-[22px] font-semibold tracking-[-0.5px] leading-[1.2]">
            Таны үнэлгээний хариу
          </h2>

          <!-- Summary -->
          <p
            class="text-[13.5px] text-muted-foreground mt-2 leading-[1.55] max-w-2xl"
          >
            {{ result.summary }}
          </p>

          <!-- Action buttons -->
          <div class="flex flex-wrap gap-2.5 mt-4">
            <!-- Pass CTAs -->
            <button
              class="inline-flex items-center gap-2 px-4 py-2 rounded-full text-[13.5px] font-semibold text-white cursor-pointer border-0"
              style="
                background: linear-gradient(
                  135deg,
                  var(--primary),
                  oklch(0.348 0.106 295)
                );
              "
              @click="emit('apply')"
            >
              {{ mode == "preview" ? "Анкет илгээх" : "Дахин илгээх" }}
              <ArrowRight class="h-3.5 w-3.5" />
            </button>
            <!-- <button
                class="inline-flex items-center gap-2 px-4 py-2 rounded-full text-[13.5px] font-semibold border border-border bg-card text-foreground cursor-pointer"
              >
                <FileText class="h-3.5 w-3.5" /> Дэлгэрэнгүй тайлан
              </button> -->
            <button
              class="inline-flex items-center gap-2 px-4 py-2 rounded-full text-[13.5px] font-semibold text-foreground cursor-pointer border-0 bg-transparent"
              @click="emit('reloadCv')"
            >
              <RefreshCw class="h-3.5 w-3.5" /> CV дахин ачаалах
            </button>

            <!-- Fail CTAs -->
            <!-- <template v-else>
              <button
                class="inline-flex items-center gap-2 px-4 py-2 rounded-full text-[13.5px] font-semibold text-white cursor-pointer border-0"
                style="
                  background: linear-gradient(
                    135deg,
                    var(--primary),
                    oklch(0.348 0.106 295)
                  );
                "
              >
                Төстэй ажил үзэх <ArrowRight class="h-3.5 w-3.5" />
              </button>
              <button
                class="inline-flex items-center gap-2 px-4 py-2 rounded-full text-[13.5px] font-semibold border border-border bg-card text-foreground cursor-pointer"
                @click="emit('reloadCv')"
              >
                <Pencil class="h-3.5 w-3.5" /> CV-гээ сайжруулах
              </button>
            </template> -->
          </div>
        </div>
      </div>

      <!-- 4-column KPI footer -->
      <div
        class="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-6 pt-5 border-t border-border"
      >
        <div>
          <p
            class="text-[10.5px] font-medium uppercase tracking-[0.6px] text-muted-foreground font-mono"
          >
            Чадварууд
          </p>
          <p
            class="text-[18px] font-semibold mt-1.5 tracking-[-0.3px] text-secondary-foreground"
          >
            {{ matchedCount }} / {{ totalSkills }}
          </p>
        </div>
        <!-- <div>
          <p
            class="text-[10.5px] font-medium uppercase tracking-[0.6px] text-muted-foreground font-mono"
          >
            Үүрэг хариуцлага
          </p>
          <p
            class="text-[18px] font-semibold mt-1.5 tracking-[-0.3px]"
            :class="isPass ? 'text-success-foreground' : 'text-destructive'"
          >
            {{ metDuties }} / {{ totalDuties }}
          </p>
        </div> -->
        <div>
          <p
            class="text-[10.5px] font-medium uppercase tracking-[0.6px] text-muted-foreground font-mono"
          >
            Шаардлагууд
          </p>
          <p
            class="text-[18px] font-semibold mt-1.5 tracking-[-0.3px] text-secondary-foreground"
          >
            {{ metReqs }} / {{ totalReqs }}
          </p>
        </div>
        <div>
          <p
            class="text-[10.5px] font-medium uppercase tracking-[0.6px] text-muted-foreground font-mono"
          >
            Нийт оноо
          </p>
          <p
            class="text-[18px] font-semibold mt-1.5 tracking-[-0.3px] text-secondary-foreground"
          >
            {{ result.overall_score }}%
          </p>
        </div>
      </div>
    </div>
  </div>

  <!-- ── Strengths + Gaps cards ────────────────────────────────── -->
  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-4">
    <!-- Strengths -->
    <div class="rounded-2xl border border-border bg-card p-5">
      <div class="flex items-center gap-2 mb-3">
        <div
          class="w-6 h-6 rounded-lg bg-success-bg text-success-foreground grid place-items-center flex-shrink-0"
        >
          <CheckCircle2 class="h-3.5 w-3.5" />
        </div>
        <span class="text-[14px] font-semibold">Хүчтэй талууд</span>
      </div>
      <div class="space-y-2.5">
        <div
          v-for="s in (result.matched_skills ?? []).slice(0, 5)"
          :key="s.skill"
          class="flex items-start gap-2 text-[13px] leading-[1.5]"
        >
          <div
            class="w-4 h-4 rounded-full bg-success-bg text-success-foreground grid place-items-center flex-shrink-0 mt-0.5"
          >
            <CheckCircle2 class="h-2.5 w-2.5" />
          </div>
          <div>
            <span class="font-medium">{{ s.skill }}</span
            ><span v-if="s.explanation" class="text-muted-foreground">
              — {{ s.explanation }}</span
            >
          </div>
        </div>
        <p
          v-if="!result.matched_skills?.length"
          class="text-[13px] text-muted-foreground italic"
        >
          Тохирсон чадвар дурдагдаагүй
        </p>
      </div>
    </div>

    <!-- Gaps / Warnings -->
    <div class="rounded-2xl border border-border bg-card p-5">
      <div class="flex items-center gap-2 mb-3">
        <div
          class="w-6 h-6 rounded-lg grid place-items-center flex-shrink-0"
          :class="
            isPass
              ? 'bg-warning-bg text-warning-foreground'
              : 'bg-destructive/10 text-destructive'
          "
        >
          <AlertCircle v-if="isPass" class="h-3.5 w-3.5" />
          <XCircle v-else class="h-3.5 w-3.5" />
        </div>
        <span class="text-[14px] font-semibold">Анхаарах талууд</span>
      </div>
      <div class="space-y-2.5">
        <div
          v-for="s in (result.missing_skills ?? []).slice(0, 5)"
          :key="s.skill"
          class="flex items-start gap-2 text-[13px] leading-[1.5]"
        >
          <div
            class="w-4 h-4 rounded-full grid place-items-center flex-shrink-0 mt-0.5"
            :class="
              isPass
                ? 'bg-warning-bg text-warning-foreground'
                : 'bg-destructive/10 text-destructive'
            "
          >
            <AlertCircle v-if="isPass" class="h-2.5 w-2.5" />
            <XCircle v-else class="h-2.5 w-2.5" />
          </div>
          <div>
            <span class="font-medium">{{ s.skill }}</span
            ><span v-if="s.explanation" class="text-muted-foreground">
              — {{ s.explanation }}</span
            >
          </div>
        </div>
        <p
          v-if="!result.missing_skills?.length"
          class="text-[13px] text-muted-foreground italic"
        >
          Дутагдаж буй чадвар дурдагдаагүй
        </p>
      </div>
    </div>
  </div>
</template>
