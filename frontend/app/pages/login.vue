<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { loginSchema } from "~/utils/schemas";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";

definePageMeta({
  layout: false,
  middleware: "guest",
});

const { login, isLoading } = useAuth();
const router = useRouter();

const unverifiedEmail = ref("");
const errorMessage = ref("");

const schema = toTypedSchema(loginSchema);

async function onSubmit(values: { email: string; password: string }) {
  errorMessage.value = "";
  unverifiedEmail.value = "";

  try {
    await login({ email: values.email, password: values.password });
    await router.push("/");
  } catch (error: any) {
    const data = error?.data as any;
    if (data?.code === "EMAIL_NOT_VERIFIED") {
      unverifiedEmail.value = data.email || values.email;
    } else {
      errorMessage.value =
        data?.message ||
        error?.message ||
        "Нэвтрэхэд алдаа гарлаа. Дахин оролдоно уу.";
    }
  }
}
</script>

<template>
  <AuthShell>
    <AuthCard
      title="Тавтай морилно уу"
      description="Бүртгэлээрээ нэвтэрч ажлын хайлтаа үргэлжлүүлнэ үү."
      footer-text="Бүртгэлгүй юу?"
      footer-link-text="Бүртгүүлэх"
      footer-link-to="/register"
    >
      <Form :validation-schema="schema" class="space-y-4" @submit="onSubmit">
        <div
          v-if="unverifiedEmail"
          class="rounded-xl px-4 py-3 text-[13px] leading-6"
          style="
            background: var(--warning-bg);
            color: var(--warning-foreground);
          "
        >
          Э-мэйл хаяг баталгаажаагүй байна.
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
          style="
            background: oklch(0.626 0.228 28 / 8%);
            color: var(--destructive);
          "
        >
          {{ errorMessage }}
        </div>

        <FormField
          v-slot="{ componentField }"
          name="email"
          :validate-on-blur="false"
          :validate-on-change="false"
          :validate-on-model-update="false"
        >
          <FormItem>
            <FormLabel class="text-xs font-medium text-muted-foreground"
              >Э-мэйл</FormLabel
            >
            <FormControl>
              <Input
                v-bind="componentField"
                type="email"
                autocomplete="email"
                placeholder="name@company.com"
                class="rounded-xl"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField
          v-slot="{ componentField }"
          name="password"
          :validate-on-blur="false"
          :validate-on-change="false"
          :validate-on-model-update="false"
        >
          <FormItem>
            <div class="flex items-center justify-between gap-3">
              <FormLabel class="text-xs font-medium text-muted-foreground"
                >Нууц үг</FormLabel
              >
              <NuxtLink
                to="/forgot-password"
                class="text-xs font-medium text-primary hover:text-primary/80"
              >
                Нууц үг мартсан уу?
              </NuxtLink>
            </div>
            <FormControl>
              <PasswordInput
                v-bind="componentField"
                autocomplete="current-password"
                placeholder="Нууц үгээ оруулна уу"
                class="rounded-xl"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <!-- <div class="rounded-xl border border-border bg-muted/40 px-4 py-3">
          <p class="text-[12.5px] leading-6 text-muted-foreground">
            Ажил горилогч болон ажил олгогч хоёулаа энд нэвтэрч болно.
          </p>
        </div> -->

        <button
          type="submit"
          :disabled="isLoading"
          class="flex w-full items-center justify-center gap-2 rounded-full py-3 text-[13.5px] font-semibold text-white transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60"
          style="
            background: linear-gradient(
              135deg,
              var(--primary),
              oklch(0.348 0.106 295)
            );
          "
        >
          {{ isLoading ? "Нэвтэрч байна..." : "Нэвтрэх" }}
        </button>
      </Form>
    </AuthCard>
  </AuthShell>
</template>
