<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: "guest",
});

const { login, isLoading } = useAuth();
const router = useRouter();

const form = reactive({ email: "", password: "" });
const errorMessage = ref("");
const unverifiedEmail = ref("");

async function submitLogin() {
  errorMessage.value = "";
  unverifiedEmail.value = "";

  try {
    await login({ email: form.email, password: form.password });
    await router.push("/");
  } catch (error: any) {
    const data = error?.data as any;
    if (data?.code === "EMAIL_NOT_VERIFIED") {
      unverifiedEmail.value = data.email || form.email;
    } else {
      errorMessage.value =
        data?.message || error?.message || "Нэвтрэхэд алдаа гарлаа. Дахин оролдоно уу.";
    }
  }
}
</script>

<template>
  <AuthShell>
    <AuthCard
      title="Тавтай морилно уу"
      description="Бүртгэлдээ нэвтэрч ажлын хайлтаа үргэлжлүүлнэ үү."
      footer-text="Бүртгэл үгүй юу?"
      footer-link-text="Бүртгүүлэх"
      footer-link-to="/register"
    >
      <form class="space-y-4" @submit.prevent="submitLogin">
        <div
          v-if="unverifiedEmail"
          class="rounded-xl px-4 py-3 text-[13px] leading-6"
          style="background: var(--warning-bg); color: var(--warning-foreground);"
        >
          И-мэйл хаяг баталгаажаагүй байна.
          <NuxtLink
            :to="`/verify-email?email=${encodeURIComponent(unverifiedEmail)}`"
            class="font-semibold underline"
          >
            Баталгаажуулах →
          </NuxtLink>
        </div>

        <div
          v-if="errorMessage"
          class="rounded-xl px-4 py-3 text-[13px] leading-6"
          style="background: oklch(0.626 0.228 28 / 8%); color: var(--destructive);"
        >
          {{ errorMessage }}
        </div>

        <div class="space-y-1.5">
          <Label for="email" class="text-[11.5px] font-medium text-muted-foreground">Э-мэйл</Label>
          <Input
            id="email"
            v-model="form.email"
            type="email"
            autocomplete="email"
            placeholder="name@company.com"
            class="rounded-xl"
          />
        </div>

        <div class="space-y-1.5">
          <div class="flex items-center justify-between gap-3">
            <Label for="password" class="text-[11.5px] font-medium text-muted-foreground">Нууц үг</Label>
            <NuxtLink
              to="/forgot-password"
              class="text-[11.5px] font-medium text-primary hover:text-primary/80"
            >
              Нууц үг мартсан уу?
            </NuxtLink>
          </div>
          <PasswordInput
            id="password"
            v-model="form.password"
            autocomplete="current-password"
            placeholder="Нууц үгээ оруулна уу"
            class="rounded-xl"
          />
        </div>

        <div class="rounded-xl border border-border bg-muted/40 px-4 py-3">
          <p class="text-[12.5px] leading-6 text-muted-foreground">
            Ажил горилогч болон ажил олгогч хоёулаа энд нэвтэрч болно.
          </p>
        </div>

        <button
          type="submit"
          :disabled="isLoading"
          class="flex w-full items-center justify-center gap-2 rounded-full py-3 text-[13.5px] font-semibold text-white transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60"
          style="background: linear-gradient(135deg, var(--primary), oklch(0.348 0.106 295))"
        >
          {{ isLoading ? "Нэвтэрч байна..." : "Нэвтрэх" }}
        </button>
      </form>
    </AuthCard>
  </AuthShell>
</template>
