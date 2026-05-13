<script setup lang="ts">
import type { AssessmentResult } from "../../composables/types";
import { CheckCircle2, XCircle, MinusCircle } from "lucide-vue-next";

const props = defineProps<{
  result: AssessmentResult;
  /** "preview" = analyze mode, "submitted" = after applying */
  mode?: "preview" | "submitted";
}>();

function scoreColor(score: number) {
  if (score >= 75) return "text-green-600";
  if (score >= 50) return "text-yellow-600";
  return "text-red-600";
}

function scoreBorder(score: number) {
  if (score >= 75) return "border-green-200 bg-green-50";
  if (score >= 50) return "border-yellow-200 bg-yellow-50";
  return "border-red-200 bg-red-50";
}

function statusIcon(status: string) {
  if (status === "met") return CheckCircle2;
  if (status === "partial") return MinusCircle;
  return XCircle;
}

function statusColor(status: string) {
  if (status === "met") return "text-green-600";
  if (status === "partial") return "text-yellow-600";
  return "text-red-500";
}

function statusBg(status: string) {
  if (status === "met") return "bg-green-50 border-green-100";
  if (status === "partial") return "bg-yellow-50 border-yellow-100";
  return "bg-red-50 border-red-100";
}
</script>

<template>
  <!-- Score header -->
  <div :class="['rounded-2xl border px-5 py-4', scoreBorder(result.overall_score)]">
    <div class="flex items-center gap-4">
      <div>
        <p class="text-xs font-medium uppercase tracking-wide text-muted-foreground">
          {{ mode === 'preview' ? 'Урьдчилсан үнэлгээ' : 'Таны үнэлгээний үр дүн' }}
        </p>
        <p :class="['text-4xl font-bold leading-none mt-1', scoreColor(result.overall_score)]">
          {{ result.overall_score }}%
        </p>
      </div>
      <p class="flex-1 text-sm leading-6 text-muted-foreground">{{ result.summary }}</p>
    </div>
  </div>
</template>
