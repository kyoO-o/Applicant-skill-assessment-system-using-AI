<script setup lang="ts">
import type { Application } from "../../../composables/types";
import { toast } from "vue-sonner";
import { Users, Loader2, ChevronLeft, CheckCircle2, XCircle } from "lucide-vue-next";

definePageMeta({ middleware: "auth" });

const route = useRoute();
const router = useRouter();
const applicationsAPI = useApplicationsAPI();
const { user } = useAuth();

const jobID = computed(() => Number(route.params.id));
const applications = ref<Application[]>([]);
const selected = ref<Application | null>(null);
const loading = ref(true);
const updatingStatus = ref(false);

const isRecruiter = computed(() => user.value?.role === "recruiter");

async function load() {
  loading.value = true;
  try {
    applications.value = await applicationsAPI.listForJob(jobID.value);
  } catch {
    toast.error("Мэдээлэл ачааллаж чадсангүй");
  } finally {
    loading.value = false;
  }
}

async function updateStatus(id: number, status: string) {
  updatingStatus.value = true;
  try {
    await applicationsAPI.updateStatus(id, status);
    await load();
    if (selected.value?.id === id) {
      selected.value = applications.value.find((a) => a.id === id) ?? null;
    }
    toast.success("Төлөв шинэчлэгдлээ");
  } catch {
    toast.error("Алдаа гарлаа");
  } finally {
    updatingStatus.value = false;
  }
}

if (!isRecruiter.value) {
  router.replace("/jobs");
}

await load();

function scoreColor(score: number) {
  if (score >= 75) return "text-green-600";
  if (score >= 50) return "text-yellow-600";
  return "text-red-600";
}

function scoreBg(score: number) {
  if (score >= 75) return "bg-green-50";
  if (score >= 50) return "bg-yellow-50";
  return "bg-red-50";
}

function statusLabel(status: Application["status"]) {
  const map: Record<string, string> = {
    pending: "Хүлээгдэж байна",
    assessed: "Үнэлэгдсэн",
    shortlisted: "Сонгогдсон",
    rejected: "Татгалзсан",
  };
  return map[status] ?? status;
}
</script>

