<script setup lang="ts">
import type { Application } from "../composables/types";
import { toast } from "vue-sonner";
import { FileText, Loader2, ChevronRight } from "lucide-vue-next";

definePageMeta({ middleware: "auth" });

const applicationsAPI = useApplicationsAPI();
const { user } = useAuth();
const router = useRouter();

const applications = ref<Application[]>([]);
const loading = ref(true);

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

function scoreColor(score: number) {
  if (score >= 75) return "text-green-600";
  if (score >= 50) return "text-yellow-600";
  return "text-red-600";
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

function statusVariant(status: Application["status"]) {
  if (status === "shortlisted") return "default";
  if (status === "rejected") return "destructive";
  if (status === "pending") return "outline";
  return "secondary";
}

function formatDate(d: string) {
  return new Intl.DateTimeFormat("mn-MN", { dateStyle: "medium" }).format(new Date(d));
}
</script>

<template>
  <div class="space-y-6">
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <p class="text-sm font-medium text-muted-foreground">Миний анкетууд</p>
      <h2 class="mt-2 text-3xl font-semibold tracking-tight">Илгээсэн анкетууд</h2>
      <p class="mt-2 text-sm leading-6 text-muted-foreground">
        Та илгээсэн анкетуудын AI үнэлгээний үр дүнг энд харж болно.
      </p>
    </section>

    <div v-if="loading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <div v-else-if="!applications.length"
         class="rounded-3xl border border-dashed border-border px-6 py-20 text-center">
      <FileText class="mx-auto h-10 w-10 text-muted-foreground" />
      <p class="mt-4 text-lg font-medium">Илгээсэн анкет байхгүй</p>
      <p class="mt-2 text-sm text-muted-foreground">Ажлын байруудыг харж анкетаа илгээгээрэй.</p>
      <Button class="mt-6 rounded-full px-5" @click="router.push('/jobs')">
        Ажлын байр харах
      </Button>
    </div>

    <div v-else class="space-y-4">
      <div v-for="app in applications" :key="app.id"
           class="cursor-pointer rounded-3xl border border-border bg-card px-6 py-5 shadow-sm transition hover:shadow-md"
           @click="router.push(`/jobs/${app.job_posting_id}`)">
        <div class="flex items-start justify-between gap-4">
          <div class="space-y-1">
            <h3 class="font-semibold">{{ app.job_title || "Ажлын байр" }}</h3>
            <p class="text-xs text-muted-foreground">{{ formatDate(app.created_at) }}</p>
          </div>
          <div class="flex items-center gap-3">
            <div v-if="app.status !== 'pending'" class="text-right">
              <p class="text-xs text-muted-foreground">Оноо</p>
              <p :class="['text-2xl font-bold', scoreColor(app.overall_score)]">{{ app.overall_score }}%</p>
            </div>
            <Badge :variant="statusVariant(app.status) as any" class="rounded-full px-3 py-1">
              {{ statusLabel(app.status) }}
            </Badge>
            <ChevronRight class="h-4 w-4 text-muted-foreground" />
          </div>
        </div>
        <p v-if="app.summary && app.status !== 'pending'" class="mt-3 text-sm text-muted-foreground line-clamp-2">
          {{ app.summary }}
        </p>
        <div v-if="app.status === 'pending'" class="mt-3 flex items-center gap-2 text-sm text-muted-foreground">
          <Loader2 class="h-3 w-3 animate-spin" /> AI үнэлгээ хийгдэж байна...
        </div>
      </div>
    </div>
  </div>
</template>
