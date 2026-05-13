<script setup lang="ts">
import { toast } from "vue-sonner";
import { Plus, X } from "lucide-vue-next";

import type { Company } from "../composables/types";
import type { SaveCompanyPayload } from "../composables/types/payload";

const { user, me } = useAuth();
const authAPI = useAuthAPI();
const companyAPI = useCompanyAPI();
const router = useRouter();
const { cities, districtsFor } = useLocationOptions();

const company = ref<Company | null>(null);
const loadingCompany = ref(false);
const isSavingCompany = ref(false);
const companyError = ref("");

const profileForm = reactive({ firstName: "", lastName: "" });
const isSavingProfile = ref(false);

function fillProfileForm() {
  const parts = (user.value?.full_name || user.value?.name || "").split(" ");
  profileForm.firstName = parts[0] || "";
  profileForm.lastName = parts.slice(1).join(" ") || "";
}

async function saveProfile() {
  if (!profileForm.firstName.trim()) {
    toast.error("Нэр шаардлагатай");
    return;
  }
  isSavingProfile.value = true;
  try {
    await authAPI.updateMe({
      first_name: profileForm.firstName.trim(),
      last_name: profileForm.lastName.trim(),
    });
    await me();
    toast.success("Профайл амжилттай шинэчлэгдлээ.");
  } catch (e: any) {
    toast.error(e?.data?.message || "Профайл хадгалахад алдаа гарлаа");
  } finally {
    isSavingProfile.value = false;
  }
}

const companyForm = reactive({
  name: "",
  description: "",
  register_id: "",
  contact_info: "",
  city: "",
  district: "",
  location_x: "",
  location_y: "",
  benefits: [] as string[],
});

const isRecruiter = computed(() => user.value?.role === "recruiter");
const hasCompany = computed(() => Boolean(user.value?.company_id));

function displayName() {
  return user.value?.name || user.value?.full_name || "Your profile";
}

