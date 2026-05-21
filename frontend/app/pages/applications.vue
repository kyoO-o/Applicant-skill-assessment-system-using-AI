<script setup lang="ts">
import type { Application } from "../composables/types";
import { toast } from "vue-sonner";
import {
  FileText,
  Loader2,
  CalendarCheck,
  ChevronLeft,
  ChevronRight,
} from "lucide-vue-next";

definePageMeta({ middleware: "auth" });

const applicationsAPI = useApplicationsAPI();
const { user } = useAuth();
const router = useRouter();

const applications = ref<Application[]>([]);
const loading = ref(true);

const mapStates = ref<
  Map<number, { show: boolean; lat: number | null; lng: number | null }>
>(new Map());

function mapState(id: number) {
  if (!mapStates.value.has(id))
    mapStates.value.set(id, { show: false, lat: null, lng: null });
  return mapStates.value.get(id)!;
}

async function toggleMap(id: number, location: string | null | undefined) {
  const state = mapState(id);
  if (state.show) {
    state.show = false;
    return;
  }
  if (state.lat !== null) {
    state.show = true;
    return;
  }
  if (!location) return;
  try {
    const results = await $fetch<Array<{ lat: string; lon: string }>>(
      `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(location)}&format=json&limit=1`,
      { headers: { "User-Agent": "SkillAssessmentApp/1.0" } },
    );
    const first = results[0];
    if (first) {
      state.lat = Number(first.lat);
      state.lng = Number(first.lon);
      state.show = true;
    }
  } catch {
    /* ignore */
  }
}

const isApplicant = computed(() => user.value?.role === "user");

async function load() {
  loading.value = true;
  try {
    applications.value = await applicationsAPI.listMine();
  } catch {
    toast.error("Мэдээлэл ачааллаж чадсангүй");
  } finally {
    loading.value = false;
  }
}

await load();

const appPage = ref(1);
const APP_PAGE_SIZE = 8;
const appTotalPages = computed(() =>
  Math.max(1, Math.ceil(applications.value.length / APP_PAGE_SIZE)),
);
const paginatedApplications = computed(() => {
  const start = (appPage.value - 1) * APP_PAGE_SIZE;
  return applications.value.slice(start, start + APP_PAGE_SIZE);
});

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

function statusLabel(status: Application["status"]) {
  const map: Record<string, string> = {
    pending: "Үнэлэгдэж байна",
    assessed: "Үнэлэгдсэн",
    shortlisted: "Сонгогдсон",
    rejected: "Татгалзсан",
  };
  return map[status] ?? status;
}

function formatDate(d: string) {
  return new Intl.DateTimeFormat("mn-MN", { dateStyle: "medium" }).format(
    new Date(d),
  );
}

function formatDateTime(d: string) {
  return new Intl.DateTimeFormat("mn-MN", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(d));
}

function gcalLink(app: Application) {
  if (!app.interview_at) return "#";
  const start = new Date(app.interview_at);
  const end = new Date(start.getTime() + 60 * 60 * 1000);
  const fmt = (d: Date) =>
    d
      .toISOString()
      .replace(/[-:]/g, "")
      .replace(/\.\d{3}/, "");
  const title = encodeURIComponent(`Ярилцлага: ${app.job_title}`);
  const details = encodeURIComponent(app.interview_note || "");
  const location = encodeURIComponent(app.interview_location || "");
  return `https://calendar.google.com/calendar/render?action=TEMPLATE&text=${title}&dates=${fmt(start)}/${fmt(end)}&details=${details}&location=${location}`;
}
</script>

