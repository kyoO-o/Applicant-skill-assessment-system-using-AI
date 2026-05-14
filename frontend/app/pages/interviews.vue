<script setup lang="ts">
import type { Application } from "../composables/types";
import { toast } from "vue-sonner";
import { CalendarCheck, Loader2, CalendarDays, MapPin, FileText } from "lucide-vue-next";

definePageMeta({ middleware: "auth" });

const { user } = useAuth();
const applicationsAPI = useApplicationsAPI();
const integrationsAPI = useIntegrationsAPI();
const router = useRouter();

const interviews = ref<Application[]>([]);
const loading = ref(true);
const gcalConnected = ref(false);

// Map toggle per card: appID → {show, lat, lng} | null
const mapStates = ref<Map<number, { show: boolean; lat: number | null; lng: number | null }>>(new Map());

function mapState(id: number) {
  if (!mapStates.value.has(id)) mapStates.value.set(id, { show: false, lat: null, lng: null });
  return mapStates.value.get(id)!;
}

async function toggleMap(id: number, location: string | null | undefined) {
  const state = mapState(id);
  if (state.show) { state.show = false; return; }
  if (state.lat !== null) { state.show = true; return; }
  if (!location) return;
  try {
    const results = await $fetch<Array<{ lat: string; lon: string }>>(
      `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(location)}&format=json&limit=1`,
      { headers: { "User-Agent": "SkillAssessmentApp/1.0" } },
    );
    if (results.length) {
      state.lat = Number(results[0].lat);
      state.lng = Number(results[0].lon);
      state.show = true;
    }
  } catch { /* ignore */ }
}

const isRecruiter = computed(() => user.value?.role === "recruiter");

async function load() {
  loading.value = true;
  try {
    if (isRecruiter.value) {
      const [data, status] = await Promise.allSettled([
        applicationsAPI.listInterviews(),
        integrationsAPI.googleCalendarStatus(),
      ]);
      if (data.status === "fulfilled") interviews.value = data.value;
      if (status.status === "fulfilled") gcalConnected.value = status.value.connected;
    } else {
      const apps = await applicationsAPI.listMine();
      interviews.value = apps.filter((a) => a.interview_at != null);
    }
  } catch {
    toast.error("Мэдээлэл ачааллаж чадсангүй");
  } finally {
    loading.value = false;
  }
}

await load();

// Group by date (YYYY-MM-DD)
const grouped = computed(() => {
  const map = new Map<string, Application[]>();
  for (const a of interviews.value) {
    const key = a.interview_at
      ? new Date(a.interview_at).toLocaleDateString("mn-MN", { year: "numeric", month: "long", day: "numeric" })
      : "Огноогүй";
    if (!map.has(key)) map.set(key, []);
    map.get(key)!.push(a);
  }
  return map;
});

function formatTime(d: string | null | undefined) {
  if (!d) return "";
  return new Intl.DateTimeFormat("mn-MN", { timeStyle: "short" }).format(new Date(d));
}

function formatDateTime(d: string | null | undefined) {
  if (!d) return "—";
  return new Intl.DateTimeFormat("mn-MN", { dateStyle: "medium", timeStyle: "short" }).format(new Date(d));
}