function fillCompanyForm(value: Company | null) {
  companyForm.name = value?.name || "";
  companyForm.description = value?.description || "";
  companyForm.register_id = value?.register_id || "";
  companyForm.contact_info = value?.contact_info || "";
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
  companyError.value = "";
  isSavingCompany.value = true;
  const wasHasCompany = hasCompany.value;

  try {
    const savedCompany = await companyAPI.save(buildCompanyPayload());
    company.value = savedCompany;
    fillCompanyForm(savedCompany);
    await me();
    toast.success(
      wasHasCompany
        ? "Компанийн профайл шинэчлэгдлээ."
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
    toast.warning("Ажлын байр нээхийн өмнө компанийн мэдээллээ нэмнэ үү.");
    return;
  }

  router.push("/jobs");
}

// ── Email change ─────────────────────────────────────────────────────────
const emailChangeStep = ref<"form" | "code">("form");
const newEmail = ref("");
const emailChangeCode = ref("");
const isSendingEmailCode = ref(false);
const isVerifyingEmail = ref(false);

async function sendEmailChangeCode() {
  if (!newEmail.value.trim()) {
    toast.error("Шинэ и-мэйл хаягаа оруулна уу");
    return;
  }
  isSendingEmailCode.value = true;
  try {
    const res = await authAPI.initiateEmailChange(newEmail.value.trim());
    toast.success(res.message);
    emailChangeStep.value = "code";
  } catch (e: any) {
    toast.error(e?.data?.message || "Код илгээхэд алдаа гарлаа");
  } finally {
    isSendingEmailCode.value = false;
  }
}

async function confirmEmailChange() {
  if (emailChangeCode.value.trim().length !== 6) {
    toast.error("6 оронтой кодыг оруулна уу");
    return;
  }
  isVerifyingEmail.value = true;
  try {
    await authAPI.verifyEmailChange(emailChangeCode.value.trim());
    await me();
    emailChangeStep.value = "form";
    newEmail.value = "";
    emailChangeCode.value = "";
    toast.success("И-мэйл хаяг амжилттай шинэчлэгдлээ.");
  } catch (e: any) {
    toast.error(e?.data?.message || "Баталгаажуулахад алдаа гарлаа");
  } finally {
    isVerifyingEmail.value = false;
  }
}

// ── Password change ───────────────────────────────────────────────────────
const passwordForm = reactive({ current: "", newPwd: "", confirm: "" });
const isChangingPassword = ref(false);

async function changePassword() {
  if (!passwordForm.current || !passwordForm.newPwd) {
    toast.error("Бүх талбарыг бөглөнө үү");
    return;
  }
  if (passwordForm.newPwd !== passwordForm.confirm) {
    toast.error("Нууц үг таарахгүй байна");
    return;
  }
  if (passwordForm.newPwd.length < 8) {
    toast.error("Нууц үг наад зах нь 8 тэмдэгттэй байна");
    return;
  }
  isChangingPassword.value = true;
  try {
    const res = await authAPI.changePassword(passwordForm.current, passwordForm.newPwd);
    toast.success(res.message);
    passwordForm.current = "";
    passwordForm.newPwd = "";
    passwordForm.confirm = "";
  } catch (e: any) {
    toast.error(e?.data?.message || "Нууц үг солиход алдаа гарлаа");
  } finally {
    isChangingPassword.value = false;
  }
}

await loadCompany();
fillProfileForm();
</script>

<template>
  <div class="space-y-6">
    <section
      class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm"
    >
      <p class="text-sm font-medium text-muted-foreground">Профайл</p>
      <h1 class="mt-2 text-3xl font-semibold tracking-tight">
        {{ displayName() }}
      </h1>
      <p class="mt-2 text-sm leading-6 text-muted-foreground">
        Бүртгэлийн мэдээлэл болон компанийн тохиргоог энд удирдана уу.
      </p>
    </section>

    <!-- Profile name edit -->
    <Card class="rounded-3xl border-border shadow-sm">
      <CardHeader>
        <CardTitle>Хувийн мэдээлэл засах</CardTitle>
        <CardDescription>Нэр, овгоо шинэчилнэ үү.</CardDescription>
      </CardHeader>
      <CardContent>
        <form class="grid gap-4 sm:grid-cols-2" @submit.prevent="saveProfile">
          <div class="space-y-2">
            <Label for="profile-first-name">Нэр</Label>
            <Input
              id="profile-first-name"
              v-model="profileForm.firstName"
              placeholder="Нэр"
              autocomplete="given-name"
            />
          </div>
          <div class="space-y-2">
            <Label for="profile-last-name">Овог</Label>
            <Input
              id="profile-last-name"
              v-model="profileForm.lastName"
              placeholder="Овог"
              autocomplete="family-name"
            />
          </div>
          <div class="sm:col-span-2 flex justify-end">
            <Button type="submit" :disabled="isSavingProfile">
              {{ isSavingProfile ? "Хадгалж байна..." : "Шинэчлэх" }}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>

    <!-- Email change card -->
    <Card class="rounded-3xl border-border shadow-sm">
      <CardHeader>
        <CardTitle>И-мэйл хаяг солих</CardTitle>
        <CardDescription>
          Шинэ и-мэйл хаяг руу баталгаажуулах код илгээгдэнэ.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="emailChangeStep === 'form'" class="grid gap-4 sm:grid-cols-[1fr_auto]">
          <div class="space-y-2">
            <Label for="new-email">Шинэ и-мэйл хаяг</Label>
            <Input
              id="new-email"
              v-model="newEmail"
              type="email"
              :placeholder="user?.email || 'шинэ@хаяг.com'"
              autocomplete="email"
            />
          </div>
          <div class="flex items-end">
            <Button :disabled="isSendingEmailCode" @click="sendEmailChangeCode">
              {{ isSendingEmailCode ? "Илгээж байна..." : "Код илгээх" }}
            </Button>
          </div>
        </div>

        <div v-else class="space-y-4">
          <p class="text-sm text-muted-foreground">
            Шинэ и-мэйл <span class="font-medium text-foreground">{{ newEmail }}</span> хаягт код илгээгдлээ.
          </p>
          <div class="grid gap-4 sm:grid-cols-[1fr_auto]">
            <div class="space-y-2">
              <Label for="email-code">Баталгаажуулах код</Label>
              <Input
                id="email-code"
                v-model="emailChangeCode"
                type="text"
                inputmode="numeric"
                maxlength="6"
                placeholder="000000"
                class="text-center text-xl tracking-[0.5em] font-bold"
              />
            </div>
            <div class="flex items-end gap-2">
              <Button variant="outline" @click="emailChangeStep = 'form'">Буцах</Button>
              <Button :disabled="isVerifyingEmail" @click="confirmEmailChange">
                {{ isVerifyingEmail ? "Баталгаажуулж байна..." : "Баталгаажуулах" }}
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Password change card -->
    <Card class="rounded-3xl border-border shadow-sm">
      <CardHeader>
        <CardTitle>Нууц үг солих</CardTitle>
        <CardDescription>Одоогийн нууц үгийг оруулж баталгаажуулна уу.</CardDescription>
      </CardHeader>
      <CardContent>
        <form class="grid gap-4 sm:grid-cols-2" @submit.prevent="changePassword">
          <div class="space-y-2 sm:col-span-2">
            <Label for="current-password">Одоогийн нууц үг</Label>
            <PasswordInput
              id="current-password"
              v-model="passwordForm.current"
              placeholder="Одоогийн нууц үг"
              autocomplete="current-password"
            />
          </div>
          <div class="space-y-2">
            <Label for="new-password">Шинэ нууц үг</Label>
            <PasswordInput
              id="new-password"
              v-model="passwordForm.newPwd"
              placeholder="Наад зах нь 8 тэмдэгт"
              autocomplete="new-password"
            />
          </div>
          <div class="space-y-2">
            <Label for="confirm-password-change">Нууц үг давтах</Label>
            <PasswordInput
              id="confirm-password-change"
              v-model="passwordForm.confirm"
              placeholder="Нууц үгийг дахин оруулна уу"
              autocomplete="new-password"
            />
          </div>
          <div class="sm:col-span-2 flex justify-end">
            <Button type="submit" :disabled="isChangingPassword">
              {{ isChangingPassword ? "Солиж байна..." : "Нууц үг солих" }}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>

    <section
      v-if="isRecruiter && !hasCompany"
      class="rounded-3xl border border-dashed border-border bg-muted/20 px-6 py-5"
    >
      <p class="text-sm font-semibold">Ажил олгогчийн тохиргоо дуусгах</p>
      <p class="mt-2 text-sm leading-6 text-muted-foreground">
        Таны бүртгэл бэлэн боллоо, гэхдээ ажлын зар нийтлэхийн өмнө компанийн мэдээллээ оруулна уу.
      </p>
    </section>

    <section class="grid gap-4 lg:grid-cols-[1.3fr_0.9fr]">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Бүртгэлийн дэлгэрэнгүй</CardTitle>
          <CardDescription>Нэвтэрсэн хэрэглэгчийн мэдээлэл.</CardDescription>
        </CardHeader>
        <CardContent class="grid gap-4 sm:grid-cols-2">
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Э-мэйл
            </p>
            <p class="mt-2 text-sm font-medium">{{ user?.email || "-" }}</p>
          </div>
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Үүрэг
            </p>
            <p class="mt-2 text-sm font-medium capitalize">{{ user?.role === "recruiter" ? "Ажил олгогч" : "Ажил горилогч" }}</p>
          </div>
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Компани
            </p>
            <p class="mt-2 text-sm font-medium">{{ company?.name || user?.company_name || "-" }}</p>
          </div>
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Албан тушаал
            </p>
            <p class="mt-2 text-sm font-medium">{{ user?.position || "-" }}</p>
          </div>
        </CardContent>
      </Card>

      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Дараагийн алхам</CardTitle>
          <CardDescription>Үүрэгт тохирсон хурдан зам.</CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <button
            type="button"
            class="block w-full rounded-2xl border border-border bg-muted/30 p-4 text-left text-sm font-medium transition hover:bg-muted/50"
            :class="isRecruiter && !hasCompany ? 'opacity-70' : ''"
            @click="openJobs"
          >
            {{ user?.role === "recruiter" ? "Ажлын байр удирдах" : "Ажлын байрууд харах" }}
          </button>
          <NuxtLink
            to="/settings"
            class="block rounded-2xl border border-border bg-muted/30 p-4 text-sm font-medium transition hover:bg-muted/50"
          >
            Бүртгэлийн тохиргоо харах
          </NuxtLink>
        </CardContent>
      </Card>
    </section>

    <Card v-if="isRecruiter" class="rounded-3xl border-border shadow-sm">
      <CardHeader>
        <CardTitle>{{ hasCompany ? "Компанийн профайл" : "Компани үүсгэх" }}</CardTitle>
        <CardDescription>
          {{ hasCompany ? "Ажил олгогчийн компанийн мэдээллийг шинэчлэн байгаарай." : "Ажлын байр нийтлэхийн өмнө энэ хэсгийг бөглөх шаардлагатай." }}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div
          v-if="companyError"
          class="mb-4 rounded-2xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
        >
          {{ companyError }}
        </div>

        <div v-if="loadingCompany" class="text-sm text-muted-foreground">
          Компанийн мэдээлэл ачааллаж байна...
        </div>

        <form v-else class="grid gap-4 sm:grid-cols-2" @submit.prevent="saveCompany">
          <div class="space-y-2 sm:col-span-2">
            <Label for="company-name">Компанийн нэр</Label>
            <Input id="company-name" v-model="companyForm.name" placeholder="Таны компани" />
          </div>

          <div class="space-y-2 sm:col-span-2">
            <Label for="company-description">Тайлбар</Label>
            <Textarea
              id="company-description"
              v-model="companyForm.description"
              rows="4"
              placeholder="Компанийхаа үйл ажиллагааны талаар бичнэ үү."
            />
          </div>

          <div class="space-y-2">
            <Label for="company-register-id">Регистрийн дугаар</Label>
            <Input id="company-register-id" v-model="companyForm.register_id" placeholder="A1234567" />
          </div>

          <div class="space-y-2">
            <Label for="company-contact-info">Холбоо барих мэдээлэл</Label>
            <Input id="company-contact-info" v-model="companyForm.contact_info" placeholder="hr@company.mn | +976 99000000" />
            <p class="text-xs text-muted-foreground">
              Ажлын зарт автоматаар нэмэгдэх холбоо барих мэдээлэл.
            </p>
          </div>

          <div class="space-y-2">
            <Label for="company-city">Хот / Аймаг</Label>
            <Select v-model="companyForm.city">
              <SelectTrigger id="company-city" class="w-full">
                <SelectValue placeholder="Хот эсвэл аймаг сонгох" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="city in cities" :key="city" :value="city">
                  {{ city }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="space-y-2">
            <Label for="company-district">Дүүрэг / Сум</Label>
            <Select v-model="companyForm.district" :disabled="!companyForm.city">
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

          <div class="sm:col-span-2 rounded-2xl border border-border bg-muted/20 p-4 space-y-2">
            <Label>Газрын зураг дээр байршил сонгох</Label>
            <LocationSearch
              v-model:model-x="companyForm.location_x"
              v-model:model-y="companyForm.location_y"
            />
          </div>

          <!-- Benefits / Incentives -->
          <div class="sm:col-span-2 space-y-3">
            <div class="flex items-center justify-between gap-3">
              <div>
                <Label>Урамшуулал / Давуу тал</Label>
                <p class="mt-1 text-xs text-muted-foreground">
                  Ажлын зарт нэмэхийн тулд компанийн давуу талуудыг жагсаана уу.
                </p>
              </div>
              <Button type="button" variant="outline" size="sm" @click="addBenefit">
                <Plus class="mr-1 h-4 w-4" />Нэмэх
              </Button>
            </div>
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
                <Button type="button" variant="outline" size="icon" @click="removeBenefit(index)">
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
          </div>

          <div class="sm:col-span-2 flex justify-end">
            <Button type="submit" :disabled="isSavingCompany">
              {{ isSavingCompany ? "Хадгалж байна..." : hasCompany ? "Шинэчлэх" : "Компани үүсгэх" }}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
