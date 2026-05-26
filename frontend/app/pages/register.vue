<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { UserRole } from "../composables/types/constants";
import { registerBaseSchema, registerRecruiterSchema } from "~/utils/schemas";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Check } from "lucide-vue-next";

definePageMeta({
  layout: false,
  middleware: "guest",
});

const selectedRole = ref<UserRole>(UserRole.User);
const { register, isLoading } = useAuth();
const router = useRouter();
const errorMessage = ref("");

const pwValue = ref("");
watch(selectedRole, () => {
  pwValue.value = "";
});

const passwordChecks = computed(() => [
  { ok: pwValue.value.length >= 8, label: "Багадаа 8 тэмдэгт" },
  { ok: /[A-Z]/.test(pwValue.value), label: "Том үсэг" },
  { ok: /[a-z]/.test(pwValue.value), label: "Жижиг үсэг" },
  { ok: /\d/.test(pwValue.value), label: "Тоо агуулсан" },
  { ok: /[^A-Za-z0-9]/.test(pwValue.value), label: "Тусгай тэмдэгт" },
]);
const passwordScore = computed(
  () => passwordChecks.value.filter((c) => c.ok).length,
);
const passwordStrength = computed(
  () =>
    ["Хэт сул", "Сул", "Дунд", "Дунд", "Хүчтэй", "Маш хүчтэй"][
      passwordScore.value
    ],
);
const passwordStrengthClass = computed(() =>
  passwordScore.value >= 4
    ? "text-emerald-600 dark:text-emerald-400"
    : passwordScore.value >= 3
      ? "text-yellow-600 dark:text-yellow-400"
      : "text-red-500",
);
const passwordBarClass = computed(() =>
  passwordScore.value >= 4
    ? "bg-emerald-500"
    : passwordScore.value >= 3
      ? "bg-yellow-500"
      : "bg-red-500",
);

const schema = computed(() =>
  toTypedSchema(
    selectedRole.value === UserRole.Recruiter
      ? registerRecruiterSchema
      : registerBaseSchema,
  ),
);

