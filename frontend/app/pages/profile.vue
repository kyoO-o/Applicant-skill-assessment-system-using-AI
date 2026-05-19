<script setup lang="ts">
import { toast } from "vue-sonner";
import {
  Plus,
  X,
  Building2,
  User,
  Mail,
  Briefcase,
  Camera,
  Globe,
} from "lucide-vue-next";

import type { Company } from "../composables/types";
import type { SaveCompanyPayload } from "../composables/types/payload";
import { toTypedSchema } from "@vee-validate/zod";
import { userProfileSchema } from "~/utils/schemas";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";

const { user, me } = useAuth();
const authAPI = useAuthAPI();
const companyAPI = useCompanyAPI();
const router = useRouter();
const { cities, districtsFor } = useLocationOptions();
const config = useRuntimeConfig();

// ── Tab switcher (recruiter only) ─────────────────────────────
const profileTab = ref<"user" | "company">("user");

// ── Company state ──────────────────────────────────────────────
const company = ref<Company | null>(null);
const loadingCompany = ref(false);
const isSavingCompany = ref(false);
const companyError = ref("");

// ── Profile form ───────────────────────────────────────────────
const profileForm = reactive({ firstName: "", lastName: "" });
const isSavingProfile = ref(false);
const profileSchema = toTypedSchema(userProfileSchema);

// ── Avatar upload ───────────────────────────────────────────────
const avatarInput = ref<HTMLInputElement | null>(null);
const pendingAvatarFile = ref<File | null>(null);
const avatarPreview = ref<string | null>(null);

const avatarURL = computed(() => {
  if (avatarPreview.value) return avatarPreview.value;
  const url = user.value?.profile_url;
  if (!url) return null;
  if (url.startsWith("http")) return url;
  return `${config.public.apiBase}${url}`;
});

function onAvatarChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (!file) return;
  if (avatarPreview.value) URL.revokeObjectURL(avatarPreview.value);
  pendingAvatarFile.value = file;
  avatarPreview.value = URL.createObjectURL(file);
}

// ── Logo upload ────────────────────────────────────────────────
const logoInput = ref<HTMLInputElement | null>(null);
const pendingLogoFile = ref<File | null>(null);
const logoPreview = ref<string | null>(null);

const logoURL = computed(() => {
  if (logoPreview.value) return logoPreview.value;
  const url = company.value?.logo_url;
  if (!url) return null;
  if (url.startsWith("http")) return url;
  return `${config.public.apiBase}${url}`;
});

function onLogoChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (!file) return;
  if (logoPreview.value) URL.revokeObjectURL(logoPreview.value);
  pendingLogoFile.value = file;
  logoPreview.value = URL.createObjectURL(file);
}

// ── Company form ───────────────────────────────────────────────
const companyForm = reactive({
  name: "",
  description: "",
  register_id: "",
  contact_info: "",
  profile_url: "",
  city: "",
  district: "",
  location_x: "",
  location_y: "",
  benefits: [] as string[],
});

const isRecruiter = computed(() => user.value?.role === "recruiter");
const hasCompany = computed(() => Boolean(user.value?.company_id));

const initials = computed(() => {
  const name = user.value?.full_name || user.value?.name || "";
  return (
    name
      .split(" ")
      .map((s) => s[0])
      .slice(0, 2)
      .join("")
      .toUpperCase() || "?"
  );
});

const companyInitial = computed(() => {
  const name = company.value?.name || companyForm.name;
  return name ? name.charAt(0).toUpperCase() : "?";
});

function displayName() {
  return user.value?.name || user.value?.full_name || "Профайл";
}

function fillProfileForm() {
  const parts = (user.value?.full_name || user.value?.name || "").split(" ");
  profileForm.firstName = parts[0] || "";
  profileForm.lastName = parts.slice(1).join(" ") || "";
}

