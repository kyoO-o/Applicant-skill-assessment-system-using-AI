<script setup lang="ts">
import type { Job } from "../../composables/types";
import type { SaveJobPayload } from "../../composables/types/payload";
import { JobStatus } from "../../composables/types";
import {
  BriefcaseBusiness,
  Eye,
  Pencil,
  Plus,
  Trash2,
  MapPin,
  Users,
} from "lucide-vue-next";

const { user } = useAuth();
const jobsAPI = useJobsAPI();

const jobs = ref<Job[]>([]);
const selectedJob = ref<Job | null>(null);
const loading = ref(false);
const dialogOpen = ref(false);
const detailsOpen = ref(false);
const deleteDialogOpen = ref(false);
const isSubmitting = ref(false);
const errorMessage = ref("");

const form = reactive({
  id: 0,
  title: "",
  location: "",
  employment_type: "Full-time",
  seniority: "Mid-level",
  status: JobStatus.Draft,
  description: "",
  requirements_text: "",
});

const isRecruiter = computed(() => user.value?.role === "recruiter");
const isEditing = computed(() => form.id > 0);
const recruiterCompanyID = computed(() => user.value?.company_id ?? 0);
const hasRecruiterCompany = computed(() => recruiterCompanyID.value > 0);

const recruiterStats = computed(() => {
  const totalJobs = jobs.value.length;
  const activeJobs = jobs.value.filter((job) => job.status === JobStatus.Posted).length;
  const totalApplicants = jobs.value.reduce(
    (sum, job) => sum + job.applicants_count,
    0,
  );

  return { totalJobs, activeJobs, totalApplicants };
});

const applicantJobs = computed(() => jobs.value.filter((job) => job.status === JobStatus.Posted));

function statusLabel(status: Job["status"]) {
  if (status === JobStatus.Posted) return "Posted";
  if (status === JobStatus.Closed) return "Closed";
  return "Draft";
}

function statusVariant(status: Job["status"]) {
  if (status === JobStatus.Posted) return "default";
  if (status === JobStatus.Closed) return "secondary";
  return "outline";
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  }).format(new Date(value));
}

function resetForm() {
  form.id = 0;
  form.title = "";
  form.location = "";
  form.employment_type = "Full-time";
  form.seniority = "Mid-level";
  form.status = JobStatus.Draft;
  form.description = "";
  form.requirements_text = "";
}

function openCreateDialog() {
  errorMessage.value = "";
  if (!hasRecruiterCompany.value) {
    errorMessage.value = "Create your company profile before adding a job post.";
    return;
  }
  resetForm();
  dialogOpen.value = true;
}

function openEditDialog(job: Job) {
  errorMessage.value = "";
  form.id = job.id;
  form.title = job.title;
  form.location = job.location;
  form.employment_type = job.employment_type || "Full-time";
  form.seniority = job.seniority || "Mid-level";
  form.status = job.status || JobStatus.Draft;
  form.description = job.description;
  form.requirements_text = job.requirements.join("\n");
  dialogOpen.value = true;
}

function openDetails(job: Job) {
  selectedJob.value = job;
  detailsOpen.value = true;
}

function openDeleteDialog(job: Job) {
  selectedJob.value = job;
  deleteDialogOpen.value = true;
}

