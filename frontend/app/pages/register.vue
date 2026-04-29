<script setup lang="ts">
definePageMeta({
  layout: "auth",
});

type RegisterRole = "applicant" | "recruiter";

const selectedRole = ref<RegisterRole>("applicant");

const form = reactive({
  fullName: "",
  email: "",
  password: "",
  confirmPassword: "",
  companyName: "",
  recruiterPosition: "",
});

function submitRegister() {
  console.log("Register form submitted", {
    role: selectedRole.value,
    ...form,
  });
}
</script>

<template>
  <AuthCard
    title="Create your account"
    description="Set up an applicant or recruiter profile to start using the platform."
    footer-text="Already have an account?"
    footer-link-text="Login"
    footer-link-to="/login"
  >
    <form class="space-y-5" @submit.prevent="submitRegister">
      <div class="space-y-3">
        <Label>Choose your role</Label>
        <Tabs v-model="selectedRole" class="w-full">
          <TabsList class="grid w-full grid-cols-2">
            <TabsTrigger value="applicant">Applicant</TabsTrigger>
            <TabsTrigger value="recruiter">Recruiter</TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      <div class="space-y-2">
        <Label for="full-name">Full name</Label>
        <Input
          id="full-name"
          v-model="form.fullName"
          type="text"
          autocomplete="name"
          placeholder="Your full name"
        />
      </div>

      <div class="space-y-2">
        <Label for="register-email">Email</Label>
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
          <Label for="register-password">Password</Label>
          <Input
            id="register-password"
            v-model="form.password"
            type="password"
            autocomplete="new-password"
            placeholder="Create a password"
          />
        </div>

        <div class="space-y-2">
          <Label for="confirm-password">Confirm password</Label>
          <Input
            id="confirm-password"
            v-model="form.confirmPassword"
            type="password"
            autocomplete="new-password"
            placeholder="Confirm your password"
          />
        </div>
      </div>

      <div v-if="selectedRole === 'recruiter'" class="grid gap-5 sm:grid-cols-2">
        <div class="space-y-2">
          <Label for="company-name">Company name</Label>
          <Input
            id="company-name"
            v-model="form.companyName"
            type="text"
            autocomplete="organization"
            placeholder="Your company"
          />
        </div>

        <div class="space-y-2">
          <Label for="recruiter-position">Recruiter position</Label>
          <Input
            id="recruiter-position"
            v-model="form.recruiterPosition"
            type="text"
            autocomplete="organization-title"
            placeholder="Talent Acquisition Specialist"
          />
        </div>
      </div>

      <Button class="w-full" type="submit">Create account</Button>
    </form>
  </AuthCard>
</template>
