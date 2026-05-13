<script setup lang="ts">
import type { Job } from "../../composables/types";
import { JobStatus } from "../../composables/types";
import { toast } from "vue-sonner";
import {
  BriefcaseBusiness,
  Eye,
  Plus,
  Trash2,
  Users,
} from "lucide-vue-next";

const { user } = useAuth();
const jobsAPI = useJobsAPI();
const router = useRouter();

const jobs = ref<Job[]>([]);
const selectedJob = ref<Job | null>(null);
const loading = ref(false);
const deleteDialogOpen = ref(false);
const isDeleting = ref(false);

const isRecruiter = computed(() => user.value?.role === "recruiter");
const recruiterCompanyID = computed(() => user.value?.company_id ?? 0);
const hasRecruiterCompany = computed(() => recruiterCompanyID.value > 0);

const recruiterStats = computed(() => {
  const totalJobs = jobs.value.length;
  const activeJobs = jobs.value.filter((job) => job.status === JobStatus.Posted).length;
  const totalApplicants = jobs.value.reduce((sum, job) => sum + job.applicants_count, 0);
  return { totalJobs, activeJobs, totalApplicants };
});

const applicantJobs = computed(() =>
  jobs.value.filter((job) => job.status === JobStatus.Posted),
);

function statusLabel(status: Job["status"]) {
  if (status === JobStatus.Posted) return "Нийтлэгдсэн";
  if (status === JobStatus.Closed) return "Хаагдсан";
  return "Ноорог";
}

function statusVariant(status: Job["status"]) {
  if (status === JobStatus.Posted) return "default";
  if (status === JobStatus.Closed) return "secondary";
  return "outline";
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("mn-MN", { dateStyle: "medium" }).format(new Date(value));
}

function openCreatePage() {
  if (!hasRecruiterCompany.value) {
    toast.warning("Ажлын байр нэмэхийн өмнө компанийн мэдээллээ бүртгэнэ үү.");
    return;
  }
  router.push("/jobs/new");
}

function openDeleteDialog(job: Job) {
  selectedJob.value = job;
  deleteDialogOpen.value = true;
}

async function loadJobs() {
  loading.value = true;
  try {
    if (isRecruiter.value && hasRecruiterCompany.value) {
      jobs.value = await jobsAPI.listByCompany(recruiterCompanyID.value);
      return;
    }
    jobs.value = await jobsAPI.list();
  } finally {
    loading.value = false;
  }
}

async function confirmDelete() {
  if (!selectedJob.value) return;

  isDeleting.value = true;
  try {
    await jobsAPI.delete(recruiterCompanyID.value, selectedJob.value.id);
    deleteDialogOpen.value = false;
    selectedJob.value = null;
    await loadJobs();
  } catch (error: any) {
    toast.error(error?.data?.message || error?.message || "Ажлын байр устгахад алдаа гарлаа.");
  } finally {
    isDeleting.value = false;
  }
}

await loadJobs();
</script>

