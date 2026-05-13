<script setup lang="ts">
definePageMeta({ layout: false, middleware: "guest" });

const authAPI = useAuthAPI();
const router = useRouter();

type Step = "email" | "code" | "reset";
const step = ref<Step>("email");

const email = ref("");
const code = ref("");
const newPassword = ref("");
const confirmPassword = ref("");

const isSending = ref(false);
const isResetting = ref(false);
const errorMessage = ref("");
const successMessage = ref("");

async function sendCode() {
  errorMessage.value = "";
  if (!email.value.trim()) {
    errorMessage.value = "И-мэйл хаягаа оруулна уу";
    return;
  }
  isSending.value = true;
  try {
    await authAPI.forgotPassword(email.value.trim());
    step.value = "code";
  } catch (e: any) {
    errorMessage.value = e?.data?.message || "Алдаа гарлаа. Дахин оролдоно уу.";
  } finally {
    isSending.value = false;
  }
}

function proceedToReset() {
  errorMessage.value = "";
  if (code.value.trim().length !== 6) {
    errorMessage.value = "6 оронтой кодыг оруулна уу";
    return;
  }
  step.value = "reset";
}

async function resetPassword() {
  errorMessage.value = "";
  if (newPassword.value !== confirmPassword.value) {
    errorMessage.value = "Нууц үг таарахгүй байна";
    return;
  }
  if (newPassword.value.length < 8) {
    errorMessage.value = "Нууц үг наад зах нь 8 тэмдэгттэй байна";
    return;
  }
  isResetting.value = true;
  try {
    const res = await authAPI.resetPassword(email.value.trim(), code.value.trim(), newPassword.value);
    successMessage.value = res.message;
    setTimeout(() => router.push("/login"), 2000);
  } catch (e: any) {
    errorMessage.value = e?.data?.message || "Нууц үг сэргээхэд алдаа гарлаа";
  } finally {
    isResetting.value = false;
  }
}

async function resendCode() {
  isSending.value = true;
  try {
    await authAPI.forgotPassword(email.value.trim());
    errorMessage.value = "";
    successMessage.value = "Шинэ код илгээгдлээ.";
  } catch {
  } finally {
    isSending.value = false;
  }
}
</script>

<template>
  <AuthShell
    title="Нууц үг сэргээх"
    description="И-мэйл хаягаа оруулахад нууц үг сэргээх код илгээгдэнэ."
    :highlights="[
      { label: 'Аюулгүй', value: 'Зөвхөн и-мэйлээр баталгаажна' },
      { label: 'Хурдан', value: '15 минутын дотор' },
      { label: 'Хялбар', value: 'Шинэ нууц үг тохируулна' },
    ]"
  >
    <AuthCard
      title="Нууц үг сэргээх"
      :description="
        step === 'email' ? 'Бүртгэлтэй и-мэйл хаягаа оруулна уу.' :
        step === 'code'  ? 'И-мэйлд ирсэн 6 оронтой кодыг оруулна уу.' :
                           'Шинэ нууц үгээ тохируулна уу.'
      "
      footer-text="Нэвтрэх хуудас руу буцах?"
      footer-link-text="Нэвтрэх"
      footer-link-to="/login"
    >
      <!-- Step 1: Enter email -->
      <form v-if="step === 'email'" class="space-y-5" @submit.prevent="sendCode">
        <div
          v-if="errorMessage"
          class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
        >
          {{ errorMessage }}
        </div>

        <div class="space-y-2">
          <Label for="forgot-email">И-мэйл</Label>
          <Input
            id="forgot-email"
            v-model="email"
            type="email"
            placeholder="name@company.com"
            autocomplete="email"
          />
        </div>

        <Button class="w-full" type="submit" :disabled="isSending">
          {{ isSending ? "Илгээж байна..." : "Код илгээх" }}
        </Button>
      </form>

      <!-- Step 2: Enter code only -->
      <form v-else-if="step === 'code'" class="space-y-5" @submit.prevent="proceedToReset">
        <div
          v-if="errorMessage"
          class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
        >
          {{ errorMessage }}
        </div>

        <div
          class="rounded-2xl border border-border bg-muted/40 px-4 py-3 text-sm text-muted-foreground"
        >
          Код илгээгдсэн хаяг: <span class="font-medium text-foreground">{{ email }}</span>
        </div>

        <div class="space-y-2">
          <Label for="reset-code">Баталгаажуулах код</Label>
          <Input
            id="reset-code"
            v-model="code"
            type="text"
            inputmode="numeric"
            maxlength="6"
            placeholder="000000"
            class="text-center text-2xl tracking-[0.5em] font-bold"
          />
        </div>

        <Button class="w-full" type="submit">
          Дараагийн алхам →
        </Button>

        <div class="text-center text-sm text-muted-foreground">
          Код ирээгүй юу?
          <button
            type="button"
            class="ml-1 font-medium text-foreground hover:text-muted-foreground transition"
            :disabled="isSending"
            @click="resendCode"
          >
            {{ isSending ? "Илгээж байна..." : "Дахин илгээх" }}
          </button>
        </div>
      </form>

      <!-- Step 3: Enter new password -->
      <form v-else class="space-y-5" @submit.prevent="resetPassword">
        <div
          v-if="successMessage"
          class="rounded-2xl border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-700"
        >
          {{ successMessage }}
        </div>

        <div
          v-if="errorMessage"
          class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
        >
          {{ errorMessage }}
        </div>

        <div class="space-y-2">
          <Label for="new-password">Шинэ нууц үг</Label>
          <PasswordInput
            id="new-password"
            v-model="newPassword"
            placeholder="Наад зах нь 8 тэмдэгт"
            autocomplete="new-password"
          />
        </div>

        <div class="space-y-2">
          <Label for="confirm-new-password">Нууц үг давтах</Label>
          <PasswordInput
            id="confirm-new-password"
            v-model="confirmPassword"
            placeholder="Нууц үгийг дахин оруулна уу"
            autocomplete="new-password"
          />
        </div>

        <Button class="w-full" type="submit" :disabled="isResetting">
          {{ isResetting ? "Шинэчилж байна..." : "Нууц үг шинэчлэх" }}
        </Button>

        <button
          type="button"
          class="w-full text-center text-sm text-muted-foreground hover:text-foreground transition"
          @click="step = 'code'"
        >
          ← Буцах
        </button>
      </form>
    </AuthCard>
  </AuthShell>
</template>
