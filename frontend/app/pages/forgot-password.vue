<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import {
  forgotPasswordEmailSchema,
  forgotPasswordCodeSchema,
  resetPasswordSchema,
} from "~/utils/schemas";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";

definePageMeta({ layout: false, middleware: "guest" });

const authAPI = useAuthAPI();
const router = useRouter();

type Step = "email" | "code" | "reset";
const step = ref<Step>("email");

const email = ref("");
const code = ref("");

const isSending = ref(false);
const isResetting = ref(false);
const errorMessage = ref("");
const successMessage = ref("");

const emailSchema = toTypedSchema(forgotPasswordEmailSchema);
const codeSchema = toTypedSchema(forgotPasswordCodeSchema);
const pwSchema = toTypedSchema(resetPasswordSchema);

async function sendCode(values: { email: string }) {
  errorMessage.value = "";
  email.value = values.email;
  isSending.value = true;
  try {
    await authAPI.forgotPassword(values.email);
    step.value = "code";
  } catch (e: any) {
    errorMessage.value = e?.data?.message || "Алдаа гарлаа. Дахин оролдоно уу.";
  } finally {
    isSending.value = false;
  }
}

function proceedToReset(values: { code: string }) {
  errorMessage.value = "";
  code.value = values.code;
  step.value = "reset";
}

async function resetPassword(values: {
  newPassword: string;
  confirmPassword: string;
}) {
  errorMessage.value = "";
  isResetting.value = true;
  try {
    const res = await authAPI.resetPassword(
      email.value,
      code.value,
      values.newPassword,
    );
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
    await authAPI.forgotPassword(email.value);
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
    description="Э-мэйл хаягаа оруулахад нууц үг сэргээх код илгээгдэнэ."
    :highlights="[
      { label: 'Аюулгүй', value: 'Зөвхөн и-мэйлээр баталгаажна' },
      { label: 'Хурдан', value: '15 минутын дотор' },
      { label: 'Хялбар', value: 'Шинэ нууц үг тохируулна' },
    ]"
  >
    <AuthCard
      title="Нууц үг сэргээх"
      :description="
        step === 'email'
          ? 'Бүртгэлтэй и-мэйл хаягаа оруулна уу.'
          : step === 'code'
            ? 'Э-мэйлд ирсэн 6 оронтой кодыг оруулна уу.'
            : 'Шинэ нууц үгээ тохируулна уу.'
      "
      footer-text="Нэвтрэх хуудас руу буцах?"
      footer-link-text="Нэвтрэх"
      footer-link-to="/login"
    >
      <!-- Step 1: Enter email -->
      <Form
        v-if="step === 'email'"
        :validation-schema="emailSchema"
        class="space-y-5"
        @submit="sendCode"
      >
        <div
          v-if="errorMessage"
          class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
        >
          {{ errorMessage }}
        </div>

        <FormField v-slot="{ componentField }" name="email">
          <FormItem class="space-y-2">
            <FormLabel>Э-мэйл</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="email"
                placeholder="name@company.com"
                autocomplete="email"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <Button class="w-full" type="submit" :disabled="isSending">
          {{ isSending ? "Илгээж байна..." : "Код илгээх" }}
        </Button>
      </Form>

      <!-- Step 2: Enter code -->
      <Form
        v-else-if="step === 'code'"
        :validation-schema="codeSchema"
        class="space-y-5"
        @submit="proceedToReset"
      >
        <div
          v-if="errorMessage"
          class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
        >
          {{ errorMessage }}
        </div>

        <div
          class="rounded-2xl border border-border bg-muted/40 px-4 py-3 text-sm text-muted-foreground"
        >
          Код илгээгдсэн хаяг:
          <span class="font-medium text-foreground">{{ email }}</span>
        </div>

        <FormField v-slot="{ componentField }" name="code">
          <FormItem class="space-y-2">
            <FormLabel>Баталгаажуулах код</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="text"
                inputmode="numeric"
                maxlength="6"
                placeholder="000000"
                class="text-center text-2xl tracking-[0.5em] font-bold"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <Button class="w-full" type="submit">Дараагийн алхам →</Button>

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
      </Form>

      <!-- Step 3: Enter new password -->
      <Form
        v-else
        :validation-schema="pwSchema"
        class="space-y-5"
        @submit="resetPassword"
      >
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

        <FormField v-slot="{ componentField }" name="newPassword">
          <FormItem class="space-y-2">
            <FormLabel>Шинэ нууц үг</FormLabel>
            <FormControl>
              <PasswordInput
                v-bind="componentField"
                placeholder="Багадаа 8 тэмдэгт"
                autocomplete="new-password"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="confirmPassword">
          <FormItem class="space-y-2">
            <FormLabel>Нууц үг давтах</FormLabel>
            <FormControl>
              <PasswordInput
                v-bind="componentField"
                placeholder="Нууц үгийг дахин оруулна уу"
                autocomplete="new-password"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

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
      </Form>
    </AuthCard>
  </AuthShell>
</template>
