<script setup lang="ts">
import { JobStatus } from "../composables/types";
import { toast } from "vue-sonner";

const { user } = useAuth();
const jobsAPI = useJobsAPI();
const router = useRouter();

const recruiterJobs = ref<Job[]>([]);
const loading = ref(false);

const recruiterStats = computed(() => {
  const totalJobs = recruiterJobs.value.length;
  const activeJobs = recruiterJobs.value.filter((job) => job.status === JobStatus.Posted).length;
  const totalApplicants = recruiterJobs.value.reduce(
    (sum, job) => sum + job.applicants_count,
    0,
  );

  return { totalJobs, activeJobs, totalApplicants };
});

const recruiterActivity = computed(() => {
  return recruiterJobs.value.slice(0, 3).map((job) => ({
    title: job.title,
    subtitle: `${job.company_name || "Your company"} · ${job.location}`,
    status: job.status,
  }));
});

const applicantHighlights = [
  {
    title: "Browse matched roles",
    description: "See recruiter posts in a cleaner workspace without changing your auth flow.",
  },
  {
    title: "Track your readiness",
    description: "Keep assessments, profile details, and job discovery grouped in one shell.",
  },
  {
    title: "Stay interview-ready",
    description: "Move from application to review with a simpler page structure.",
  },
];

const recruiterNeedsCompany = computed(
  () => user.value?.role === "recruiter" && !user.value?.company_id,
);

function openRecruiterJobs() {
  if (recruiterNeedsCompany.value) {
    toast.warning("Add your company first before managing recruiter jobs.");
    router.push("/profile");
    return;
  }

  router.push("/jobs");
}

async function loadRecruiterDashboard() {
  if (user.value?.role !== "recruiter") return;

  loading.value = true;
  try {
    recruiterJobs.value = await jobsAPI.list();
  } finally {
    loading.value = false;
  }
}

await loadRecruiterDashboard();
</script>

<template>
  <div v-if="user?.role === 'recruiter'" class="space-y-6">
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p class="text-sm font-medium text-muted-foreground">Recruiter dashboard</p>
          <h2 class="mt-2 text-3xl font-semibold tracking-tight">
            Welcome back, {{ user?.name }}
          </h2>
          <p class="mt-2 text-sm leading-6 text-muted-foreground">
            Your job posting workspace now follows the wireframe structure and is connected to live backend data.
          </p>
        </div>
        <Button class="rounded-full px-5" @click="openRecruiterJobs">
          Manage job posts
        </Button>
      </div>
    </section>

    <section class="grid gap-4 md:grid-cols-3">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Total jobs posted</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.totalJobs }}</CardTitle>
        </CardHeader>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Active listings</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.activeJobs }}</CardTitle>
        </CardHeader>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Applicants tracked</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.totalApplicants }}</CardTitle>
        </CardHeader>
      </Card>
    </section>

    <section class="grid gap-4 lg:grid-cols-[1.3fr_0.9fr]">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Job momentum</CardTitle>
          <CardDescription>
            {{ loading ? "Refreshing recruiter jobs..." : "Latest recruiter-owned postings from the backend." }}
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <div
            v-for="job in recruiterJobs.slice(0, 4)"
            :key="job.id"
            class="rounded-2xl border border-border bg-muted/20 p-4"
          >
            <div class="flex items-start justify-between gap-4">
              <div>
                <p class="text-sm font-semibold">{{ job.title }}</p>
                <p class="mt-1 text-sm text-muted-foreground">
                  {{ job.location }} · {{ job.employment_type || "Not set" }}
                </p>
              </div>
              <Badge variant="secondary" class="rounded-full px-3 py-1 capitalize">
                {{ job.status }}
              </Badge>
            </div>
          </div>
          <div
            v-if="!recruiterJobs.length && !loading"
            class="rounded-2xl border border-dashed border-border px-4 py-10 text-center text-sm text-muted-foreground"
          >
            Create your first job post to populate the recruiter dashboard.
          </div>
        </CardContent>
      </Card>

      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Recent activity</CardTitle>
          <CardDescription>Wireframe-style summary without inbox pages.</CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <div
            v-for="item in recruiterActivity"
            :key="item.title"
            class="rounded-2xl border border-border bg-muted/20 p-4"
          >
            <p class="text-sm font-semibold">{{ item.title }}</p>
            <p class="mt-1 text-sm text-muted-foreground">{{ item.subtitle }}</p>
            <Badge variant="outline" class="mt-3 rounded-full px-3 py-1 capitalize">
              {{ item.status }}
            </Badge>
          </div>
          <div
            v-if="!recruiterActivity.length && !loading"
            class="rounded-2xl border border-dashed border-border px-4 py-10 text-center text-sm text-muted-foreground"
          >
            Recent activity appears here once jobs are created.
          </div>
        </CardContent>
      </Card>
    </section>
  </div>

  <div v-else class="space-y-6">
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <p class="text-sm font-medium text-muted-foreground">Applicant dashboard</p>
      <h2 class="mt-2 text-3xl font-semibold tracking-tight">
        Welcome back, {{ user?.name }}
      </h2>
      <p class="mt-2 max-w-2xl text-sm leading-6 text-muted-foreground">
        The main workspace now follows the wireframe direction while leaving your core applicant logic ready for future integrations.
      </p>
    </section>

    <section class="grid gap-4 md:grid-cols-3">
      <Card
        v-for="item in applicantHighlights"
        :key="item.title"
        class="rounded-3xl border-border shadow-sm"
      >
        <CardHeader>
          <CardTitle class="text-xl">{{ item.title }}</CardTitle>
          <CardDescription class="leading-6">{{ item.description }}</CardDescription>
        </CardHeader>
      </Card>
    </section>
  </div>
</template>