async function onSubmit(values: Record<string, any>) {
  errorMessage.value = "";
  try {
    const res = await register({
      first_name: values.firstName.trim(),
      last_name: values.lastName.trim(),
      company_name:
        selectedRole.value === UserRole.Recruiter
          ? values.companyName?.trim() || undefined
          : undefined,
      email: values.email,
      position: values.recruiterPosition?.trim() || undefined,
      password: values.password,
      role:
        selectedRole.value === UserRole.Recruiter
          ? UserRole.Recruiter
          : UserRole.User,
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
      description="Үүргээ сонгоод платформыг ашиглаж эхлэхийн тулд мэдээллийг бөглөнө үү."
      footer-text="Бүртгэл байна уу?"
      footer-link-text="Нэвтрэх"
      footer-link-to="/login"
    >
      <Form
        :key="selectedRole"
        :validation-schema="schema"
        class="space-y-5"
        @submit="onSubmit"
      >
        <div
          v-if="errorMessage"
          class="rounded-2xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm leading-6 text-destructive"
        >
          {{ errorMessage }}
        </div>

        <div class="space-y-3">
          <Label>Үүргээ сонгоно уу</Label>
          <Tabs v-model="selectedRole" class="w-full">
            <TabsList class="grid w-full grid-cols-2">
              <TabsTrigger value="user">Ажил горилогч</TabsTrigger>
              <TabsTrigger value="recruiter">Ажил олгогч</TabsTrigger>
            </TabsList>
          </Tabs>
        </div>

        <div class="grid gap-5 sm:grid-cols-2">
          <FormField
            v-slot="{ componentField }"
            name="firstName"
            :validate-on-blur="false"
            :validate-on-change="false"
            :validate-on-model-update="false"
          >
            <FormItem>
              <FormLabel class="text-xs text-muted-foreground">Нэр*</FormLabel>
              <FormControl>
                <Input
                  v-bind="componentField"
                  type="text"
                  autocomplete="given-name"
                  placeholder="Нэр"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <FormField
            v-slot="{ componentField }"
            name="lastName"
            :validate-on-blur="false"
            :validate-on-change="false"
            :validate-on-model-update="false"
          >
            <FormItem>
              <FormLabel class="text-xs text-muted-foreground">Овог*</FormLabel>
              <FormControl>
                <Input
                  v-bind="componentField"
                  type="text"
                  autocomplete="family-name"
                  placeholder="Овог"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
        </div>

        <FormField
          v-slot="{ componentField }"
          name="email"
          :validate-on-blur="false"
          :validate-on-change="false"
          :validate-on-model-update="false"
        >
          <FormItem>
            <FormLabel class="text-xs text-muted-foreground">Э-мэйл*</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="email"
                autocomplete="email"
                placeholder="name@company.com"
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
            <FormLabel class="text-xs text-muted-foreground"
              >Нууц үг*</FormLabel
            >
            <FormControl>
              <PasswordInput
                v-bind="componentField"
                autocomplete="new-password"
                placeholder="Нууц үгээ оруулна уу"
                @update:model-value="pwValue = $event as string"
              />
            </FormControl>
            <template v-if="pwValue">
              <div class="mt-2 flex gap-1">
                <div
                  v-for="i in 5"
                  :key="i"
                  class="h-1 flex-1 rounded-full transition-colors"
                  :class="i <= passwordScore ? passwordBarClass : 'bg-border'"
                />
              </div>
              <p
                class="mt-1 text-xs font-semibold"
                :class="passwordStrengthClass"
              >
                {{ passwordStrength }}
              </p>
              <div class="mt-2 grid grid-cols-2 gap-1">
                <div
                  v-for="check in passwordChecks"
                  :key="check.label"
                  class="flex items-center gap-1.5 text-xs"
                  :class="
                    check.ok
                      ? 'text-emerald-600 dark:text-emerald-400'
                      : 'text-muted-foreground'
                  "
                >
                  <Check v-if="check.ok" class="h-3 w-3" />
                  <span
                    v-else
                    class="inline-block h-2 w-2 rounded-full border border-muted-foreground/40"
                  />
                  {{ check.label }}
                </div>
              </div>
            </template>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField
          v-slot="{ componentField }"
          name="confirmPassword"
          :validate-on-blur="false"
          :validate-on-change="false"
          :validate-on-model-update="false"
        >
          <FormItem>
            <FormLabel class="text-xs text-muted-foreground"
              >Нууц үг давтах*</FormLabel
            >
            <FormControl>
              <PasswordInput
                v-bind="componentField"
                autocomplete="new-password"
                placeholder="Нууц үгээ дахин оруулна уу"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <div v-if="selectedRole === UserRole.Recruiter" class="space-y-5">
          <FormField
            v-slot="{ componentField }"
            name="recruiterPosition"
            :validate-on-blur="false"
            :validate-on-change="false"
            :validate-on-model-update="false"
          >
            <FormItem>
              <FormLabel class="text-xs text-muted-foreground"
                >Ажил олгогчийн албан тушаал*</FormLabel
              >
              <FormControl>
                <Input
                  v-bind="componentField"
                  type="text"
                  autocomplete="organization-title"
                  placeholder="Talent Acquisition Specialist"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- <FormField
            v-slot="{ componentField }"
            name="companyName"
            :validate-on-blur="false"
            :validate-on-change="false"
            :validate-on-model-update="false"
          >
            <FormItem>
              <FormLabel class="text-xs text-muted-foreground"
                >Компаний нэр</FormLabel
              >
              <FormControl>
                <Input
                  v-bind="componentField"
                  type="text"
                  autocomplete="organization"
                  placeholder="Таны компани"
                />
              </FormControl>
              <p class="text-xs text-muted-foreground">
                Ажил олгогчийн бүртгэл үүсгэсний дараа нэмж болно.
              </p>
              <FormMessage />
            </FormItem>
          </FormField> -->
        </div>

        <Button class="w-full mt-2" type="submit" :disabled="isLoading">
          {{ isLoading ? "Бүртгэл үүсгэж байна..." : "Бүртгүүлэх" }}
        </Button>
      </Form>
    </AuthCard>
  </AuthShell>
</template>