<template>
  <div class="space-y-5">
    <!-- Page header -->
    <div>
      <h1 class="text-[26px] font-semibold tracking-[-0.6px]">
        Илгээсэн анкетууд
      </h1>
      <p class="mt-1 text-[13.5px] text-muted-foreground">
        AI үнэлгээний үр дүнг энд харна уу.
      </p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="h-7 w-7 animate-spin text-primary" />
    </div>

    <!-- Empty state -->
    <div
      v-else-if="!applications.length"
      class="flex flex-col items-center py-20 text-center"
    >
      <div
        class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-muted"
      >
        <FileText class="h-7 w-7 text-muted-foreground" />
      </div>
      <p class="text-[16px] font-semibold">Илгээсэн анкет байхгүй</p>
      <p class="mt-2 text-[13.5px] text-muted-foreground">
        Ажлын байруудыг харж анкетаа илгээгээрэй.
      </p>
      <button
        class="mt-6 inline-flex items-center gap-2 rounded-full px-5 py-2.5 text-[13.5px] font-semibold text-white transition hover:opacity-90"
        style="
          background: linear-gradient(
            135deg,
            var(--primary),
            oklch(0.348 0.106 295)
          );
        "
        @click="router.push('/jobs')"
      >
        Ажлын байр харах
      </button>
    </div>

    <!-- Application cards -->
    <div v-else class="space-y-3.5">
      <div
        v-for="app in paginatedApplications"
        :key="app.id"
        class="group cursor-pointer rounded-2xl border border-border bg-card p-5 transition-all hover:border-primary/30 hover:shadow-sm"
        @click="router.push(`/jobs/${app.job_posting_id}`)"
      >
        <!-- Top row -->
        <div class="flex items-start justify-between gap-4">
          <div class="flex-1 min-w-0">
            <h3 class="text-[15px] font-semibold leading-snug">
              {{ app.job_title || "Ажлын байр" }}
            </h3>
            <p class="mt-0.5 text-[12px] text-muted-foreground">
              {{ formatDate(app.created_at) }}
            </p>
          </div>
          <div class="flex shrink-0 items-center gap-2.5">
            <!-- Score circle -->
            <div
              v-if="app.status !== 'pending'"
              class="flex h-12 w-12 items-center justify-center rounded-full border-[3px] shrink-0"
              :class="scoreBorderBg(app.overall_score)"
            >
              <span
                class="text-[13px] font-bold leading-none"
                :class="scoreColor(app.overall_score)"
              >
                {{ app.overall_score }}<span class="text-[9px]">%</span>
              </span>
            </div>
            <!-- Status badge -->
            <span
              class="inline-flex items-center rounded-full px-3 py-1 text-[12px] font-semibold"
              :class="{
                'badge-success': app.status === 'shortlisted',
                'badge-info': app.status === 'assessed',
                'bg-destructive/10 text-destructive': app.status === 'rejected',
                'bg-muted text-muted-foreground': app.status === 'pending',
              }"
            >
              {{ statusLabel(app.status) }}
            </span>
          </div>
        </div>

        <!-- AI pending indicator -->
        <div
          v-if="app.status === 'pending'"
          class="mt-3 flex items-center gap-2 text-[13px] text-muted-foreground"
        >
          <Loader2 class="h-3.5 w-3.5 animate-spin text-primary" />
          <span>AI үнэлгээ хийгдэж байна...</span>
        </div>

        <!-- Summary -->
        <p
          v-if="app.summary && app.status !== 'pending'"
          class="mt-3 text-[13.5px] leading-[1.55] text-muted-foreground line-clamp-2"
        >
          {{ app.summary }}
        </p>

        <!-- Skill chips -->
        <!-- <div
          v-if="
            (app.matched_skills?.length || app.missing_skills?.length) &&
            app.status !== 'pending'
          "
          class="mt-3 flex flex-wrap gap-1.5"
        >
          <span
            v-for="s in app.matched_skills?.slice(0, 4)"
            :key="s.skill"
            class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-[11.5px] font-medium badge-success"
          >
            ✓ {{ s.skill }}
          </span>
          <span
            v-for="s in app.missing_skills?.slice(0, 3)"
            :key="'m-' + s.skill"
            class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-[11.5px] font-medium bg-destructive/10 text-destructive"
          >
            ✗ {{ s.skill }}
          </span>
        </div> -->

        <!-- Interview section -->
        <div
          v-if="app.interview_at"
          class="mt-4 rounded-xl border p-4 ai-surface"
          @click.stop
        >
          <div class="flex items-center gap-1.5 mb-2">
            <CalendarCheck class="h-3.5 w-3.5 text-primary" />
            <p
              class="text-[11.5px] font-semibold tracking-[0.5px] text-primary"
            >
              Ярилцлагын хуваарь
            </p>
          </div>
          <p class="text-[13px] font-semibold">
            {{ formatDateTime(app.interview_at) }}
          </p>
          <div
            v-if="app.interview_location"
            class="mt-2 flex items-start justify-between gap-2"
          >
            <p class="text-[12.5px] text-muted-foreground">
              {{ app.interview_location }}
            </p>
            <label
              class="flex cursor-pointer items-center gap-1 text-[12px] text-muted-foreground whitespace-nowrap"
            >
              <Checkbox
                :checked="mapState(app.id).show"
                @update:checked="toggleMap(app.id, app.interview_location)"
              />
              Зураг
            </label>
          </div>
          <div
            v-if="mapState(app.id).show && mapState(app.id).lat !== null"
            class="mt-2"
          >
            <LocationMap
              :lat="mapState(app.id).lat!"
              :lng="mapState(app.id).lng!"
              :label="app.interview_location ?? undefined"
            />
          </div>
          <p
            v-if="app.interview_note"
            class="mt-1 text-[12px] text-muted-foreground italic"
          >
            {{ app.interview_note }}
          </p>
          <!-- <a
            :href="gcalLink(app)"
            target="_blank"
            rel="noopener"
            class="mt-3 inline-flex items-center gap-1.5 text-[12.5px] font-semibold text-primary hover:underline"
          >
            <CalendarCheck class="h-3.5 w-3.5" />
            Google Calendar-д нэмэх
          </a> -->
        </div>
      </div>

      <!-- Applications pagination -->
      <div class="flex items-center justify-center gap-1 pt-2">
        <button
          class="flex h-8 w-8 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
          :disabled="appPage === 1"
          @click="appPage--"
        >
          <ChevronLeft class="h-4 w-4" />
        </button>
        <button
          v-for="p in appTotalPages"
          :key="p"
          :class="[
            'h-8 w-8 rounded-full text-[13px] font-medium transition',
            p === appPage
              ? 'bg-primary text-background'
              : 'text-muted-foreground hover:bg-muted',
          ]"
          @click="appPage = p"
        >
          {{ p }}
        </button>
        <button
          class="flex h-8 w-8 items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition hover:bg-muted disabled:opacity-40"
          :disabled="appPage === appTotalPages"
          @click="appPage++"
        >
          <ChevronRight class="h-4 w-4" />
        </button>
      </div>
    </div>
  </div>
</template>
