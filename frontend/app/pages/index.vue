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
    title: "Тохирох ажлыг олоорой",
    description: "Ажил олгогчдийн зарласан нээлттэй байр дундаас өөрт тохирохыг нь олоорой.",
  },
  {
    title: "Бэлэн байдлаа хянаарай",
    description: "Үнэлгээ, профайл мэдээлэл болон ажлын нээлтийг нэг дор хадгалаарай.",
  },
  {
    title: "Ярилцлагад бэлэн байгаарай",
    description: "Анкетаас эхлээд дүгнэлт хүртэлх бүх шатыг энэ платформоор дамжуулан хянаарай.",
  },
];

const recruiterNeedsCompany = computed(
  () => user.value?.role === "recruiter" && !user.value?.company_id,
);

function openRecruiterJobs() {
  if (recruiterNeedsCompany.value) {
    toast.warning("Ажлын байр удирдахаас өмнө компанийн мэдээллээ бүртгэнэ үү.");
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
          <p class="text-sm font-medium text-muted-foreground">Ажил олгогчийн самбар</p>
          <h2 class="mt-2 text-3xl font-semibold tracking-tight">
            Тавтай морилно уу, {{ user?.name }}
          </h2>
          <p class="mt-2 text-sm leading-6 text-muted-foreground">
            Ажлын байрны зарыг удирдах, нэр дэвшигчидтэй ажиллах ажлын орчин.
          </p>
        </div>
        <Button class="rounded-full px-5" @click="openRecruiterJobs">
          Ажлын байр удирдах
        </Button>
      </div>
    </section>

    <section class="grid gap-4 md:grid-cols-3">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Нийт ажлын байр</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.totalJobs }}</CardTitle>
        </CardHeader>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Идэвхтэй зар</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.activeJobs }}</CardTitle>
        </CardHeader>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardDescription>Бүртгэгдсэн горилогч</CardDescription>
          <CardTitle class="text-4xl">{{ recruiterStats.totalApplicants }}</CardTitle>
        </CardHeader>
      </Card>
    </section>

    <section class="grid gap-4 lg:grid-cols-[1.3fr_0.9fr]">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Ажлын байрны хөдөлгөөн</CardTitle>
          <CardDescription>
            {{ loading ? "Ачааллаж байна..." : "Таны сүүлийн ажлын байрны зарууд." }}
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
                  {{ job.location }} · {{ job.employment_type || "Тогтоогдоогүй" }}
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
            Самбарыг дүүргэхийн тулд эхний ажлын зараа оруулна уу.
          </div>
        </CardContent>
      </Card>

      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Сүүлийн үйл ажиллагаа</CardTitle>
          <CardDescription>Ажлын байрны сүүлийн хөдөлгөөний товч тойм.</CardDescription>
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
            Ажлын байр үүсгэсний дараа энд үйл ажиллагаа харагдана.
          </div>
        </CardContent>
      </Card>
    </section>
  </div>

  <div v-else class="space-y-6">
    <section class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm">
      <p class="text-sm font-medium text-muted-foreground">Ажил горилогчийн самбар</p>
      <h2 class="mt-2 text-3xl font-semibold tracking-tight">
        Тавтай морилно уу, {{ user?.name }}
      </h2>
      <p class="mt-2 max-w-2xl text-sm leading-6 text-muted-foreground">
        Анкетуудаа хянаж, тохирох ажлыг олж, дараагийн алхамдаа бэлэн байгаарай.
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