async function saveProfile(values: Record<string, any>) {
  isSavingProfile.value = true;
  try {
    if (pendingAvatarFile.value) {
      await authAPI.uploadAvatar(pendingAvatarFile.value);
      if (avatarPreview.value) URL.revokeObjectURL(avatarPreview.value);
      pendingAvatarFile.value = null;
      avatarPreview.value = null;
      if (avatarInput.value) avatarInput.value.value = "";
    }
    await authAPI.updateMe({
      first_name: values.firstName.trim(),
      last_name: values.lastName?.trim() ?? "",
    });
    await me();
    toast.success("Профайл амжилттай шинэчлэгдлээ.");
  } catch (e: any) {
    toast.error(e?.data?.message || "Профайл хадгалахад алдаа гарлаа");
  } finally {
    isSavingProfile.value = false;
  }
}

function fillCompanyForm(value: Company | null) {
  companyForm.name = value?.name || "";
  companyForm.description = value?.description || "";
  companyForm.register_id = value?.register_id || "";
  companyForm.contact_info = value?.contact_info || "";
  companyForm.profile_url = value?.profile_url || "";
  companyForm.city = value?.city || "";
  companyForm.district = value?.district || "";
  companyForm.location_x =
    typeof value?.location_x === "number" ? String(value.location_x) : "";
  companyForm.location_y =
    typeof value?.location_y === "number" ? String(value.location_y) : "";
  companyForm.benefits = value?.benefits?.map((b) => b.description) || [];
}

function addBenefit() {
  companyForm.benefits.push("");
}

function removeBenefit(index: number) {
  companyForm.benefits.splice(index, 1);
}

function buildCompanyPayload(): SaveCompanyPayload {
  return {
    name: companyForm.name.trim(),
    description: companyForm.description.trim(),
    register_id: companyForm.register_id.trim(),
    contact_info: companyForm.contact_info.trim(),
    profile_url: companyForm.profile_url.trim() || undefined,
    city: companyForm.city.trim(),
    district: companyForm.district.trim(),
    location_x: companyForm.location_x.trim()
      ? Number(companyForm.location_x.trim())
      : undefined,
    location_y: companyForm.location_y.trim()
      ? Number(companyForm.location_y.trim())
      : undefined,
    benefits: companyForm.benefits.map((b) => b.trim()).filter(Boolean),
  };
}

const districtOptions = computed(() => districtsFor(companyForm.city));

async function loadCompany() {
  if (!isRecruiter.value || !hasCompany.value) {
    company.value = null;
    fillCompanyForm(null);
    return;
  }
  loadingCompany.value = true;
  try {
    company.value = await companyAPI.get();
    fillCompanyForm(company.value);
  } catch {
    company.value = null;
  } finally {
    loadingCompany.value = false;
  }
}

async function saveCompany() {
  console.log("clicked");
  companyError.value = "";
  isSavingCompany.value = true;
  const wasHasCompany = hasCompany.value;
  try {
    const savedCompany = await companyAPI.save(buildCompanyPayload());
    company.value = savedCompany;
    fillCompanyForm(savedCompany);
    await me();

    if (pendingLogoFile.value) {
      const updated = await companyAPI.uploadLogo(pendingLogoFile.value);
      company.value = updated;
      if (logoPreview.value) URL.revokeObjectURL(logoPreview.value);
      pendingLogoFile.value = null;
      logoPreview.value = null;
      if (logoInput.value) logoInput.value.value = "";
    }

    toast.success(
      wasHasCompany
        ? "Компаний профайл шинэчлэгдлээ."
        : "Компани үүсгэгдлээ. Одоо ажлын байр удирдах боломжтой.",
    );
  } catch (error: any) {
    companyError.value =
      error?.data?.message || error?.message || "Unable to save the company.";
    toast.error(companyError.value);
  } finally {
    isSavingCompany.value = false;
  }
}

function openJobs() {
  if (isRecruiter.value && !hasCompany.value) {
    toast.warning("Ажлын байр нээхийн өмнө компаний мэдээллээ нэмнэ үү.");
    return;
  }
  router.push("/jobs");
}

await loadCompany();
fillProfileForm();
</script>

