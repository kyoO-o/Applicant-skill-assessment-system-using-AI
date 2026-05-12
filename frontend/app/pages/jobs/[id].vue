<script setup lang="ts">
import type { Job, Application } from "../../composables/types";
import { toast } from "vue-sonner";
import { BriefcaseBusiness, MapPin, Upload, CheckCircle2, XCircle, ChevronRight, Loader2 } from "lucide-vue-next";

definePageMeta({ middleware: "auth" });

const route = useRoute();
const { user } = useAuth();
const jobsAPI = useJobsAPI();
const applicationsAPI = useApplicationsAPI();

const job = ref<Job | null>(null);
const application = ref<Application | null>(null);
const loading = ref(true);
const applyOpen = ref(false);
const cvFile = ref<File | null>(null);
const isSubmitting = ref(false);
const pollTimer = ref<ReturnType<typeof setInterval> | null>(null);

const jobID = computed(() => Number(route.params.id));
const isApplicant = computed(() => user.value?.role === "user");

async function loadJob() {
  try {
    const jobs = await jobsAPI.list();
    job.value = jobs.find((j) => j.id === jobID.value) ?? null;
  } catch {
    toast.error("Ажлын байрны мэдээлэл ачааллаж чадсангүй");
  }
}

async function loadMyApplication() {
  if (!isApplicant.value) return;
  try {
    const apps = await applicationsAPI.listMine();
    application.value = apps.find((a) => a.job_posting_id === jobID.value) ?? null;
  } catch { /* not applied yet */ }
}

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement;
  cvFile.value = input.files?.[0] ?? null;
}

async function submitApplication() {
  if (!cvFile.value) {
    toast.error("CV файлаа сонгоно уу");
    return;
  }
  isSubmitting.value = true;
  try {
    await applicationsAPI.applyToJob(jobID.value, cvFile.value);
    applyOpen.value = false;
    toast.success("Анкет амжилттай илгээгдлээ! AI үнэлгээ хийгдэж байна...");
    await loadMyApplication();
    startPolling();
  } catch (e: any) {
    toast.error(e?.data?.message ?? "Анкет илгээхэд алдаа гарлаа");
  } finally {
    isSubmitting.value = false;
  }
}

function startPolling() {
  stopPolling();
  pollTimer.value = setInterval(async () => {
    if (application.value?.status !== "pending") {
      stopPolling();
      return;
    }
    const apps = await applicationsAPI.listMine().catch(() => []);
    const updated = apps.find((a) => a.job_posting_id === jobID.value);
    if (updated) application.value = updated;
    if (updated?.status !== "pending") stopPolling();
  }, 4000);
}

function stopPolling() {
  if (pollTimer.value) {
    clearInterval(pollTimer.value);
    pollTimer.value = null;
  }
}

onUnmounted(stopPolling);

loading.value = true;
await loadJob();
await loadMyApplication();
loading.value = false;

if (application.value?.status === "pending") startPolling();

function scoreColor(score: number) {
  if (score >= 75) return "text-green-600";
  if (score >= 50) return "text-yellow-600";
  return "text-red-600";
}

function scoreBg(score: number) {
  if (score >= 75) return "bg-green-50 border-green-200";
  if (score >= 50) return "bg-yellow-50 border-yellow-200";
  return "bg-red-50 border-red-200";
}
</script>