<template>
  <div class="space-y-6">
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <div class="flex items-center gap-3">
        <Button variant="ghost" size="icon" class="rounded-xl" @click="router.back()">
          <ChevronLeft class="h-4 w-4" />
        </Button>
        <div>
          <p class="text-sm font-medium text-muted-foreground">Ирсэн анкетууд</p>
          <h2 class="text-2xl font-semibold tracking-tight">
            {{ applications.length }} анкет ирсэн
          </h2>
        </div>
      </div>
    </section>

    <div v-if="loading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <div v-else-if="!applications.length"
         class="rounded-3xl border border-dashed border-border px-6 py-20 text-center">
      <Users class="mx-auto h-10 w-10 text-muted-foreground" />
      <p class="mt-4 text-lg font-medium">Анкет ирээгүй байна</p>
    </div>

    <div v-else class="grid gap-6 lg:grid-cols-5">
      <!-- List -->
      <div class="space-y-3 lg:col-span-2">
        <div v-for="app in applications" :key="app.id"
             :class="['cursor-pointer rounded-3xl border px-4 py-4 shadow-sm transition',
                      selected?.id === app.id ? 'border-primary bg-primary/5' : 'border-border bg-card hover:shadow-md']"
             @click="selected = app">
          <div class="flex items-start justify-between gap-2">
            <div class="min-w-0">
              <p class="truncate font-medium">{{ app.applicant_name || "Хэрэглэгч" }}</p>
              <p class="truncate text-xs text-muted-foreground">{{ app.applicant_email }}</p>
            </div>
            <div :class="['shrink-0 text-right', scoreBg(app.overall_score), 'rounded-xl px-2 py-1']">
              <p :class="['text-lg font-bold', scoreColor(app.overall_score)]">{{ app.overall_score }}%</p>
            </div>
          </div>
          <div class="mt-2 flex items-center justify-between">
            <Badge variant="outline" class="rounded-full px-2 text-xs">{{ statusLabel(app.status) }}</Badge>
          </div>
        </div>
      </div>

      <!-- Detail -->
      <div class="lg:col-span-3">
        <div v-if="!selected" class="rounded-3xl border border-dashed border-border px-6 py-20 text-center text-sm text-muted-foreground">
          Анкет сонгоно уу
        </div>

        <div v-else class="space-y-4">
          <Card class="rounded-3xl border-border shadow-sm">
            <CardHeader>
              <div class="flex items-start justify-between">
                <div>
                  <CardTitle>{{ selected.applicant_name || "Хэрэглэгч" }}</CardTitle>
                  <CardDescription>{{ selected.applicant_email }}</CardDescription>
                </div>
                <p :class="['text-4xl font-bold', scoreColor(selected.overall_score)]">
                  {{ selected.overall_score }}%
                </p>
              </div>
            </CardHeader>
            <CardContent class="space-y-4">
              <p class="text-sm leading-6 text-muted-foreground">{{ selected.summary }}</p>

              <!-- Status actions -->
              <div v-if="selected.status === 'assessed'" class="flex gap-2">
                <Button class="rounded-full" :disabled="updatingStatus" @click="updateStatus(selected.id, 'shortlisted')">
                  Сонгох
                </Button>
                <Button variant="outline" class="rounded-full" :disabled="updatingStatus" @click="updateStatus(selected.id, 'rejected')">
                  Татгалзах
                </Button>
              </div>
              <div v-else>
                <Badge :variant="selected.status === 'shortlisted' ? 'default' : selected.status === 'rejected' ? 'destructive' : 'secondary'" class="rounded-full px-3">
                  {{ statusLabel(selected.status) }}
                </Badge>
              </div>
            </CardContent>
          </Card>

          <div class="grid gap-4 sm:grid-cols-2">
            <Card class="rounded-3xl border-border shadow-sm">
              <CardHeader>
                <CardTitle class="flex items-center gap-2 text-sm text-green-700">
                  <CheckCircle2 class="h-4 w-4" /> Тохирсон ур чадварууд
                </CardTitle>
              </CardHeader>
              <CardContent class="space-y-2">
                <div v-for="s in selected.matched_skills" :key="s.skill" class="rounded-xl bg-green-50 px-3 py-2">
                  <p class="text-xs font-medium text-green-800">✓ {{ s.skill }}</p>
                  <p class="mt-0.5 text-xs text-green-600">{{ s.explanation }}</p>
                </div>
                <p v-if="!selected.matched_skills?.length" class="text-xs text-muted-foreground">Байхгүй</p>
              </CardContent>
            </Card>

            <Card class="rounded-3xl border-border shadow-sm">
              <CardHeader>
                <CardTitle class="flex items-center gap-2 text-sm text-red-700">
                  <XCircle class="h-4 w-4" /> Дутуу ур чадварууд
                </CardTitle>
              </CardHeader>
              <CardContent class="space-y-2">
                <div v-for="s in selected.missing_skills" :key="s.skill" class="rounded-xl bg-red-50 px-3 py-2">
                  <p class="text-xs font-medium text-red-800">✗ {{ s.skill }}</p>
                  <p class="mt-0.5 text-xs text-red-600">{{ s.explanation }}</p>
                </div>
                <p v-if="!selected.missing_skills?.length" class="text-xs text-muted-foreground">Байхгүй</p>
              </CardContent>
            </Card>
          </div>

          <Card v-if="selected.recommendations?.length" class="rounded-3xl border-border shadow-sm">
            <CardHeader><CardTitle class="text-sm">Зөвлөмж</CardTitle></CardHeader>
            <CardContent>
              <ul class="space-y-1">
                <li v-for="r in selected.recommendations" :key="r" class="text-sm text-muted-foreground">• {{ r }}</li>
              </ul>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  </div>
</template>
