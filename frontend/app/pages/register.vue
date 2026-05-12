<script setup lang="ts">
import { UserRole } from "../composables/types/constants";

definePageMeta({
  layout: false,
  middleware: "guest",
});

const selectedRole = ref<UserRole>(UserRole.User);
const { register, isLoading } = useAuth();
const router = useRouter();

const form = reactive({
  fullName: "",
  email: "",
  password: "",
  confirmPassword: "",
  companyName: "",
  recruiterPosition: "",
});

const errorMessage = ref("");

async function submitRegister() {
  errorMessage.value = "";

  if (form.password !== form.confirmPassword) {
    errorMessage.value = "Passwords do not match.";
    return;
  }

  try {
    await register({
      name: form.fullName,
      company_name:
        selectedRole.value === "recruiter"
          ? form.companyName.trim()
          : undefined,
      email: form.email,
      position: form.recruiterPosition.trim() || undefined,
      password: form.password,
      role:
        selectedRole.value === "recruiter" ? UserRole.Recruiter : UserRole.User,
    });

    await router.push("/");
  } catch (error: any) {
    errorMessage.value =
      error?.data?.message ||
      error?.message ||
      "Registration failed. Please try again.";
  }
}
</script>

<template>
  <AuthShell
    title="Create your account"
    description="Set up an applicant or recruiter profile and get into the platform without using the main app layout."
    :highlights="[
      { label: 'Onboarding', value: 'Fast setup for both roles' },
      { label: 'Profile', value: 'Role-based account creation' },
      { label: 'Hiring', value: 'Start assessments in minutes' },
    ]"
  >
    <AuthCard
      title="Create your account"
      description="Choose your role and fill in the details to begin using the platform."
      footer-text="Already have an account?"
      footer-link-text="Login"
      footer-link-to="/login"
    >
      <form class="space-y-5" @submit.prevent="submitRegister">
        <div
          v-if="errorMessage"
          class="rounded-2xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm leading-6 text-destructive"
        >
          {{ errorMessage }}
        </div>

        <div class="space-y-3">
          <Label>Choose your role</Label>
          <Tabs v-model="selectedRole" class="w-full">
            <TabsList class="grid w-full grid-cols-2">
              <TabsTrigger value="user">Ажил горилогч</TabsTrigger>
              <TabsTrigger value="recruiter">Ажил олгогч</TabsTrigger>
            </TabsList>
          </Tabs>
        </div>

        <div class="space-y-2">
          <Label for="full-name">Овог нэр*</Label>
          <Input
            id="full-name"
            v-model="form.fullName"
            type="text"
            autocomplete="name"
            placeholder="Your full name"
          />
        </div>

        <div class="space-y-2">
          <Label for="register-email">Э-мэйл*</Label>
          <Input
            id="register-email"
            v-model="form.email"
            type="email"
            autocomplete="email"
            placeholder="name@company.com"
          />
        </div>

        <div class="grid gap-5 sm:grid-cols-2">
          <div class="space-y-2">
            <Label for="register-password">Нууц үг*</Label>
            <Input
              id="register-password"
              v-model="form.password"
              type="password"
              autocomplete="new-password"
              placeholder="Нууц үгээ оруулна уу"
            />
          </div>

          <div class="space-y-2">
            <Label for="confirm-password"></Label>
            <Input
              id="confirm-password"
              v-model="form.confirmPassword"
              type="password"
              autocomplete="new-password"
              placeholder="Нууц үгээ дахин оруулна уу"
            />
          </div>
        </div>

        <div
          v-if="selectedRole === 'recruiter'"
          class="grid gap-5 sm:grid-cols-2"
        >
          <div class="space-y-2">
            <Label for="company-name">Компаний нэр</Label>
            <Input
              id="company-name"
              v-model="form.companyName"
              type="text"
              autocomplete="organization"
              placeholder="Your company"
            />
            <p class="text-xs text-muted-foreground">
              You can also add this after creating your recruiter account.
            </p>
          </div>

          <div class="space-y-2">
            <Label for="recruiter-position">Ажил олгогчийн албан тушаал*</Label>
            <Input
              id="recruiter-position"
              v-model="form.recruiterPosition"
              type="text"
              autocomplete="organization-title"
              placeholder="Talent Acquisition Specialist"
            />
          </div>
        </div>

        <Button class="w-full" type="submit" :disabled="isLoading">
          {{ isLoading ? "Creating account..." : "Create account" }}
        </Button>
      </form>
    </AuthCard>
  </AuthShell>
</template>
