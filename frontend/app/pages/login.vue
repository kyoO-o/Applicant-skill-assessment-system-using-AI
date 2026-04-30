<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: "guest",
});

const { login, isLoading } = useAuth();
const router = useRouter();

const form = reactive({
  email: "",
  password: "",
});

const errorMessage = ref("");

async function submitLogin() {
  errorMessage.value = "";

  try {
    await login({
      email: form.email,
      password: form.password,
    });

    await router.push("/");
  } catch (error: any) {
    errorMessage.value =
      error?.data?.message || error?.message || "Login failed. Please try again.";
  }
}
</script>

<template>
  <AuthShell
    title="Welcome back"
    description="Sign in to continue your AI-powered applicant assessment workflow."
    :highlights="[
      { label: 'Access', value: 'Single login for both roles' },
      { label: 'Insights', value: 'Continue your AI hiring flow' },
      { label: 'Security', value: 'Protected account access' },
    ]"
  >
    <AuthCard
      title="Welcome back"
      description="Use your account to manage candidate screening, assessments, and applications."
      footer-text="Don't have an account?"
      footer-link-text="Create one"
      footer-link-to="/register"
    >
      <form class="space-y-5" @submit.prevent="submitLogin">
        <div
          v-if="errorMessage"
          class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm leading-6 text-red-700"
        >
          {{ errorMessage }}
        </div>

        <div class="space-y-2">
          <Label for="email">Email</Label>
          <Input
            id="email"
            v-model="form.email"
            type="email"
            autocomplete="email"
            placeholder="name@company.com"
          />
        </div>

        <div class="space-y-2">
          <div class="flex items-center justify-between gap-3">
            <Label for="password">Password</Label>
            <NuxtLink
              to="/register"
              class="text-sm font-medium text-foreground transition hover:text-muted-foreground"
            >
              Need an account?
            </NuxtLink>
          </div>
          <Input
            id="password"
            v-model="form.password"
            type="password"
            autocomplete="current-password"
            placeholder="Enter your password"
          />
        </div>

        <div class="rounded-2xl border border-border bg-muted/40 px-4 py-3">
          <p class="text-sm leading-6 text-muted-foreground">
            Applicants and recruiters can both sign in here.
          </p>
        </div>

        <Button class="w-full" type="submit" :disabled="isLoading">
          {{ isLoading ? "Signing in..." : "Login" }}
        </Button>
      </form>
    </AuthCard>
  </AuthShell>
</template>