function buildPayload(): SaveJobPayload {
  return {
    title: form.title.trim(),
    location: form.location.trim(),
    employment_type: form.employment_type.trim(),
    seniority: form.seniority.trim(),
    status: form.status,
    description: form.description.trim(),
    requirements: form.requirements_text
      .split("\n")
      .map((item) => item.trim())
      .filter(Boolean),
  };
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

async function submitJob() {
  errorMessage.value = "";
  isSubmitting.value = true;

  try {
    const payload = buildPayload();

    if (isEditing.value) {
      if (!hasRecruiterCompany.value) {
        throw new Error("A company profile is required to update a job.");
      }
      await jobsAPI.update(recruiterCompanyID.value, form.id, payload);
    } else {
      await jobsAPI.create(payload);
    }

    dialogOpen.value = false;
    await loadJobs();
  } catch (error: any) {
    errorMessage.value =
      error?.data?.message || error?.message || "Unable to save the job.";
  } finally {
    isSubmitting.value = false;
  }
}

async function confirmDelete() {
  if (!selectedJob.value) return;

  isSubmitting.value = true;
  try {
    if (!hasRecruiterCompany.value) {
      throw new Error("A company profile is required to delete a job.");
    }
    await jobsAPI.delete(recruiterCompanyID.value, selectedJob.value.id);
    deleteDialogOpen.value = false;
    detailsOpen.value = false;
    selectedJob.value = null;
    await loadJobs();
  } catch (error: any) {
    errorMessage.value =
      error?.data?.message || error?.message || "Unable to delete the job.";
  } finally {
    isSubmitting.value = false;
  }
}

await loadJobs();
</script>

<template>
  <div v-if="isRecruiter" class="space-y-6">
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p class="text-sm font-medium text-muted-foreground">Recruiter jobs</p>
          <h2 class="mt-2 text-3xl font-semibold tracking-tight">
            Job management
          </h2>
          <p class="mt-2 text-sm leading-6 text-muted-foreground">
            This mirrors the wireframe recruiter table pattern and is connected
            to your backend CRUD instead of static data.
          </p>
        </div>
        <Button
          class="rounded-full px-5"
          :disabled="!hasRecruiterCompany"
          @click="openCreateDialog"
        >
          <Plus class="mr-2 h-4 w-4" />
          Create Job
        </Button>
      </div>
    </section>

    <div
      v-if="!hasRecruiterCompany"
      class="rounded-3xl border border-dashed border-border bg-muted/20 px-6 py-4 text-sm text-muted-foreground"
    >
      Add your company profile first from the company section before creating or managing job posts.
    </div>

    <section class="grid gap-4 md:grid-cols-3">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Total jobs posted</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.totalJobs }}</CardTitle>
        </CardHeader>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Active jobs</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.activeJobs }}</CardTitle>
        </CardHeader>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Total applicants</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.totalApplicants }}</CardTitle>
        </CardHeader>
      </Card>
    </section>

    <Card class="rounded-3xl border-border shadow-sm">
      <CardHeader class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <CardTitle>Your job posts</CardTitle>
          <CardDescription>
            {{ loading ? "Loading recruiter jobs..." : `${jobs.length} jobs connected from the backend.` }}
          </CardDescription>
        </div>
      </CardHeader>
      <CardContent>
        <div v-if="jobs.length" class="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Job Title</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Applicants</TableHead>
                <TableHead>Created</TableHead>
                <TableHead>Updated</TableHead>
                <TableHead class="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="job in jobs" :key="job.id" class="cursor-pointer" @click="openDetails(job)">
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
                <TableCell class="text-muted-foreground">
                  {{ formatDate(job.created_at) }}
                </TableCell>
                <TableCell class="text-muted-foreground">
                  {{ formatDate(job.updated_at) }}
                </TableCell>
                <TableCell class="text-right">
                  <div class="flex justify-end gap-2">
                    <Button
                      variant="outline"
                      size="icon"
                      class="rounded-xl"
                      @click.stop="openDetails(job)"
                    >
                      <Eye class="h-4 w-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="icon"
                      class="rounded-xl"
                      @click.stop="openEditDialog(job)"
                    >
                      <Pencil class="h-4 w-4" />
                    </Button>
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
          <p class="mt-4 text-lg font-medium">No job posts yet</p>
          <p class="mt-2 text-sm text-muted-foreground">
            Create your first role to start using the recruiter pipeline.
          </p>
          <Button
            class="mt-6 rounded-full px-5"
            :disabled="!hasRecruiterCompany"
            @click="openCreateDialog"
          >
            <Plus class="mr-2 h-4 w-4" />
            Create Job
          </Button>
        </div>
      </CardContent>
    </Card>

    <Dialog v-model:open="dialogOpen">
      <DialogContent class="max-w-3xl rounded-3xl">
        <DialogHeader>
          <DialogTitle>
            {{ isEditing ? "Edit job post" : "Create a new job post" }}
          </DialogTitle>
          <DialogDescription>
            Keep the structure from the wireframe and save directly to the backend.
          </DialogDescription>
        </DialogHeader>

        <form class="grid gap-4 py-2" @submit.prevent="submitJob">
          <div
            v-if="errorMessage"
            class="rounded-2xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
          >
            {{ errorMessage }}
          </div>

          <div class="grid gap-4 sm:grid-cols-2">
            <div class="space-y-2">
              <Label for="job-title">Job title</Label>
              <Input id="job-title" v-model="form.title" placeholder="Frontend Developer" />
            </div>
            <div class="space-y-2">
              <Label for="job-location">Location</Label>
              <Input id="job-location" v-model="form.location" placeholder="Ulaanbaatar, MN" />
            </div>
          </div>

          <div class="grid gap-4 sm:grid-cols-3">
            <div class="space-y-2">
              <Label for="job-type">Employment type</Label>
              <Input id="job-type" v-model="form.employment_type" placeholder="Full-time" />
            </div>
            <div class="space-y-2">
              <Label for="job-level">Seniority</Label>
              <Input id="job-level" v-model="form.seniority" placeholder="Mid-level" />
            </div>
            <div class="space-y-2">
              <Label for="job-status">Status</Label>
              <Select v-model="form.status">
                <SelectTrigger id="job-status" class="w-full">
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem :value="JobStatus.Draft">Draft</SelectItem>
                  <SelectItem :value="JobStatus.Posted">Posted</SelectItem>
                  <SelectItem :value="JobStatus.Closed">Closed</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div class="space-y-2">
            <Label for="job-description">Job description</Label>
            <Textarea
              id="job-description"
              v-model="form.description"
              rows="5"
              placeholder="Describe the role, scope, and outcomes."
            />
          </div>

          <div class="space-y-2">
            <Label for="job-requirements">Requirements</Label>
            <Textarea
              id="job-requirements"
              v-model="form.requirements_text"
              rows="6"
              placeholder="One requirement per line"
            />
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" @click="dialogOpen = false">
              Cancel
            </Button>
            <Button type="submit" :disabled="isSubmitting">
              {{ isSubmitting ? "Saving..." : isEditing ? "Update Job" : "Create Job" }}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="detailsOpen">
      <DialogContent class="max-w-3xl rounded-3xl">
        <DialogHeader>
          <DialogTitle>{{ selectedJob?.title }}</DialogTitle>
          <DialogDescription>
            {{ selectedJob?.company_name || user?.company_name || "Your company" }}
          </DialogDescription>
        </DialogHeader>

        <div v-if="selectedJob" class="space-y-6 py-2">
          <div class="flex flex-wrap gap-2">
            <Badge class="rounded-full px-3 py-1 capitalize">
              {{ statusLabel(selectedJob.status) }}
            </Badge>
            <Badge variant="outline" class="rounded-full px-3 py-1">
              <MapPin class="mr-1 h-3.5 w-3.5" />
              {{ selectedJob.location }}
            </Badge>
            <Badge variant="outline" class="rounded-full px-3 py-1">
              {{ selectedJob.employment_type || "Not set" }}
            </Badge>
            <Badge variant="outline" class="rounded-full px-3 py-1">
              {{ selectedJob.seniority || "Not set" }}
            </Badge>
            <Badge variant="outline" class="rounded-full px-3 py-1">
              <Users class="mr-1 h-3.5 w-3.5" />
              {{ selectedJob.applicants_count }} applicants
            </Badge>
          </div>

          <div class="rounded-2xl border border-border bg-muted/20 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Description
            </p>
            <p class="mt-3 whitespace-pre-line text-sm leading-6">
              {{ selectedJob.description }}
            </p>
          </div>

          <div class="rounded-2xl border border-border bg-muted/20 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Requirements
            </p>
            <ul class="mt-3 space-y-2 text-sm leading-6">
              <li v-for="requirement in selectedJob.requirements" :key="requirement">
                {{ requirement }}
              </li>
            </ul>
          </div>

          <DialogFooter>
            <Button variant="outline" @click="openEditDialog(selectedJob)">
              <Pencil class="mr-2 h-4 w-4" />
              Edit Job
            </Button>
            <Button variant="outline" @click="openDeleteDialog(selectedJob)">
              <Trash2 class="mr-2 h-4 w-4" />
              Delete Job
            </Button>
          </DialogFooter>
        </div>
      </DialogContent>
    </Dialog>

    <AlertDialog v-model:open="deleteDialogOpen">
      <AlertDialogContent class="rounded-3xl">
        <AlertDialogHeader>
          <AlertDialogTitle>Delete this job?</AlertDialogTitle>
          <AlertDialogDescription>
            This removes <span class="font-medium text-foreground">{{ selectedJob?.title }}</span>
            from your recruiter workspace.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel :disabled="isSubmitting">Cancel</AlertDialogCancel>
          <AlertDialogAction :disabled="isSubmitting" @click="confirmDelete">
            {{ isSubmitting ? "Deleting..." : "Delete Job" }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>

  <div v-else class="space-y-6">
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <p class="text-sm font-medium text-muted-foreground">Job matches</p>
      <h2 class="mt-2 text-3xl font-semibold tracking-tight">
        Browse open roles
      </h2>
      <p class="mt-2 text-sm leading-6 text-muted-foreground">
        Applicants now see the same shell language, while the job data is read
        from the live backend endpoint.
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
              <CardDescription class="mt-1">
                {{ job.company_name || "Company" }}
              </CardDescription>
            </div>
            <Badge class="rounded-full px-3 py-1 capitalize">
              {{ statusLabel(job.status) }}
            </Badge>
          </div>
          <div class="flex flex-wrap gap-2">
            <Badge variant="outline" class="rounded-full px-3 py-1">
              {{ job.location }}
            </Badge>
            <Badge variant="outline" class="rounded-full px-3 py-1">
              {{ job.employment_type || "Not set" }}
            </Badge>
            <Badge variant="outline" class="rounded-full px-3 py-1">
              {{ job.seniority || "Not set" }}
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          <p class="text-sm leading-6 text-muted-foreground">
            {{ job.description }}
          </p>
        </CardContent>
      </Card>

      <div
        v-if="!applicantJobs.length && !loading"
        class="rounded-3xl border border-dashed border-border px-6 py-16 text-center text-sm text-muted-foreground lg:col-span-2"
      >
        No active jobs are available yet.
      </div>
    </section>
  </div>
</template>
