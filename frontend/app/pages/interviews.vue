<script setup lang="ts">
import type { Application } from "../composables/types";
import { toast } from "vue-sonner";
import {
  CalendarCheck,
  Loader2,
  CalendarDays,
  MapPin,
  FileText,
  Video,
} from "lucide-vue-next";

definePageMeta({ middleware: "auth" });

const { user } = useAuth();
const applicationsAPI = useApplicationsAPI();
const integrationsAPI = useIntegrationsAPI();
const router = useRouter();

const interviews = ref<Application[]>([]);
const loading = ref(true);
const gcalConnected = ref(false);

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
      if (status.status === "fulfilled")
        gcalConnected.value = status.value.connected;
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

const grouped = computed(() => {
  const map = new Map<string, Application[]>();
  for (const a of interviews.value) {
    const key = a.interview_at
      ? new Date(a.interview_at).toLocaleDateString("mn-MN", {
          year: "numeric",
          month: "long",
          day: "numeric",
        })
      : "Огноогүй";
    if (!map.has(key)) map.set(key, []);
    map.get(key)!.push(a);
  }
  return map;
});

function formatTime(d: string | null | undefined) {
  if (!d) return "";
  return new Intl.DateTimeFormat("mn-MN", { timeStyle: "short" }).format(
    new Date(d),
  );
}

function isMeetUrl(val: string | null | undefined): boolean {
  return !!val && val.startsWith("https://meet.google.com");
}