<template>
  <div v-if="loading" class="flex justify-center py-20">
    <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
  </div>

  <div v-else-if="!job" class="rounded-3xl border border-dashed border-border px-6 py-20 text-center text-muted-foreground">
    Ажлын байр олдсонгүй
  </div>

  <div v-else class="space-y-6">
    <!-- Job header -->
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
        <div class="space-y-2">
          <h1 class="text-3xl font-semibold tracking-tight">{{ job.title }}</h1>
          <div class="flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
            <span class="flex items-center gap-1">
              <MapPin class="h-4 w-4" /> {{ job.location }}
            </span>
            <Badge variant="outline" class="rounded-full px-3">{{ job.type || job.employment_type }}</Badge>
            <Badge variant="outline" class="rounded-full px-3">{{ job.level || job.seniority }}</Badge>
          </div>
          <p v-if="job.min_salary || job.max_salary" class="text-sm font-medium">
            {{ job.min_salary?.toLocaleString() }}₮ – {{ job.max_salary?.toLocaleString() }}₮
          </p>
        </div>

        <!-- Apply button -->
        <div v-if="isApplicant">
          <div v-if="!application">
            <Button class="rounded-full px-6" @click="applyOpen = true">
              <Upload class="mr-2 h-4 w-4" /> Анкет илгээх
            </Button>
          </div>
          <div v-else-if="application.status === 'pending'" class="flex items-center gap-2 text-sm text-muted-foreground">
            <Loader2 class="h-4 w-4 animate-spin" /> AI үнэлгээ хийгдэж байна...
          </div>
          <div v-else>
            <Badge class="rounded-full px-4 py-1.5" variant="secondary">Анкет илгээсэн</Badge>
          </div>
        </div>
      </div>
    </section>

    <!-- Job details -->
    <div class="grid gap-6 lg:grid-cols-3">
      <div class="space-y-6 lg:col-span-2">
        <Card class="rounded-3xl border-border shadow-sm">
          <CardHeader><CardTitle>Тайлбар</CardTitle></CardHeader>
          <CardContent>
            <p class="text-sm leading-7 text-muted-foreground whitespace-pre-line">{{ job.additional_info || job.description }}</p>
          </CardContent>
        </Card>

        <Card v-if="job.duties?.length" class="rounded-3xl border-border shadow-sm">
          <CardHeader><CardTitle>Үүрэг хариуцлага</CardTitle></CardHeader>
          <CardContent>
            <ul class="space-y-2">
              <li v-for="duty in job.duties" :key="duty" class="flex items-start gap-2 text-sm text-muted-foreground">
                <ChevronRight class="mt-0.5 h-4 w-4 shrink-0 text-primary" /> {{ duty }}
              </li>
            </ul>
          </CardContent>
        </Card>

        <Card v-if="job.requirements?.length" class="rounded-3xl border-border shadow-sm">
          <CardHeader><CardTitle>Шаардлагууд</CardTitle></CardHeader>
          <CardContent>
            <ul class="space-y-2">
              <li v-for="req in job.requirements" :key="req" class="flex items-start gap-2 text-sm text-muted-foreground">
                <ChevronRight class="mt-0.5 h-4 w-4 shrink-0 text-primary" /> {{ req }}
              </li>
            </ul>
          </CardContent>
        </Card>
      </div>

      <div class="space-y-6">
        <Card v-if="job.skills?.length" class="rounded-3xl border-border shadow-sm">
          <CardHeader><CardTitle>Шаардлагатай ур чадвар</CardTitle></CardHeader>
          <CardContent class="flex flex-wrap gap-2">
            <Badge v-for="skill in job.skills" :key="skill" variant="secondary" class="rounded-full px-3">
              {{ skill }}
            </Badge>
          </CardContent>
        </Card>

        <Card v-if="job.bonuses?.length" class="rounded-3xl border-border shadow-sm">
          <CardHeader><CardTitle>Нэмэлт давуу тал</CardTitle></CardHeader>
          <CardContent>
            <ul class="space-y-1">
              <li v-for="bonus in job.bonuses" :key="bonus" class="text-sm text-muted-foreground">✓ {{ bonus }}</li>
            </ul>
          </CardContent>
        </Card>

        <Card v-if="job.contact_info" class="rounded-3xl border-border shadow-sm">
          <CardHeader><CardTitle>Холбоо барих</CardTitle></CardHeader>
          <CardContent>
            <p class="text-sm text-muted-foreground">{{ job.contact_info }}</p>
          </CardContent>
        </Card>
      </div>
    </div>

    <!-- AI Assessment results -->
    <div v-if="application && application.status !== 'pending'" class="space-y-4">
      <h2 class="text-xl font-semibold">Таны AI үнэлгээний үр дүн</h2>

      <div :class="['rounded-3xl border px-6 py-6 shadow-sm', scoreBg(application.overall_score)]">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-muted-foreground">Нийт тохирлын хувь</p>
            <p :class="['text-5xl font-bold', scoreColor(application.overall_score)]">
              {{ application.overall_score }}%
            </p>
          </div>
          <div class="text-right">
            <Badge :class="[
              'rounded-full px-4 py-1.5 text-sm',
              application.status === 'shortlisted' ? 'bg-green-100 text-green-800' :
              application.status === 'rejected' ? 'bg-red-100 text-red-800' :
              'bg-blue-100 text-blue-800'
            ]">
              {{ application.status === 'shortlisted' ? 'Сонгогдсон' :
                 application.status === 'rejected' ? 'Татгалзсан' : 'Үнэлэгдсэн' }}
            </Badge>
          </div>
        </div>
        <p class="mt-4 text-sm leading-6">{{ application.summary }}</p>
      </div>

      <div class="grid gap-4 lg:grid-cols-2">
        <Card class="rounded-3xl border-border shadow-sm">
          <CardHeader>
            <CardTitle class="flex items-center gap-2 text-green-700">
              <CheckCircle2 class="h-5 w-5" /> Тохирсон ур чадварууд
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div v-if="application.matched_skills?.length" class="space-y-3">
              <div v-for="skill in application.matched_skills" :key="skill.skill"
                   class="rounded-2xl bg-green-50 px-4 py-3">
                <p class="text-sm font-medium text-green-800">✓ {{ skill.skill }}</p>
                <p class="mt-1 text-xs text-green-600">{{ skill.explanation }}</p>
              </div>
            </div>
            <p v-else class="text-sm text-muted-foreground">Тохирсон ур чадвар олдсонгүй</p>
          </CardContent>
        </Card>

        <Card class="rounded-3xl border-border shadow-sm">
          <CardHeader>
            <CardTitle class="flex items-center gap-2 text-red-700">
              <XCircle class="h-5 w-5" /> Дутуу ур чадварууд
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div v-if="application.missing_skills?.length" class="space-y-3">
              <div v-for="skill in application.missing_skills" :key="skill.skill"
                   class="rounded-2xl bg-red-50 px-4 py-3">
                <p class="text-sm font-medium text-red-800">✗ {{ skill.skill }}</p>
                <p class="mt-1 text-xs text-red-600">{{ skill.explanation }}</p>
              </div>
            </div>
            <p v-else class="text-sm text-muted-foreground">Бүх ур чадвар хангасан</p>
          </CardContent>
        </Card>
      </div>

      <Card v-if="application.recommendations?.length" class="rounded-3xl border-border shadow-sm">
        <CardHeader><CardTitle>Сайжруулах зөвлөмж</CardTitle></CardHeader>
        <CardContent>
          <ul class="space-y-2">
            <li v-for="rec in application.recommendations" :key="rec"
                class="flex items-start gap-2 text-sm text-muted-foreground">
              <ChevronRight class="mt-0.5 h-4 w-4 shrink-0 text-primary" /> {{ rec }}
            </li>
          </ul>
        </CardContent>
      </Card>
    </div>
  </div>

  <!-- Apply dialog -->
  <Dialog v-model:open="applyOpen">
    <DialogContent class="rounded-3xl sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Анкет илгээх</DialogTitle>
        <DialogDescription>
          CV файлаа PDF хэлбэрээр оруулна уу. AI таны CV-г ажлын байрны шаардлагатай харьцуулж үнэлнэ.
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-2">
        <div class="space-y-2">
          <Label>CV файл (PDF)</Label>
          <Input type="file" accept=".pdf" @change="onFileChange" />
          <p v-if="cvFile" class="text-xs text-muted-foreground">Сонгосон: {{ cvFile.name }}</p>
        </div>
      </div>
      <DialogFooter>
        <Button variant="outline" @click="applyOpen = false" :disabled="isSubmitting">Болих</Button>
        <Button @click="submitApplication" :disabled="isSubmitting || !cvFile">
          <Loader2 v-if="isSubmitting" class="mr-2 h-4 w-4 animate-spin" />
          {{ isSubmitting ? "Илгээж байна..." : "Илгээх" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
