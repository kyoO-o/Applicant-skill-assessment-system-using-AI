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
  firstName: "",
  lastName: "",
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
    errorMessage.value = "Нууц үг таарахгүй байна.";
    return;
  }

  try {
    const res = await register({
      first_name: form.firstName.trim(),
      last_name: form.lastName.trim(),
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

    await router.push(`/verify-email?email=${encodeURIComponent(res.email)}`);
  } catch (error: any) {
    errorMessage.value =
      error?.data?.message ||
      error?.message ||
      "Бүртгэл амжилтгүй болов. Дахин оролдоно уу.";
  }
}
</script>

<template>
  <AuthShell
    title="Бүртгэл үүсгэх"
    description="Ажил горилогч эсвэл ажил олгогч профайл тохируулж платформд нэвтэрнэ үү."
    :highlights="[
      { label: 'Бүртгэл', value: 'Хоёр үүрэгт хурдан бүртгэл' },
      { label: 'Профайл', value: 'Үүрэгт суурилсан бүртгэл үүсгэлт' },
      { label: 'Ажилд авах', value: 'Минутын дотор үнэлгээ эхлүүлэх' },
    ]"
  >
    <AuthCard
      title="Бүртгэл үүсгэх"
      description="Үүрэгээ сонгоод платформыг ашиглаж эхлэхийн тулд мэдээллийг бөглөнө үү."
      footer-text="Бүртгэл байна уу?"
      footer-link-text="Нэвтрэх"
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
          <Label>Үүрэгээ сонгоно уу</Label>
          <Tabs v-model="selectedRole" class="w-full">
            <TabsList class="grid w-full grid-cols-2">
              <TabsTrigger value="user">Ажил горилогч</TabsTrigger>
              <TabsTrigger value="recruiter">Ажил олгогч</TabsTrigger>
            </TabsList>
          </Tabs>
        </div>

        <div class="grid gap-5 sm:grid-cols-2">
          <div class="space-y-2">
            <Label for="first-name">Нэр*</Label>
            <Input
              id="first-name"
              v-model="form.firstName"
              type="text"
              autocomplete="given-name"
              placeholder="Нэр"
            />
          </div>
          <div class="space-y-2">
            <Label for="last-name">Овог*</Label>
            <Input
              id="last-name"
              v-model="form.lastName"
              type="text"
              autocomplete="family-name"
              placeholder="Овог"
            />
          </div>
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
            <PasswordInput
              id="register-password"
              v-model="form.password"
              autocomplete="new-password"
              placeholder="Нууц үгээ оруулна уу"
            />
          </div>

          <div class="space-y-2">
            <Label for="confirm-password">Нууц үг давтах*</Label>
            <PasswordInput
              id="confirm-password"
              v-model="form.confirmPassword"
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
              placeholder="Таны компани"
            />
            <p class="text-xs text-muted-foreground">
              Ажил олгогчийн бүртгэл үүсгэсний дараа нэмж болно.
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
          {{ isLoading ? "Бүртгэл үүсгэж байна..." : "Бүртгүүлэх" }}
        </Button>
      </form>
    </AuthCard>
  </AuthShell>
</template>
