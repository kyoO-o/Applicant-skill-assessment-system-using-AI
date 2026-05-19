<script setup lang="ts">
import { toast } from "vue-sonner";
import type { Company } from "../../composables/types";

definePageMeta({ middleware: "auth" });

const { user } = useAuth();
const jobsAPI = useJobsAPI();
const companyAPI = useCompanyAPI();
const router = useRouter();

const company = ref<Company | null>(null);
const isSubmitting = ref(false);
const errorMessage = ref("");

const isRecruiter = computed(() => user.value?.role === "recruiter");

if (!isRecruiter.value) {
  router.replace("/jobs");
}

try {
  if (user.value?.company_id) {
    company.value = await companyAPI.get();
  }
} catch {
  /* company not yet set up */
}

async function handleSubmit(payload: any) {
  errorMessage.value = "";
  isSubmitting.value = true;
  try {
    await jobsAPI.create(payload);
    toast.success("Ажлын байр амжилттай үүсгэгдлээ.");
    router.push("/jobs");
  } catch (error: any) {
    errorMessage.value =
      error?.data?.message ||
      error?.message ||
      "Ажлын байр хадгалахад алдаа гарлаа.";
  } finally {
    isSubmitting.value = false;
  }
}
</script>

<template>
  <div class="space-y-6">
    <section class="">
      <!-- <div class="flex items-center gap-3">
        <NuxtLink
          to="/jobs"
          class="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
        >
          ← Ажлын байрууд
        </NuxtLink>
      </div> -->
      <h1 class="text-[26px] font-semibold tracking-tight">
        Шинэ ажлын байр нэмэх
      </h1>
      <p class="mt-2 text-sm leading-6 text-muted-foreground">
        Ажлын байрны мэдээллийг бөглөж хадгална уу.
      </p>
    </section>

    <Card class="rounded-3xl border-border">
      <CardContent class="pt-6">
        <JobForm
          :company="company"
          :is-submitting="isSubmitting"
          :error-message="errorMessage"
          @submit="handleSubmit"
        >
          <template #actions>
            <Button
              type="button"
              variant="outline"
              @click="router.push('/jobs')"
              >Болих</Button
            >
            <Button type="submit" :disabled="isSubmitting">
              {{ isSubmitting ? "Хадгалж байна..." : "Ажлын байр үүсгэх" }}
            </Button>
          </template>
        </JobForm>
      </CardContent>
    </Card>
  </div>
</template>