<template>
  <div v-if="isRecruiter" class="space-y-6">
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p class="text-sm font-medium text-muted-foreground">Ажил олгогчийн ажлын байрууд</p>
          <h2 class="mt-2 text-3xl font-semibold tracking-tight">Ажлын байрны удирдлага</h2>
          <p class="mt-2 text-sm leading-6 text-muted-foreground">
            Нийтлэгдсэн болон ноорог ажлын зарыг удирдаж, нэр дэвшигчдийн ирсэн анкетыг хянана уу.
          </p>
        </div>
        <Button class="rounded-full px-5" @click="openCreatePage">
          <Plus class="mr-2 h-4 w-4" />
          Ажлын байр нэмэх
        </Button>
      </div>
    </section>

    <div
      v-if="!hasRecruiterCompany"
      class="rounded-3xl border border-dashed border-border bg-muted/20 px-6 py-4 text-sm text-muted-foreground"
    >
      Ажлын байр нэмэх эсвэл удирдахаас өмнө профайл хэсгээс компанийн мэдээллээ бүртгэнэ үү.
    </div>

    <section class="grid gap-4 md:grid-cols-3">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Нийт ажлын байр</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.totalJobs }}</CardTitle>
        </CardHeader>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Идэвхтэй ажлын байр</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.activeJobs }}</CardTitle>
        </CardHeader>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Нийт горилогч</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.totalApplicants }}</CardTitle>
        </CardHeader>
      </Card>
    </section>

    <Card class="rounded-3xl border-border shadow-sm">
      <CardHeader class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <CardTitle>Таны ажлын байрууд</CardTitle>
          <CardDescription>
            {{ loading ? "Ачааллаж байна..." : `${jobs.length} ажлын байр байна.` }}
          </CardDescription>
        </div>
      </CardHeader>
      <CardContent>
        <div v-if="jobs.length" class="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Ажлын байрны нэр</TableHead>
                <TableHead>Төлөв</TableHead>
                <TableHead>Горилогч</TableHead>
                <TableHead>Үүсгэсэн</TableHead>
                <TableHead>Шинэчлэгдсэн</TableHead>
                <TableHead class="text-right">Үйлдэл</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow
                v-for="job in jobs"
                :key="job.id"
                class="cursor-pointer"
                @click="router.push(`/jobs/${job.id}`)"
              >
                <TableCell class="font-medium">{{ job.title }}</TableCell>
                <TableCell>
                  <Badge
                    :variant="statusVariant(job.status) as any"
                    class="rounded-full px-3 py-1 capitalize"
                  >
                    {{ statusLabel(job.status) }}
                  </Badge>
                </TableCell>
                <TableCell>{{ job.applicants_count }}</TableCell>
                <TableCell class="text-muted-foreground">{{ formatDate(job.created_at) }}</TableCell>
                <TableCell class="text-muted-foreground">{{ formatDate(job.updated_at) }}</TableCell>
                <TableCell class="text-right">
                  <div class="flex justify-end gap-2">
                    <NuxtLink :to="`/jobs/${job.id}/applicants`" @click.stop>
                      <Button variant="outline" size="icon" class="rounded-xl" title="Ирсэн анкетууд">
                        <Users class="h-4 w-4" />
                      </Button>
                    </NuxtLink>
                    <NuxtLink :to="`/jobs/${job.id}`" @click.stop>
                      <Button variant="outline" size="icon" class="rounded-xl" title="Харах / Засах">
                        <Eye class="h-4 w-4" />
                      </Button>
                    </NuxtLink>
                    <Button
                      variant="outline"
                      size="icon"
                      class="rounded-xl"
                      @click.stop="openDeleteDialog(job)"
                    >
                      <Trash2 class="h-4 w-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>

        <div
          v-else-if="!loading"
          class="rounded-3xl border border-dashed border-border px-6 py-16 text-center"
        >
          <BriefcaseBusiness class="mx-auto h-10 w-10 text-muted-foreground" />
          <p class="mt-4 text-lg font-medium">Ажлын байр байхгүй</p>
          <p class="mt-2 text-sm text-muted-foreground">
            Ажилд авах урсгалаа эхлүүлэхийн тулд эхний ажлын байраа нэмнэ үү.
          </p>
          <Button class="mt-6 rounded-full px-5" @click="openCreatePage">
            <Plus class="mr-2 h-4 w-4" />
            Ажлын байр нэмэх
          </Button>
        </div>
      </CardContent>
    </Card>

    <AlertDialog v-model:open="deleteDialogOpen">
      <AlertDialogContent class="rounded-3xl">
        <AlertDialogHeader>
          <AlertDialogTitle>Ажлын байр устгах уу?</AlertDialogTitle>
          <AlertDialogDescription>
            <span class="font-medium text-foreground">{{ selectedJob?.title }}</span>
            ажлын байрыг таны ажлын орчноос устгана.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel :disabled="isDeleting">Болих</AlertDialogCancel>
          <AlertDialogAction :disabled="isDeleting" @click="confirmDelete">
            {{ isDeleting ? "Устгаж байна..." : "Устгах" }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>

  <div v-else class="space-y-6">
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <p class="text-sm font-medium text-muted-foreground">Ажлын байрны тохирол</p>
      <h2 class="mt-2 text-3xl font-semibold tracking-tight">Нээлттэй ажлын байрууд</h2>
      <p class="mt-2 text-sm leading-6 text-muted-foreground">
        Идэвхтэй нийтлэгдсэн ажлын байруудыг харж, анкетаа илгээгээрэй.
      </p>
    </section>

    <section class="grid gap-4 lg:grid-cols-2">
      <Card
        v-for="job in applicantJobs"
        :key="job.id"
        class="rounded-3xl border-border shadow-sm"
      >
        <CardHeader class="space-y-3">
          <div class="flex items-start justify-between gap-4">
            <div>
              <CardTitle class="text-2xl">{{ job.title }}</CardTitle>
              <CardDescription class="mt-1">{{ job.company_name || "Компани" }}</CardDescription>
            </div>
            <Badge class="rounded-full px-3 py-1 capitalize">{{ statusLabel(job.status) }}</Badge>
          </div>
          <div class="flex flex-wrap gap-2">
            <Badge variant="outline" class="rounded-full px-3 py-1">{{ job.location }}</Badge>
            <Badge variant="outline" class="rounded-full px-3 py-1">
              {{ job.employment_type || "Тогтоогдоогүй" }}
            </Badge>
            <Badge variant="outline" class="rounded-full px-3 py-1">
              {{ job.seniority || "Тогтоогдоогүй" }}
            </Badge>
          </div>
        </CardHeader>
        <CardContent class="space-y-3">
          <p class="line-clamp-3 text-sm leading-6 text-muted-foreground">{{ job.description }}</p>
          <NuxtLink :to="`/jobs/${job.id}`">
            <Button class="w-full rounded-full" variant="outline">Дэлгэрэнгүй харах</Button>
          </NuxtLink>
        </CardContent>
      </Card>

      <div
        v-if="!applicantJobs.length && !loading"
        class="rounded-3xl border border-dashed border-border px-6 py-16 text-center text-sm text-muted-foreground lg:col-span-2"
      >
        Идэвхтэй ажлын байр байхгүй байна.
      </div>
    </section>
  </div>
</template>