function formatDateTime(d: string | null | undefined) {
  if (!d) return "—";
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
  <div class="space-y-5">
    <!-- Page header -->
    <div class="flex items-end justify-between gap-4">
      <div>
        <h1 class="text-[26px] font-semibold tracking-[-0.6px]">
          {{ isRecruiter ? "Ярилцлагын урсгал" : "Миний ярилцлагууд" }}
        </h1>
        <p class="mt-1 text-[13.5px] text-muted-foreground">
          {{
            isRecruiter
              ? "Компанийн бүх товлогдсон ярилцлагуудыг огноогоор эрэмбэлэн харна уу."
              : "Ажил олгогч тань товлосон ярилцлагуудыг энд харна уу."
          }}
        </p>
      </div>
      <div
        v-if="isRecruiter && gcalConnected"
        class="inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-[12px] font-semibold badge-success"
      >
        <CalendarCheck class="h-3.5 w-3.5" />
        Google Calendar холбогдсон
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="h-7 w-7 animate-spin text-primary" />
    </div>

    <!-- Empty -->
    <div
      v-else-if="interviews.length === 0"
      class="flex flex-col items-center py-20 text-center"
    >
      <div
        class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-muted"
      >
        <CalendarDays class="h-7 w-7 text-muted-foreground" />
      </div>
      <p class="text-[16px] font-semibold">Товлогдсон ярилцлага байхгүй</p>
      <p class="mt-2 text-[13.5px] text-muted-foreground">
        {{
          isRecruiter
            ? "Горилогчид ярилцлага товлох үед энд харагдана."
            : "Ажил олгогч таны ярилцлага товлоход энд харагдана."
        }}
      </p>
    </div>

    <!-- Grouped by date -->
    <div v-else class="space-y-7">
      <div
        v-for="[dateLabel, items] in grouped"
        :key="dateLabel"
        class="space-y-3"
      >
        <!-- Date group header -->
        <div class="flex items-center gap-3">
          <div
            class="flex h-7 w-7 items-center justify-center rounded-lg bg-primary/10"
          >
            <CalendarCheck class="h-3.5 w-3.5 text-primary" />
          </div>
          <h2 class="text-[13.5px] font-semibold text-primary">
            {{ dateLabel }}
          </h2>
          <div class="flex-1 border-t border-primary/15" />
          <span class="text-[12px] font-medium text-muted-foreground"
            >{{ items.length }} ярилцлага</span
          >
        </div>

        <!-- Cards -->
        <div class="grid gap-3.5 sm:grid-cols-2 lg:grid-cols-3">
          <div
            v-for="app in items"
            :key="app.id"
            class="rounded-2xl border border-border bg-card p-5 transition hover:border-primary/30 hover:shadow-sm"
          >
            <!-- Time badge -->
            <div class="mb-3 flex items-center justify-between gap-2">
              <span
                class="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-[12px] font-semibold bg-primary/10 text-primary"
              >
                <CalendarCheck class="h-3 w-3" />
                {{ formatTime(app.interview_at) }}
              </span>
              <!-- <span
                v-if="isRecruiter && gcalConnected"
                class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-[11px] font-semibold badge-success"
                title="Google Calendar-д автоматаар нэмэгдсэн"
              >
                ✓ GCal
              </span> -->
            </div>

            <!-- Job title -->
            <p class="text-[14.5px] font-semibold leading-snug">
              {{ app.job_title || "Ажлын байр" }}
            </p>

            <!-- Applicant (recruiter view) -->
            <template v-if="isRecruiter">
              <p class="mt-1.5 text-[13px] font-medium text-foreground">
                {{ app.applicant_name || "Горилогч" }}
              </p>
              <p class="text-[12px] text-muted-foreground">
                {{ app.applicant_email }}
              </p>
            </template>

            <!-- Location -->
            <div v-if="app.interview_location" class="mt-3">
              <!-- Google Meet link -->
              <a
                v-if="isMeetUrl(app.interview_location)"
                :href="app.interview_location"
                target="_blank"
                rel="noopener"
                class="inline-flex items-center gap-1.5 rounded-xl border border-border px-3 py-1.5 text-[12.5px] font-semibold text-primary transition hover:bg-primary/5"
              >
                <Video class="h-3.5 w-3.5 shrink-0" />
                Google Meet нэгдэх
              </a>

              <!-- Onsite location -->
              <template v-else>
                <div class="space-y-2">
                  <div class="flex items-start justify-between gap-2">
                    <div
                      class="flex items-start gap-1.5 text-[12.5px] text-muted-foreground"
                    >
                      <MapPin class="mt-0.5 h-3.5 w-3.5 shrink-0" />
                      <span>{{ app.interview_location }}</span>
                    </div>
                    <label
                      class="flex cursor-pointer items-center gap-1 text-[12px] text-muted-foreground whitespace-nowrap"
                    >
                      <Checkbox
                        :checked="mapState(app.id).show"
                        @update:checked="
                          toggleMap(app.id, app.interview_location)
                        "
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
              </template>
            </div>

            <!-- Note -->
            <div
              v-if="app.interview_note"
              class="mt-2.5 flex items-start gap-1.5 text-[12.5px] text-muted-foreground italic"
            >
              <FileText class="mt-0.5 h-3.5 w-3.5 shrink-0" />
              <span class="line-clamp-2">{{ app.interview_note }}</span>
            </div>

            <!-- Full datetime -->
            <p class="mt-3 text-[12px] text-muted-foreground">
              {{ formatDateTime(app.interview_at) }}
            </p>

            <!-- Calendar CTA -->
            <!-- <template v-if="!(isRecruiter && gcalConnected)">
              <a
                :href="gcalLink(app)"
                target="_blank"
                rel="noopener"
                class="mt-3 flex w-full items-center justify-center gap-1.5 rounded-xl border p-2 text-[12.5px] font-semibold transition ai-surface hover:opacity-80 text-primary"
              >
                <CalendarCheck class="h-3.5 w-3.5" />
                Google Calendar-д нэмэх
              </a>
            </template> -->
            <!-- <template v-else>
              <p
                class="mt-3 flex items-center justify-center gap-1.5 text-[12px] text-success-foreground font-medium"
              >
                <CalendarCheck class="h-3.5 w-3.5" />
                Google Calendar-д нэмэгдсэн
              </p>
            </template> -->
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