<template>
  <div>
    <!-- Page header -->
    <div class="mb-6 flex items-end justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Профайл</h1>
        <p class="mt-1.5 text-sm text-muted-foreground">
          <template v-if="isRecruiter && profileTab === 'company'">
            Ажил горилогчид компаний мэдээллийг харуулна.
          </template>
          <template v-else> Таны бүртгэл. </template>
        </p>
      </div>
    </div>

    <!-- Tab switcher (recruiter only) -->
    <div v-if="isRecruiter" class="mb-5">
      <div
        class="inline-flex gap-0.5 rounded-full border border-border bg-muted/50 p-1"
      >
        <button
          v-for="tab in [
            { value: 'user', label: 'Таны профайл', icon: User },
            { value: 'company', label: 'Компаний профайл', icon: Building2 },
          ]"
          :key="tab.value"
          class="flex items-center gap-2 rounded-full px-4 py-2 text-sm font-semibold transition-all duration-100"
          :class="
            profileTab === tab.value
              ? 'bg-background text-foreground shadow-sm ring-1 ring-border'
              : 'text-muted-foreground hover:text-foreground'
          "
          @click="profileTab = tab.value as 'user' | 'company'"
        >
          <component :is="tab.icon" class="h-3.5 w-3.5" />
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- ── User profile section ──────────────────────────────── -->
    <template v-if="!isRecruiter || profileTab === 'user'">
      <div class="grid gap-5 lg:grid-cols-[1fr_1.6fr]">
        <!-- Avatar + account info sidebar -->
        <Card class="rounded-3xl border-border shadow-none">
          <CardContent class="pt-6">
            <div class="flex flex-col items-center gap-3">
              <label class="relative group block h-24 w-24 cursor-pointer">
                <div
                  v-if="avatarURL"
                  class="h-24 w-24 rounded-full overflow-hidden border border-border"
                >
                  <img
                    :src="avatarURL"
                    alt="Avatar"
                    class="h-full w-full object-cover"
                  />
                </div>
                <div
                  v-else
                  class="flex h-24 w-24 items-center justify-center rounded-full bg-muted text-2xl font-bold text-muted-foreground"
                >
                  {{ initials }}
                </div>
                <div
                  class="absolute inset-0 flex items-center justify-center rounded-full bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  <Camera class="h-5 w-5 text-white" />
                </div>
                <input
                  ref="avatarInput"
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  class="hidden"
                  :disabled="isSavingProfile"
                  @change="onAvatarChange"
                />
              </label>
              <div class="text-center">
                <div class="text-lg font-semibold tracking-tight">
                  {{ displayName() }}
                </div>
                <div class="mt-1 text-sm text-muted-foreground">
                  {{
                    isRecruiter
                      ? user?.position || "Ажил олгогч"
                      : user?.position || "Ажил горилогч"
                  }}
                </div>
              </div>
            </div>

            <Separator class="my-5" />

            <div class="flex flex-col gap-3">
              <div class="flex items-center gap-3">
                <div
                  class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-muted"
                >
                  <Mail class="h-3.5 w-3.5 text-muted-foreground" />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-xs text-muted-foreground">Э-мэйл</p>
                  <p class="truncate text-sm font-medium">
                    {{ user?.email || "–" }}
                  </p>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <div
                  class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-muted"
                >
                  <User class="h-3.5 w-3.5 text-muted-foreground" />
                </div>
                <div class="flex-1">
                  <p class="text-xs text-muted-foreground">Үүрэг</p>
                  <p class="text-sm font-medium">
                    {{
                      user?.role === "recruiter"
                        ? "Ажил олгогч"
                        : "Ажил горилогч"
                    }}
                  </p>
                </div>
              </div>
              <div v-if="user?.company_name" class="flex items-center gap-3">
                <div
                  class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-muted"
                >
                  <Building2 class="h-3.5 w-3.5 text-muted-foreground" />
                </div>
                <div class="flex-1">
                  <p class="text-xs text-muted-foreground">Компани</p>
                  <p class="text-sm font-medium">{{ user.company_name }}</p>
                </div>
              </div>
              <div v-if="user?.position" class="flex items-center gap-3">
                <div
                  class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-muted"
                >
                  <Briefcase class="h-3.5 w-3.5 text-muted-foreground" />
                </div>
                <div class="flex-1">
                  <p class="text-xs text-muted-foreground">Албан тушаал</p>
                  <p class="text-sm font-medium">{{ user.position }}</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Right column -->
        <div class="flex flex-col gap-4 h-full">
          <!-- Personal info form -->
          <Card class="rounded-3xl border-border shadow-none h-full">
            <CardHeader>
              <CardTitle>Хувийн мэдээлэл</CardTitle>
              <CardDescription>Нэр, овгоо шинэчилнэ үү.</CardDescription>
            </CardHeader>
            <Form
              :validation-schema="profileSchema"
              :initial-values="profileForm"
              class="contents"
              @submit="saveProfile"
            >
              <CardContent class="h-full">
                <div class="grid gap-3 sm:grid-cols-2">
                  <FormField v-slot="{ componentField }" name="firstName">
                    <FormItem class="space-y-1.5">
                      <FormLabel>Нэр</FormLabel>
                      <FormControl>
                        <Input
                          v-bind="componentField"
                          placeholder="Нэр"
                          autocomplete="given-name"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  </FormField>

                  <FormField v-slot="{ componentField }" name="lastName">
                    <FormItem class="space-y-1.5">
                      <FormLabel>Овог</FormLabel>
                      <FormControl>
                        <Input
                          v-bind="componentField"
                          placeholder="Овог"
                          autocomplete="family-name"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  </FormField>

                  <div class="space-y-1.5">
                    <Label>Э-мэйл</Label>
                    <div
                      class="rounded-xl border border-border bg-muted/30 px-3 py-2.5 text-sm text-muted-foreground"
                    >
                      {{ user?.email }}
                    </div>
                    <p class="text-xs text-muted-foreground">
                      Э-мэйл солихыг
                      <NuxtLink
                        to="/settings"
                        class="text-primary underline-offset-4 hover:underline"
                        >Тохиргоо</NuxtLink
                      >-с хийнэ үү.
                    </p>
                  </div>
                  <div v-if="user?.phone_number" class="space-y-1.5">
                    <Label>Утас</Label>
                    <div
                      class="rounded-xl border border-border bg-muted/30 px-3 py-2.5 text-sm text-muted-foreground"
                    >
                      {{ user.phone_number }}
                    </div>
                  </div>
                </div>
              </CardContent>
              <CardFooter class="justify-end items-end">
                <Button type="submit" :disabled="isSavingProfile">
                  {{
                    isSavingProfile ? "Хадгалж байна..." : "Өөрчлөлт хадгалах"
                  }}
                </Button>
              </CardFooter>
            </Form>
          </Card>
        </div>
      </div>
    </template>

    <!-- ── Company profile section (recruiter) ──────────────── -->
    <template v-else-if="isRecruiter && profileTab === 'company'">
      <section
        v-if="!hasCompany"
        class="mb-5 rounded-3xl border border-dashed border-border bg-muted/20 px-6 py-5"
      >
        <p class="text-sm font-semibold">Ажил олгогчийн тохиргоо дуусгах</p>
        <p class="mt-2 text-sm leading-6 text-muted-foreground">
          Таны бүртгэл бэлэн боллоо, гэхдээ ажлын зар нийтлэхийн өмнө компанийн
          мэдээллээ оруулна уу.
        </p>
      </section>

      <div class="grid gap-5 lg:grid-cols-[1fr_1.6fr]">
        <!-- Company mark sidebar -->
        <Card class="rounded-3xl border-border shadow-none">
          <CardContent class="pt-6">
            <div class="flex flex-col items-center gap-3">
              <label
                class="relative group block h-24 w-24"
                :class="hasCompany ? 'cursor-pointer' : 'cursor-default'"
              >
                <div
                  v-if="logoURL"
                  class="h-24 w-24 rounded-full overflow-hidden border border-border"
                >
                  <img
                    :src="logoURL"
                    alt="Logo"
                    class="h-full w-full object-cover"
                  />
                </div>
                <div
                  v-else
                  class="flex h-24 w-24 items-center justify-center rounded-full bg-muted text-3xl font-bold text-muted-foreground"
                >
                  {{ companyInitial }}
                </div>
                <div
                  v-if="hasCompany"
                  class="absolute inset-0 flex items-center justify-center rounded-full bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  <Camera class="h-5 w-5 text-white" />
                </div>
                <input
                  ref="logoInput"
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  class="hidden"
                  :disabled="!hasCompany || isSavingCompany"
                  @change="onLogoChange"
                />
              </label>
              <div class="text-center">
                <div class="text-lg font-semibold tracking-tight">
                  {{ company?.name || companyForm.name || "Компани" }}
                </div>
                <div class="mt-1 text-sm text-muted-foreground">
                  {{
                    [companyForm.city, companyForm.district]
                      .filter(Boolean)
                      .join(" · ") || "Байршил тодорхойгүй"
                  }}
                </div>
              </div>
            </div>

            <Separator class="my-5" />

            <div class="flex flex-col gap-3">
              <div class="flex items-center gap-3">
                <div
                  class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-muted"
                >
                  <Briefcase class="h-3.5 w-3.5 text-muted-foreground" />
                </div>
                <div class="flex-1">
                  <p class="text-xs text-muted-foreground">Регистр</p>
                  <p class="text-sm font-medium">
                    {{ company?.register_id || companyForm.register_id || "–" }}
                  </p>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <div
                  class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-muted"
                >
                  <Mail class="h-3.5 w-3.5 text-muted-foreground" />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-xs text-muted-foreground">Холбоо барих</p>
                  <p class="truncate text-sm font-medium">
                    {{
                      company?.contact_info || companyForm.contact_info || "–"
                    }}
                  </p>
                </div>
              </div>
              <div
                v-if="company?.profile_url || companyForm.profile_url"
                class="flex items-center gap-3"
              >
                <div
                  class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-muted"
                >
                  <Globe class="h-3.5 w-3.5 text-muted-foreground" />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-xs text-muted-foreground">Вэбсайт</p>
                  <a
                    :href="company?.profile_url || companyForm.profile_url"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="truncate text-sm font-medium text-primary hover:underline underline-offset-2"
                  >
                    {{ company?.profile_url || companyForm.profile_url }}
                  </a>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Company form -->
        <div class="flex flex-col gap-4">
          <div
            v-if="companyError"
            class="rounded-2xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
          >
            {{ companyError }}
          </div>

          <div v-if="loadingCompany" class="text-sm text-muted-foreground">
            Компанийн мэдээлэл ачааллаж байна...
          </div>

          <template v-else>
            <!-- About company -->
            <Card class="rounded-3xl border-border shadow-none">
              <CardHeader>
                <CardTitle>Компанийн тухай</CardTitle>
              </CardHeader>
              <CardContent class="grid gap-3 sm:grid-cols-2">
                <div class="space-y-1.5 sm:col-span-2">
                  <Label for="company-name">Компанийн нэр</Label>
                  <Input
                    id="company-name"
                    v-model="companyForm.name"
                    placeholder="Таны компани"
                  />
                </div>
                <div class="space-y-1.5">
                  <Label for="company-register-id">Регистрийн дугаар</Label>
                  <Input
                    id="company-register-id"
                    v-model="companyForm.register_id"
                    placeholder="A1234567"
                  />
                </div>
                <div class="space-y-1.5">
                  <Label for="company-contact-info">Холбоо барих</Label>
                  <Input
                    id="company-contact-info"
                    v-model="companyForm.contact_info"
                    placeholder="hr@company.mn | +976 99000000"
                  />
                </div>
                <div class="space-y-1.5 sm:col-span-2">
                  <Label for="company-profile-url"
                    >Вэбсайт / Профайл холбоос</Label
                  >
                  <Input
                    id="company-profile-url"
                    v-model="companyForm.profile_url"
                    type="url"
                    placeholder="https://company.mn"
                  />
                </div>
                <div class="space-y-1.5 sm:col-span-2">
                  <Label for="company-description">Тайлбар</Label>
                  <Textarea
                    id="company-description"
                    v-model="companyForm.description"
                    :rows="3"
                    placeholder="Компанийхаа үйл ажиллагааны талаар бичнэ үү."
                  />
                </div>
              </CardContent>
            </Card>

            <!-- Office location -->
            <Card
              class="overflow-hidden rounded-3xl border-border p-0 shadow-none"
            >
              <div
                class="flex items-center justify-between border-b border-border px-6 py-4"
              >
                <div>
                  <p class="text-sm font-semibold">Оффисын байршил</p>
                  <p class="mt-0.5 text-xs text-muted-foreground">
                    Хот, дүүрэг болон газрын зурган дээрх байршлыг тохируулна
                    уу.
                  </p>
                </div>
              </div>
              <div class="grid gap-3 p-6 sm:grid-cols-2">
                <div class="space-y-1.5">
                  <Label for="company-city">Хот / Аймаг</Label>
                  <Select v-model="companyForm.city">
                    <SelectTrigger id="company-city" class="w-full">
                      <SelectValue placeholder="Хот эсвэл аймаг сонгох" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="city in cities"
                        :key="city"
                        :value="city"
                      >
                        {{ city }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div class="space-y-1.5">
                  <Label for="company-district">Дүүрэг / Сум</Label>
                  <Select
                    v-model="companyForm.district"
                    :disabled="!companyForm.city"
                  >
                    <SelectTrigger id="company-district" class="w-full">
                      <SelectValue placeholder="Дүүрэг эсвэл сум сонгох" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="district in districtOptions"
                        :key="district"
                        :value="district"
                      >
                        {{ district }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div
                  class="space-y-2 rounded-2xl border border-border bg-muted/20 p-4 sm:col-span-2"
                >
                  <Label>Газрын зураг дээр байршил сонгох</Label>
                  <LocationSearch
                    v-model:model-x="companyForm.location_x"
                    v-model:model-y="companyForm.location_y"
                  />
                </div>
              </div>
            </Card>

            <!-- Benefits -->
            <Card class="rounded-3xl border-border shadow-none">
              <CardHeader>
                <div class="flex items-center justify-between">
                  <div>
                    <CardTitle>Урамшуулал / Давуу тал</CardTitle>
                    <CardDescription class="mt-0.5">
                      Ажлын зарт нэмэхийн тулд компаний давуу талуудыг жагсаана
                      уу.
                    </CardDescription>
                  </div>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    @click="addBenefit"
                  >
                    <Plus class="mr-1 h-4 w-4" />Нэмэх
                  </Button>
                </div>
              </CardHeader>
              <CardContent class="space-y-2">
                <div v-if="companyForm.benefits.length" class="space-y-2">
                  <div
                    v-for="(_, index) in companyForm.benefits"
                    :key="`benefit-${index}`"
                    class="flex items-center gap-2"
                  >
                    <Input
                      v-model="companyForm.benefits[index]"
                      placeholder="Жишээ: Эрүүл мэндийн даатгал, Уян хатан цаг..."
                    />
                    <Button
                      type="button"
                      variant="outline"
                      size="icon"
                      @click="removeBenefit(index)"
                    >
                      <X class="h-4 w-4" />
                    </Button>
                  </div>
                </div>
                <div
                  v-else
                  class="rounded-2xl border border-dashed border-border px-4 py-3 text-sm text-muted-foreground"
                >
                  Одоогоор давуу тал нэмэгдээгүй байна.
                </div>
              </CardContent>
            </Card>

            <div class="flex justify-end">
              <Button :disabled="isSavingCompany" @click="saveCompany">
                {{
                  isSavingCompany
                    ? "Хадгалж байна..."
                    : hasCompany
                      ? "Өөрчлөлт хадгалах"
                      : "Компани үүсгэх"
                }}
              </Button>
            </div>
          </template>
        </div>
      </div>
    </template>
  </div>
</template>
