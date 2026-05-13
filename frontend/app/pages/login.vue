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
      errorMessage.value = data?.message || error?.message || "Нэвтрэхэд алдаа гарлаа. Дахин оролдоно уу.";
    }
  }
}
</script>

<template>
  <AuthShell
    title="Тавтай морилно уу"
    description="AI-д суурилсан ажилд авах үнэлгээний системд нэвтэрнэ үү."
    :highlights="[
      { label: 'Хандалт', value: 'Нэг нэвтрэлтээр хоёр үүрэгт' },
      { label: 'Шинжилгээ', value: 'AI ажилд авах урсгалаа үргэлжлүүлнэ' },
      { label: 'Аюулгүй байдал', value: 'Хамгаалагдсан бүртгэлийн хандалт' },
    ]"
  >
    <AuthCard
      title="Тавтай морилно уу"
      description="Нэр дэвшигчийн дэлгэрэнгүй, үнэлгээ болон анкетыг удирдах бүртгэлдээ нэвтэрнэ үү."
      footer-text="Бүртгэл үгүй юу?"
      footer-link-text="Бүртгүүлэх"
      footer-link-to="/register"
    >
      <form class="space-y-5" @submit.prevent="submitLogin">
        <div
          v-if="unverifiedEmail"
          class="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm leading-6 text-amber-800"
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
          class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm leading-6 text-red-700"
        >
          {{ errorMessage }}
        </div>

        <div class="space-y-2">
          <Label for="email">Э-мэйл</Label>
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
            <Label for="password">Нууц үг</Label>
            <NuxtLink
              to="/forgot-password"
              class="text-sm font-medium text-foreground transition hover:text-muted-foreground"
            >
              Нууц үг мартсан уу?
            </NuxtLink>
          </div>
          <PasswordInput
            id="password"
            v-model="form.password"
            autocomplete="current-password"
            placeholder="Нууц үгээ оруулна уу"
          />
        </div>

        <div class="rounded-2xl border border-border bg-muted/40 px-4 py-3">
          <p class="text-sm leading-6 text-muted-foreground">
            Ажил горилогч болон ажил олгогч хоёулаа энд нэвтэрч болно.
          </p>
        </div>

        <Button class="w-full" type="submit" :disabled="isLoading">
          {{ isLoading ? "Нэвтэрч байна..." : "Нэвтрэх" }}
        </Button>
      </form>
    </AuthCard>
  </AuthShell>
</template>