function gcalLink(app: Application) {
  if (!app.interview_at) return "#";
  const start = new Date(app.interview_at);
  const end = new Date(start.getTime() + 60 * 60 * 1000);
  const fmt = (d: Date) =>
    d.toISOString().replace(/[-:]/g, "").replace(/\.\d{3}/, "");
  const title = encodeURIComponent(`Ярилцлага: ${app.job_title}`);
  const details = encodeURIComponent(
    [
      isRecruiter.value ? `Горилогч: ${app.applicant_name}` : "",
      app.interview_note || "",
    ]
      .filter(Boolean)
      .join("\n"),
  );
  const location = encodeURIComponent(app.interview_location || "");
  return `https://calendar.google.com/calendar/render?action=TEMPLATE&text=${title}&dates=${fmt(start)}/${fmt(end)}&details=${details}&location=${location}`;
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <div class="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
        <div class="space-y-2">
          <p class="text-sm font-medium text-muted-foreground">Ярилцлага</p>
          <h1 class="text-3xl font-semibold tracking-tight">
            {{ isRecruiter ? "Ярилцлагын урсгал" : "Миний ярилцлагууд" }}
          </h1>
          <p class="max-w-2xl text-sm leading-6 text-muted-foreground">
            {{ isRecruiter
              ? "Компанийн бүх товлогдсон ярилцлагуудыг огноогоор эрэмбэлэн харна уу."
              : "Ажил олгогч тань товлосон ярилцлагуудыг энд харна уу." }}
          </p>
        </div>
      </div>
    </section>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <!-- Empty -->
    <div
      v-else-if="interviews.length === 0"
      class="rounded-3xl border border-dashed border-border px-6 py-20 text-center"
    >
      <CalendarDays class="mx-auto h-10 w-10 text-muted-foreground" />
      <p class="mt-4 text-lg font-medium">Товлогдсон ярилцлага байхгүй</p>
      <p class="mt-2 text-sm text-muted-foreground">
        {{ isRecruiter
          ? "Горилогчид ярилцлага товлох үед энд харагдана."
          : "Ажил олгогч таны ярилцлага товлоход энд харагдана." }}
      </p>
    </div>

    <!-- Grouped by date -->
    <div v-else class="space-y-8">
      <div v-for="[dateLabel, items] in grouped" :key="dateLabel" class="space-y-3">
        <!-- Date heading -->
        <div class="flex items-center gap-3">
          <CalendarCheck class="h-4 w-4 text-purple-600" />
          <h2 class="text-sm font-semibold text-purple-700">{{ dateLabel }}</h2>
          <div class="flex-1 border-t border-purple-100" />
        </div>

        <!-- Interview cards -->
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div
            v-for="app in items"
            :key="app.id"
            class="rounded-3xl border border-border bg-card px-5 py-4 shadow-sm transition hover:shadow-md"
          >
            <!-- Time badge -->
            <div class="mb-3 flex items-center justify-between">
              <span class="inline-flex items-center gap-1.5 rounded-full bg-purple-50 border border-purple-200 px-3 py-1 text-xs font-semibold text-purple-700">
                <CalendarCheck class="h-3 w-3" />
                {{ formatTime(app.interview_at) }}
              </span>
              <span
                v-if="isRecruiter && gcalConnected"
                class="inline-flex items-center gap-1 rounded-full bg-green-50 border border-green-200 px-2 py-0.5 text-xs text-green-700"
                title="Google Calendar-д автоматаар нэмэгдсэн"
              >
                ✓ GCal
              </span>
              <a
                v-else
                :href="gcalLink(app)"
                target="_blank"
                rel="noopener"
                class="text-xs text-muted-foreground underline underline-offset-2 hover:text-foreground"
                title="Google Calendar-д нэмэх"
              >
                GCal
              </a>
            </div>

            <!-- Job title -->
            <p class="font-semibold text-sm leading-snug">{{ app.job_title || "Ажлын байр" }}</p>

            <!-- Applicant (recruiter view) -->
            <template v-if="isRecruiter">
              <p class="mt-1 text-sm text-muted-foreground font-medium">{{ app.applicant_name || "Горилогч" }}</p>
              <p class="text-xs text-muted-foreground">{{ app.applicant_email }}</p>
            </template>

            <!-- Location -->
            <div v-if="app.interview_location" class="mt-2 space-y-2">
              <div class="flex items-start justify-between gap-2">
                <div class="flex items-start gap-1.5 text-xs text-muted-foreground">
                  <MapPin class="mt-0.5 h-3 w-3 shrink-0" />
                  <span>{{ app.interview_location }}</span>
                </div>
                <label class="flex cursor-pointer items-center gap-1 text-xs text-muted-foreground whitespace-nowrap">
                  <Checkbox
                    :checked="mapState(app.id).show"
                    @update:checked="toggleMap(app.id, app.interview_location)"
                  />
                  Зураг
                </label>
              </div>
              <LocationMap
                v-if="mapState(app.id).show && mapState(app.id).lat !== null"
                :lat="mapState(app.id).lat!"
                :lng="mapState(app.id).lng!"
                :label="app.interview_location"
              />
            </div>

            <!-- Note -->
            <div v-if="app.interview_note" class="mt-2 flex items-start gap-1.5 text-xs text-muted-foreground italic">
              <FileText class="mt-0.5 h-3 w-3 shrink-0" />
              <span class="line-clamp-2">{{ app.interview_note }}</span>
            </div>

            <!-- Full date/time -->
            <p class="mt-3 text-xs text-muted-foreground">{{ formatDateTime(app.interview_at) }}</p>

            <!-- Google Calendar link (full) — hidden if recruiter already auto-synced -->
            <template v-if="!(isRecruiter && gcalConnected)">
              <a
                :href="gcalLink(app)"
                target="_blank"
                rel="noopener"
                class="mt-3 inline-flex w-full items-center justify-center gap-1.5 rounded-2xl border border-purple-200 bg-purple-50 px-3 py-2 text-xs font-medium text-purple-700 transition hover:bg-purple-100"
              >
                <CalendarCheck class="h-3.5 w-3.5" />
                Google Calendar-д нэмэх
              </a>
            </template>
            <template v-else>
              <p class="mt-3 flex items-center justify-center gap-1.5 text-xs text-green-600">
                <CalendarCheck class="h-3.5 w-3.5" />
                Google Calendar-д автоматаар нэмэгдсэн
              </p>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
