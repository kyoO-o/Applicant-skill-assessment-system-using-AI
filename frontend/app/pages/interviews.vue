<script setup lang="ts">
import type { Application } from "../composables/types";
import type { DateValue } from "@internationalized/date";
import { CalendarRoot } from "reka-ui";
import { toast } from "vue-sonner";
import {
  CalendarCheck,
  Loader2,
  CalendarDays,
  MapPin,
  FileText,
  Video,
  X,
} from "lucide-vue-next";

definePageMeta({ middleware: "auth", ssr: false });

const { user } = useAuth();
const applicationsAPI = useApplicationsAPI();
const integrationsAPI = useIntegrationsAPI();

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

// Calendar event state
const selectedCalendarDate = ref<DateValue | undefined>(undefined);

const interviewDateSet = computed(() => {
  const set = new Set<string>();
  for (const a of interviews.value) {
    if (a.interview_at) {
      const d = new Date(a.interview_at);
      set.add(`${d.getFullYear()}-${d.getMonth() + 1}-${d.getDate()}`);
    }
  }
  return set;
});

function hasInterview(date: DateValue): boolean {
  return interviewDateSet.value.has(`${date.year}-${date.month}-${date.day}`);
}

const filteredInterviews = computed(() => {
  if (!selectedCalendarDate.value) return interviews.value;
  const sel = selectedCalendarDate.value;
  return interviews.value.filter((a) => {
    if (!a.interview_at) return false;
    const d = new Date(a.interview_at);
    return (
      d.getFullYear() === sel.year &&
      d.getMonth() + 1 === sel.month &&
      d.getDate() === sel.day
    );
  });
});

const grouped = computed(() => {
  const map = new Map<string, Application[]>();
  for (const a of filteredInterviews.value) {
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

const selectedDateLabel = computed(() => {
  if (!selectedCalendarDate.value) return null;
  const { year, month, day } = selectedCalendarDate.value;
  return new Date(year, month - 1, day).toLocaleDateString("mn-MN", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
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

    <!-- Main content: Calendar + List -->
    <div v-else class="flex flex-col lg:flex-row gap-6 items-start">
      <!-- Calendar panel -->
      <div
        class="w-full lg:w-[400px] lg:shrink-0 lg:sticky lg:top-0 lg:h-[calc(100vh-12rem)]"
      >
        <div
          class="rounded-2xl border border-border bg-card p-5 h-full flex flex-col"
        >
          <p
            class="mb-3 text-[11px] font-semibold uppercase tracking-widest text-muted-foreground shrink-0"
          >
            Ярилцлагын хуанли
          </p>

          <CalendarRoot
            v-slot="{ grid, weekDays }"
            v-model="selectedCalendarDate"
            locale="mn-MN"
            :week-starts-on="1"
            class="w-full flex-1 flex flex-col"
          >
            <CalendarHeader class="pt-0 mb-1 shrink-0">
              <nav
                class="flex items-center justify-between absolute top-0 inset-x-0"
              >
                <CalendarPrevButton />
                <CalendarNextButton />
              </nav>
              <CalendarHeading />
            </CalendarHeader>

            <div class="mt-3 flex-1 flex flex-col">
              <CalendarGrid
                v-for="month in grid"
                :key="month.value.toString()"
                class="w-full flex-1 flex flex-col"
              >
                <CalendarGridHead class="shrink-0">
                  <CalendarGridRow>
                    <CalendarHeadCell v-for="day in weekDays" :key="day">
                      {{ day }}
                    </CalendarHeadCell>
                  </CalendarGridRow>
                </CalendarGridHead>
                <CalendarGridBody class="flex-1 flex flex-col">
                  <CalendarGridRow
                    v-for="(weekDates, idx) in month.rows"
                    :key="`week-${idx}`"
                    class="flex-1 w-full"
                  >
                    <CalendarCell
                      v-for="weekDate in weekDates"
                      :key="weekDate.toString()"
                      :date="weekDate"
                      class="h-full"
                    >
                      <CalendarCellTrigger
                        :day="weekDate"
                        :month="month.value"
                        class="size-full rounded-lg"
                      >
                        <span
                          class="relative flex flex-col items-center pb-1.5"
                        >
                          {{ weekDate.day }}
                          <span
                            v-if="hasInterview(weekDate)"
                            class="absolute bottom-0 h-1.5 w-1.5 rounded-full bg-primary"
                          />
                        </span>
                      </CalendarCellTrigger>
                    </CalendarCell>
                  </CalendarGridRow>
                </CalendarGridBody>
              </CalendarGrid>
            </div>
          </CalendarRoot>

          <!-- Legend -->
          <div
            class="mt-4 shrink-0 flex items-center gap-2 text-[12px] text-muted-foreground"
          >
            <span class="inline-block h-1.5 w-1.5 rounded-full bg-primary" />
            Ярилцлага товлогдсон
          </div>

          <!-- Clear filter -->
          <button
            v-if="selectedCalendarDate"
            class="mt-3 shrink-0 flex w-full items-center justify-center gap-1.5 rounded-lg border border-border py-1.5 text-[12px] font-medium text-muted-foreground transition hover:bg-muted"
            @click="selectedCalendarDate = undefined"
          >
            <X class="h-3 w-3" />
            Бүгдийг харах
          </button>
        </div>
      </div>

      <!-- Interview list -->
      <div class="flex-1 min-w-0 space-y-7">
        <!-- Selected date banner -->
        <div
          v-if="selectedCalendarDate"
          class="flex items-center gap-2 rounded-xl border border-primary/20 bg-primary/5 px-4 py-2.5"
        >
          <CalendarCheck class="h-4 w-4 shrink-0 text-primary" />
          <span class="text-[13px] font-semibold text-primary">
            {{ selectedDateLabel }}
          </span>
          <span class="text-[13px] text-primary/70">
            — {{ filteredInterviews.length }} ярилцлага
          </span>
        </div>

        <!-- No interviews for selected date -->
        <div
          v-if="selectedCalendarDate && filteredInterviews.length === 0"
          class="flex flex-col items-center py-16 text-center"
        >
          <div
            class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-muted"
          >
            <CalendarDays class="h-5 w-5 text-muted-foreground" />
          </div>
          <p class="text-[14px] font-semibold">Энэ өдөр ярилцлага байхгүй</p>
          <p class="mt-1 text-[13px] text-muted-foreground">
            Ярилцлага байгаа өдрийг хуанлиас сонгоно уу.
          </p>
        </div>

        <!-- Grouped by date -->
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
          <div class="grid gap-3.5 sm:grid-cols-2 lg:grid-cols-2">
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
                      v-if="
                        mapState(app.id).show && mapState(app.id).lat !== null
                      "
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
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
