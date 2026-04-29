<script setup lang="ts">
definePageMeta({
  layout: "auth",
});

const form = reactive({
  email: "",
  password: "",
});

function submitLogin() {
  console.log("Login form submitted", { ...form });
}
</script>

<template>
  <AuthCard
    title="Welcome back"
    description="Sign in to continue your AI-powered applicant assessment workflow."
    footer-text="Don't have an account?"
    footer-link-text="Create one"
    footer-link-to="/register"
  >
    <form class="space-y-5" @submit.prevent="submitLogin">
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
        <Label for="password">Password</Label>
        <Input
          id="password"
          v-model="form.password"
          type="password"
          autocomplete="current-password"
          placeholder="Enter your password"
        />
      </div>

      <div class="rounded-2xl border border-blue-100 bg-blue-50 px-4 py-3">
        <p class="text-sm leading-6 text-blue-900">
          Applicant or Recruiter can use the same login page.
        </p>
      </div>

      <Button class="w-full" type="submit">Login</Button>
    </form>
  </AuthCard>
</template>
