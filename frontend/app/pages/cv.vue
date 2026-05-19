<script setup lang="ts">
import { toast } from "vue-sonner";
import * as z from "zod";
import {
  Loader2,
  Plus,
  Trash2,
  Upload,
  Printer,
  Save,
  ScrollText,
  Briefcase,
  GraduationCap,
  Brain,
  BookOpen,
  FlaskConical,
  Building2,
  Trophy,
} from "lucide-vue-next";
import type {
  CVProfile,
  WorkExperience,
  Education,
  Language,
  Training,
  Exam,
  Internship,
  Award,
} from "../composables/types";
import {
  emptyCV,
  emptyWorkExperience,
  emptyEducation,
  emptyLanguage,
  emptyTraining,
  emptyExam,
  emptyInternship,
  emptyAward,
} from "../composables/types";

definePageMeta({ middleware: "auth" });

const cvAPI = useCVAPI();
const { user } = useAuth();

const loading = ref(true);
const saving = ref(false);
const parsing = ref(false);

const cv = ref<CVProfile>(emptyCV());

const cvRequiredSchema = z.object({
  last_name: z.string().min(1, "Овог шаардлагатай"),
  first_name: z.string().min(1, "Нэр шаардлагатай"),
});

const DRIVER_LICENSES = ["A", "B", "C", "D", "E", "M"];
const GENDERS = ["Эрэгтэй", "Эмэгтэй"];
const MARITAL_STATUSES = ["Гэрлэсэн", "Гэрлээгүй", "Салсан", "Бэлэвсэн"];
const LANGUAGE_LEVELS = [
  "Анхан шатны",
  "Дунд",
  "Дэвшилтэт",
  "Чөлөөт",
  "Эх хэл",
];

async function load() {
  try {
    const data = await cvAPI.get();
    cv.value = { ...emptyCV(), ...data };
  } catch {
    toast.error("CV мэдээлэл ачааллаж чадсангүй");
  } finally {
    loading.value = false;
  }
}

await load();

async function save() {
  const validation = cvRequiredSchema.safeParse(cv.value);
  if (!validation.success) {
    toast.error(validation.error.issues[0]?.message ?? "Мэдээлэл дутуу байна");
    return;
  }
  saving.value = true;
  try {
    const saved = await cvAPI.save(cv.value);
    cv.value = { ...emptyCV(), ...saved };
    toast.success("CV хадгалагдлаа");
  } catch {
    toast.error("Хадгалахад алдаа гарлаа");
  } finally {
    saving.value = false;
  }
}

const fileInput = ref<HTMLInputElement | null>(null);

async function onImportPDF(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  parsing.value = true;
  try {
    const parsed = await cvAPI.parse(file);
    cv.value = { ...cv.value, ...parsed };
    toast.success("CV амжилттай уншлаа. Мэдээллийг шалгаад хадгалаарай.");
  } catch {
    toast.error("PDF уншихад алдаа гарлаа");
  } finally {
    parsing.value = false;
    input.value = "";
  }
}

function print() {
  window.print();
}

// Driver licenses
function toggleLicense(lic: string) {
  const idx = cv.value.driver_licenses.indexOf(lic);
  if (idx === -1) cv.value.driver_licenses.push(lic);
  else cv.value.driver_licenses.splice(idx, 1);
}

// Generic list item operations
function addItem<T>(list: T[], factory: () => T) {
  list.push(factory());
}
function removeItem<T>(list: T[], index: number) {
  list.splice(index, 1);
}

// Skill tag helpers
const skillInput = ref<Record<string, string>>({
  personal: "",
  professional: "",
  computer: "",
  art: "",
  sport: "",
});

function addSkill(field: keyof CVProfile, inputKey: string) {
  const val = (skillInput.value[inputKey] ?? "").trim();
  if (!val) return;
  const arr = cv.value[field] as string[];
  if (!arr.includes(val)) arr.push(val);
  skillInput.value[inputKey] = "";
}

function removeSkill(field: keyof CVProfile, index: number) {
  const arr = cv.value[field] as string[];
  arr.splice(index, 1);
}

function onSkillKeydown(
  e: KeyboardEvent,
  field: keyof CVProfile,
  inputKey: string,
) {
  if (e.key === "Enter") {
    e.preventDefault();
    addSkill(field, inputKey);
  }
}

// Sections for sidebar nav
const sections = [
  { id: "general", label: "Ерөнхий мэдээлэл", icon: ScrollText },
  { id: "about", label: "Миний тухай", icon: BookOpen },
  { id: "work", label: "Ажлын туршлага", icon: Briefcase },
  { id: "education", label: "Боловсрол", icon: GraduationCap },
  { id: "skills", label: "Ур чадвар", icon: Brain },
  { id: "trainings", label: "Сургалт, сертификат", icon: BookOpen },
  { id: "exams", label: "Шалгалт, тест", icon: FlaskConical },
  { id: "internships", label: "Дадлага", icon: Building2 },
  { id: "awards", label: "Шагнал урамшуулал", icon: Trophy },
];

function scrollTo(id: string) {
  document
    .getElementById(id)
    ?.scrollIntoView({ behavior: "smooth", block: "start" });
}

const fullName = computed(() =>
  [cv.value.last_name, cv.value.first_name].filter(Boolean).join(" "),
);

const currentJobTitle = computed(
  () => cv.value.work_experiences[0]?.position ?? "",
);

const printDate = computed(() => {
  const d = new Date();
  return d
    .toLocaleDateString("mn-MN", { month: "long", year: "numeric" })
    .toUpperCase();
});

const cvPhoto = ref<string>("");
const photoInput = ref<HTMLInputElement | null>(null);

function onPhotoChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  const reader = new FileReader();
  reader.onload = (e) => {
    cvPhoto.value = (e.target?.result as string) ?? "";
  };
  reader.readAsDataURL(file);
}
</script>

<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-20">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <!-- ===== SCREEN VIEW ===== -->
      <div class="print:hidden space-y-6">
        <!-- Header -->
        <div
          class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between"
        >
          <div class="space-y-1">
            <h1 class="text-[26px] font-semibold tracking-tight">Миний CV</h1>
            <p class="text-sm text-muted-foreground">
              Мэдээллээ бөглөөд PDF болгон татаж авах эсвэл PDF-ийг импортлон
              автоматаар бөглүүлэх боломжтой.
            </p>
          </div>
          <div class="flex flex-wrap gap-2">
            <input
              ref="fileInput"
              type="file"
              accept=".pdf"
              class="hidden"
              @change="onImportPDF"
            />
            <Button
              variant="outline"
              class="gap-2 rounded-full"
              :disabled="parsing"
              @click="fileInput?.click()"
            >
              <Loader2 v-if="parsing" class="h-4 w-4 animate-spin" />
              <Upload v-else class="h-4 w-4" />
              {{ parsing ? "Уншиж байна..." : "PDF импортлох" }}
            </Button>
            <Button variant="outline" class="gap-2 rounded-full" @click="print">
              <Printer class="h-4 w-4" />
              PDF татах
            </Button>
            <Button class="gap-2 rounded-full" :disabled="saving" @click="save">
              <Loader2 v-if="saving" class="h-4 w-4 animate-spin" />
              <Save v-else class="h-4 w-4" />
              {{ saving ? "Хадгалж байна..." : "Хадгалах" }}
            </Button>
          </div>
        </div>

        <div class="flex gap-6">
          <!-- Sidebar nav -->
          <aside class="hidden w-52 shrink-0 lg:block">
            <div
              class="sticky top-6 space-y-1 rounded-3xl border border-border bg-card px-3 py-3"
            >
              <button
                v-for="s in sections"
                :key="s.id"
                type="button"
                class="flex w-full items-center gap-2.5 rounded-2xl px-3 py-2.5 text-left text-sm text-muted-foreground transition hover:bg-muted hover:text-foreground"
                @click="scrollTo(s.id)"
              >
                <component :is="s.icon" class="h-4 w-4 shrink-0" />
                {{ s.label }}
              </button>
            </div>
          </aside>

          <!-- Form -->
          <div class="min-w-0 flex-1 space-y-6">
            <!-- General Info -->
            <section
              id="general"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <h2 class="mb-5 text-lg font-semibold flex items-center gap-2">
                <ScrollText class="h-5 w-5 text-muted-foreground" />
                Ерөнхий мэдээлэл
              </h2>
              <!-- Photo upload -->
              <div class="mb-5 flex items-center gap-4">
                <input
                  ref="photoInput"
                  type="file"
                  accept="image/*"
                  class="hidden"
                  @change="onPhotoChange"
                />
                <button
                  type="button"
                  class="w-20 h-20 rounded-full border-2 border-dashed border-border bg-muted flex flex-col items-center justify-center text-muted-foreground text-xs text-center gap-0.5 overflow-hidden hover:border-primary/50 transition shrink-0"
                  @click="photoInput?.click()"
                >
                  <template v-if="cvPhoto">
                    <img :src="cvPhoto" class="w-full h-full object-cover" />
                  </template>
                  <template v-else>
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mb-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"/><path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
                    <span>Зураг</span>
                  </template>
                </button>
                <div class="text-sm text-muted-foreground">
                  <p class="font-medium text-foreground">Профайл зураг</p>
                  <p class="text-xs mt-0.5">Заавал биш — PDF-д харагдана</p>
                  <button
                    v-if="cvPhoto"
                    type="button"
                    class="mt-1.5 text-xs text-destructive hover:underline"
                    @click="cvPhoto = ''"
                  >Устгах</button>
                </div>
              </div>

              <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Овог</label>
                  <Input v-model="cv.last_name" placeholder="Овог" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Нэр</label>
                  <Input v-model="cv.first_name" placeholder="Нэр" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Төрсөн огноо</label>
                  <Input v-model="cv.date_of_birth" type="date" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Хүйс</label>
                  <Select v-model="cv.gender">
                    <SelectTrigger class="rounded-xl">
                      <SelectValue placeholder="Сонгох" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem v-for="g in GENDERS" :key="g" :value="g">{{
                        g
                      }}</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Регистрийн дугаар</label>
                  <Input v-model="cv.national_id" placeholder="АА00000000" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Гэрлэлтийн байдал</label>
                  <Select v-model="cv.marital_status">
                    <SelectTrigger class="rounded-xl">
                      <SelectValue placeholder="Сонгох" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="m in MARITAL_STATUSES"
                        :key="m"
                        :value="m"
                        >{{ m }}</SelectItem
                      >
                    </SelectContent>
                  </Select>
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Утасны дугаар</label>
                  <Input v-model="cv.phone" placeholder="+976 9900 0000" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-sm font-medium">Э-мэйл</label>
                  <Input
                    v-model="cv.email"
                    type="email"
                    placeholder="mail@example.com"
                  />
                </div>
                <div class="space-y-1.5 sm:col-span-2">
                  <label class="text-sm font-medium">Хаяг</label>
                  <Input
                    v-model="cv.address"
                    placeholder="Улаанбаатар, Сүхбаатар дүүрэг..."
                  />
                </div>
                <div class="space-y-2 sm:col-span-2">
                  <label class="text-sm font-medium">Жолооны үнэмлэх</label>
                  <div class="flex flex-wrap gap-2">
                    <button
                      v-for="lic in DRIVER_LICENSES"
                      :key="lic"
                      type="button"
                      class="h-9 w-9 rounded-xl border text-sm font-semibold transition"
                      :class="
                        cv.driver_licenses.includes(lic)
                          ? 'border-primary bg-primary text-primary-foreground'
                          : 'border-border bg-muted text-muted-foreground hover:bg-muted/80'
                      "
                      @click="toggleLicense(lic)"
                    >
                      {{ lic }}
                    </button>
                  </div>
                </div>
              </div>
            </section>

            <!-- About -->
            <section
              id="about"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <h2 class="mb-5 text-lg font-semibold flex items-center gap-2">
                <BookOpen class="h-5 w-5 text-muted-foreground" />
                Миний тухай
              </h2>
              <Textarea
                v-model="cv.about"
                placeholder="Өөрийнхөө тухай товч танилцуулга бичнэ үү..."
                class="min-h-28 resize-y rounded-xl"
              />
            </section>

            <!-- Work Experience -->
            <section
              id="work"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <Briefcase class="h-5 w-5 text-muted-foreground" />
                  Ажлын туршлага
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.work_experiences, emptyWorkExperience)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(exp, i) in cv.work_experiences"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.work_experiences, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-2">
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Байгууллага</label
                      >
                      <Input
                        v-model="exp.company"
                        placeholder="Компанийн нэр"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Албан тушаал</label
                      >
                      <Input v-model="exp.position" placeholder="Ажлын нэр" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Эхлэх огноо</label
                      >
                      <Input v-model="exp.start_date" type="month" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Дуусах огноо</label
                      >
                      <Input
                        v-model="exp.end_date"
                        type="month"
                        :disabled="exp.current"
                      />
                    </div>
                    <div class="flex items-center gap-2">
                      <Checkbox
                        :id="`work-current-${i}`"
                        v-model:checked="exp.current"
                      />
                      <label :for="`work-current-${i}`" class="text-sm"
                        >Одоогоор ажиллаж байна</label
                      >
                    </div>
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Тайлбар</label
                      >
                      <Textarea
                        v-model="exp.description"
                        placeholder="Гүйцэтгэсэн үүрэг, амжилтууд..."
                        class="min-h-20 resize-y rounded-xl"
                      />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.work_experiences.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Ажлын туршлага нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Education -->
            <section
              id="education"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <GraduationCap class="h-5 w-5 text-muted-foreground" />
                  Боловсрол
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.education, emptyEducation)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(edu, i) in cv.education"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.education, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-2">
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Сургуулийн нэр</label
                      >
                      <Input
                        v-model="edu.school"
                        placeholder="Монгол Улсын Их Сургууль"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Зэрэг</label
                      >
                      <Input
                        v-model="edu.degree"
                        placeholder="Бакалавр, Магистр..."
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Чиглэл / Мэргэжил</label
                      >
                      <Input
                        v-model="edu.field"
                        placeholder="Компьютерийн шинжлэх ухаан"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Элссэн огноо</label
                      >
                      <Input v-model="edu.start_date" type="month" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Төгссөн огноо</label
                      >
                      <Input
                        v-model="edu.end_date"
                        type="month"
                        :disabled="edu.current"
                      />
                    </div>
                    <div class="flex items-center gap-2">
                      <Checkbox
                        :id="`edu-current-${i}`"
                        v-model:checked="edu.current"
                      />
                      <label :for="`edu-current-${i}`" class="text-sm"
                        >Одоогоор суралцаж байна</label
                      >
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >GPA / Дундаж оноо</label
                      >
                      <Input v-model="edu.gpa" placeholder="3.8 / 4.0" />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.education.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Боловсролын мэдээлэл нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Skills -->
            <section
              id="skills"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <h2 class="mb-5 text-lg font-semibold flex items-center gap-2">
                <Brain class="h-5 w-5 text-muted-foreground" />
                Мэдлэг ур чадвар
              </h2>
              <div class="space-y-6">
                <!-- Simple tag skill sections -->
                <div
                  v-for="(cfg, idx) in [
                    {
                      label: 'Хувийн ур чадвар',
                      field: 'personal_skills' as keyof CVProfile,
                      key: 'personal',
                      placeholder: 'Харилцааны ур чадвар...',
                    },
                    {
                      label: 'Мэргэжлийн ур чадвар',
                      field: 'professional_skills' as keyof CVProfile,
                      key: 'professional',
                      placeholder: 'JavaScript, Python...',
                    },
                    {
                      label: 'Компьютерын программ',
                      field: 'computer_skills' as keyof CVProfile,
                      key: 'computer',
                      placeholder: 'Microsoft Office, Photoshop...',
                    },
                    {
                      label: 'Урлагийн ур чадвар',
                      field: 'art_skills' as keyof CVProfile,
                      key: 'art',
                      placeholder: 'Дуу хөгжим, Зураг...',
                    },
                    {
                      label: 'Спортын ур чадвар',
                      field: 'sport_skills' as keyof CVProfile,
                      key: 'sport',
                      placeholder: 'Бөх, Сагсан бөмбөг...',
                    },
                  ]"
                  :key="cfg.key"
                >
                  <div class="space-y-2">
                    <label class="text-sm font-medium">{{ cfg.label }}</label>
                    <div class="flex flex-wrap gap-1.5 mb-2">
                      <span
                        v-for="(skill, si) in cv[cfg.field] as string[]"
                        :key="si"
                        class="inline-flex items-center gap-1 rounded-full border border-border bg-muted px-2.5 py-0.5 text-xs font-medium"
                      >
                        {{ skill }}
                        <button
                          type="button"
                          class="ml-0.5 text-muted-foreground hover:text-destructive"
                          @click="removeSkill(cfg.field, si)"
                        >
                          ×
                        </button>
                      </span>
                    </div>
                    <div class="flex gap-2">
                      <Input
                        v-model="skillInput[cfg.key]"
                        :placeholder="cfg.placeholder"
                        class="flex-1"
                        @keydown="
                          (e: KeyboardEvent) =>
                            onSkillKeydown(e, cfg.field, cfg.key)
                        "
                      />
                      <Button
                        type="button"
                        size="sm"
                        variant="outline"
                        class="shrink-0 rounded-xl"
                        @click="addSkill(cfg.field, cfg.key)"
                      >
                        <Plus class="h-3.5 w-3.5" />
                      </Button>
                    </div>
                  </div>
                </div>

                <!-- Languages -->
                <div class="space-y-2">
                  <div class="flex items-center justify-between">
                    <label class="text-sm font-medium"
                      >Гадаад хэлний мэдлэг</label
                    >
                    <Button
                      size="sm"
                      variant="outline"
                      class="gap-1 rounded-full"
                      @click="addItem(cv.languages, emptyLanguage)"
                    >
                      <Plus class="h-3.5 w-3.5" /> Нэмэх
                    </Button>
                  </div>
                  <div class="space-y-2">
                    <div
                      v-for="(lang, i) in cv.languages"
                      :key="i"
                      class="flex items-center gap-2"
                    >
                      <Input
                        v-model="lang.name"
                        placeholder="Англи, Орос..."
                        class="flex-1"
                      />
                      <Select v-model="lang.level">
                        <SelectTrigger class="w-40 rounded-xl">
                          <SelectValue placeholder="Түвшин" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem
                            v-for="lvl in LANGUAGE_LEVELS"
                            :key="lvl"
                            :value="lvl"
                            >{{ lvl }}</SelectItem
                          >
                        </SelectContent>
                      </Select>
                      <button
                        type="button"
                        class="p-1 text-muted-foreground hover:text-destructive"
                        @click="removeItem(cv.languages, i)"
                      >
                        <Trash2 class="h-4 w-4" />
                      </button>
                    </div>
                    <p
                      v-if="!cv.languages.length"
                      class="text-sm text-muted-foreground"
                    >
                      Хэлний мэдлэг нэмэгдээгүй байна
                    </p>
                  </div>
                </div>
              </div>
            </section>

            <!-- Trainings -->
            <section
              id="trainings"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <BookOpen class="h-5 w-5 text-muted-foreground" />
                  Сургалт, сертификат
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.trainings, emptyTraining)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(tr, i) in cv.trainings"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.trainings, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-2">
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Сургалтын нэр</label
                      >
                      <Input v-model="tr.name" placeholder="Сургалтын нэр" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Байгууллага</label
                      >
                      <Input
                        v-model="tr.organization"
                        placeholder="Зохион байгуулагч"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Огноо</label
                      >
                      <Input v-model="tr.date" type="month" />
                    </div>
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Сертификат / Тэмдэглэл</label
                      >
                      <Input
                        v-model="tr.certificate"
                        placeholder="Сертификатын дугаар эсвэл нэр"
                      />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.trainings.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Сургалт нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Exams -->
            <section
              id="exams"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <FlaskConical class="h-5 w-5 text-muted-foreground" />
                  Шалгалт, тест
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.exams, emptyExam)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(ex, i) in cv.exams"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.exams, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-3">
                    <div class="space-y-1 sm:col-span-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Шалгалтын нэр</label
                      >
                      <Input
                        v-model="ex.name"
                        placeholder="IELTS, TOEFL, HSK..."
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Оноо / Дүн</label
                      >
                      <Input v-model="ex.score" placeholder="7.5, B2..." />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Огноо</label
                      >
                      <Input v-model="ex.date" type="month" />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.exams.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Шалгалт нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Internships -->
            <section
              id="internships"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <Building2 class="h-5 w-5 text-muted-foreground" />
                  Дадлага
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.internships, emptyInternship)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(int, i) in cv.internships"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.internships, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-2">
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Байгууллага</label
                      >
                      <Input
                        v-model="int.company"
                        placeholder="Компанийн нэр"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Албан тушаал</label
                      >
                      <Input
                        v-model="int.position"
                        placeholder="Дадлагажигч, Туслах..."
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Эхлэх огноо</label
                      >
                      <Input v-model="int.start_date" type="month" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Дуусах огноо</label
                      >
                      <Input v-model="int.end_date" type="month" />
                    </div>
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Тайлбар</label
                      >
                      <Textarea
                        v-model="int.description"
                        placeholder="Дадлагын үеийн туршлага, сурсан зүйлс..."
                        class="min-h-20 resize-y rounded-xl"
                      />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.internships.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Дадлага нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Awards -->
            <section
              id="awards"
              class="rounded-3xl border border-border bg-card px-6 py-6"
            >
              <div class="mb-5 flex items-center justify-between">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                  <Trophy class="h-5 w-5 text-muted-foreground" />
                  Шагнал урамшуулал
                </h2>
                <Button
                  size="sm"
                  variant="outline"
                  class="gap-1.5 rounded-full"
                  @click="addItem(cv.awards, emptyAward)"
                >
                  <Plus class="h-3.5 w-3.5" /> Нэмэх
                </Button>
              </div>
              <div class="space-y-4">
                <div
                  v-for="(aw, i) in cv.awards"
                  :key="i"
                  class="relative rounded-2xl border border-border p-4"
                >
                  <button
                    type="button"
                    class="absolute right-3 top-3 rounded-lg p-1 text-muted-foreground hover:text-destructive"
                    @click="removeItem(cv.awards, i)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                  <div class="grid gap-3 sm:grid-cols-2">
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Шагналын нэр</label
                      >
                      <Input v-model="aw.name" placeholder="Шагналын нэр" />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Байгууллага</label
                      >
                      <Input
                        v-model="aw.organization"
                        placeholder="Шагнал олгосон байгууллага"
                      />
                    </div>
                    <div class="space-y-1">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Огноо</label
                      >
                      <Input v-model="aw.date" type="month" />
                    </div>
                    <div class="space-y-1 sm:col-span-2">
                      <label class="text-xs font-medium text-muted-foreground"
                        >Тайлбар</label
                      >
                      <Textarea
                        v-model="aw.description"
                        placeholder="Шагналын тухай дэлгэрэнгүй мэдээлэл..."
                        class="min-h-20 resize-y rounded-xl"
                      />
                    </div>
                  </div>
                </div>
                <p
                  v-if="!cv.awards.length"
                  class="text-center text-sm text-muted-foreground py-4"
                >
                  Шагнал нэмэгдээгүй байна
                </p>
              </div>
            </section>

            <!-- Bottom save -->
            <div class="flex justify-end gap-2 pb-6">
              <Button
                variant="outline"
                class="gap-2 rounded-full"
                @click="print"
              >
                <Printer class="h-4 w-4" /> PDF татах
              </Button>
              <Button
                class="gap-2 rounded-full"
                :disabled="saving"
                @click="save"
              >
                <Loader2 v-if="saving" class="h-4 w-4 animate-spin" />
                <Save v-else class="h-4 w-4" />
                {{ saving ? "Хадгалж байна..." : "Хадгалах" }}
              </Button>
            </div>
          </div>
        </div>
      </div>

      <!-- ===== PRINT VIEW ===== -->
      <div class="hidden print:block">
        <div class="cv-page">
          <div class="cv-page-inner">

            <!-- HEADER: with photo -->
            <header v-if="cvPhoto" class="cv-head cv-head--photo">
              <div class="cv-photo">
                <img :src="cvPhoto" alt="profile" />
              </div>
              <div>
                <div class="cv-name">{{ fullName || "Нэр оруулаагүй" }}</div>
                <div v-if="currentJobTitle" class="cv-role">{{ currentJobTitle }}</div>
              </div>
              <div class="cv-contact">
                <div v-if="cv.email" class="cv-contact-row"><span class="cv-lbl">Имэйл</span><span>{{ cv.email }}</span></div>
                <div v-if="cv.phone" class="cv-contact-row"><span class="cv-lbl">Утас</span><span>{{ cv.phone }}</span></div>
                <div v-if="cv.address" class="cv-contact-row"><span class="cv-lbl">Хаяг</span><span>{{ cv.address }}</span></div>
                <div v-if="cv.date_of_birth" class="cv-contact-row"><span class="cv-lbl">Төрсөн</span><span>{{ cv.date_of_birth }}</span></div>
                <div v-if="cv.national_id" class="cv-contact-row"><span class="cv-lbl">Рег</span><span>{{ cv.national_id }}</span></div>
                <div v-if="cv.driver_licenses.length" class="cv-contact-row"><span class="cv-lbl">Жолоо</span><span>{{ cv.driver_licenses.join(", ") }}</span></div>
              </div>
            </header>

            <!-- HEADER: without photo -->
            <header v-else class="cv-head cv-head--no-photo">
              <div>
                <div class="cv-name">{{ fullName || "Нэр оруулаагүй" }}</div>
                <div v-if="currentJobTitle" class="cv-role">{{ currentJobTitle }}</div>
              </div>
              <div class="cv-contact">
                <div v-if="cv.email" class="cv-contact-row"><span class="cv-lbl">Имэйл</span><span>{{ cv.email }}</span></div>
                <div v-if="cv.phone" class="cv-contact-row"><span class="cv-lbl">Утас</span><span>{{ cv.phone }}</span></div>
                <div v-if="cv.address" class="cv-contact-row"><span class="cv-lbl">Хаяг</span><span>{{ cv.address }}</span></div>
                <div v-if="cv.date_of_birth" class="cv-contact-row"><span class="cv-lbl">Төрсөн</span><span>{{ cv.date_of_birth }}</span></div>
                <div v-if="cv.national_id" class="cv-contact-row"><span class="cv-lbl">Рег</span><span>{{ cv.national_id }}</span></div>
                <div v-if="cv.driver_licenses.length" class="cv-contact-row"><span class="cv-lbl">Жолоо</span><span>{{ cv.driver_licenses.join(", ") }}</span></div>
              </div>
            </header>

            <!-- ABOUT -->
            <section v-if="cv.about" class="cv-sec">
              <div class="cv-sec-h"><span class="cv-dot"></span><span class="cv-sec-title">Миний тухай</span><span class="cv-rule"></span></div>
              <p class="cv-about">{{ cv.about }}</p>
            </section>

            <!-- TWO COLUMNS -->
            <div class="cv-body">

              <!-- LEFT -->
              <div class="cv-col">

                <section v-if="cv.professional_skills.length || cv.personal_skills.length || cv.computer_skills.length || cv.art_skills.length || cv.sport_skills.length || cv.languages.length" class="cv-sec">
                  <div class="cv-sec-h"><span class="cv-dot"></span><span class="cv-sec-title">Ур чадвар</span><span class="cv-rule"></span></div>
                  <div v-if="cv.professional_skills.length" class="cv-skill-group">
                    <div class="cv-skill-label">Мэргэжлийн</div>
                    <div class="cv-skill-tags"><span v-for="s in cv.professional_skills" :key="s" class="cv-skill-tag">{{ s }}</span></div>
                  </div>
                  <div v-if="cv.personal_skills.length" class="cv-skill-group">
                    <div class="cv-skill-label">Хувийн</div>
                    <div class="cv-skill-tags"><span v-for="s in cv.personal_skills" :key="s" class="cv-skill-tag">{{ s }}</span></div>
                  </div>
                  <div v-if="cv.computer_skills.length" class="cv-skill-group">
                    <div class="cv-skill-label">Компьютер</div>
                    <div class="cv-skill-tags"><span v-for="s in cv.computer_skills" :key="s" class="cv-skill-tag">{{ s }}</span></div>
                  </div>
                  <div v-if="cv.languages.length" class="cv-skill-group">
                    <div class="cv-skill-label">Хэл</div>
                    <div class="cv-skill-tags">
                      <span v-for="l in cv.languages" :key="l.name" class="cv-skill-tag">{{ l.name }}<template v-if="l.level"> ({{ l.level }})</template></span>
                    </div>
                  </div>
                  <div v-if="cv.art_skills.length" class="cv-skill-group">
                    <div class="cv-skill-label">Урлаг</div>
                    <div class="cv-skill-tags"><span v-for="s in cv.art_skills" :key="s" class="cv-skill-tag">{{ s }}</span></div>
                  </div>
                  <div v-if="cv.sport_skills.length" class="cv-skill-group">
                    <div class="cv-skill-label">Спорт</div>
                    <div class="cv-skill-tags"><span v-for="s in cv.sport_skills" :key="s" class="cv-skill-tag">{{ s }}</span></div>
                  </div>
                </section>

                <section v-if="cv.trainings.length" class="cv-sec">
                  <div class="cv-sec-h"><span class="cv-dot"></span><span class="cv-sec-title">Сургалт</span><span class="cv-rule"></span></div>
                  <div v-for="tr in cv.trainings" :key="tr.name" class="cv-row-item">
                    <div class="cv-row-top"><div class="cv-row-title">{{ tr.name }}</div><div class="cv-row-meta">{{ tr.date }}</div></div>
                    <div class="cv-row-sub">{{ tr.organization }}<template v-if="tr.certificate"> · {{ tr.certificate }}</template></div>
                  </div>
                </section>

                <section v-if="cv.exams.length" class="cv-sec">
                  <div class="cv-sec-h"><span class="cv-dot"></span><span class="cv-sec-title">Шалгалт</span><span class="cv-rule"></span></div>
                  <div v-for="ex in cv.exams" :key="ex.name" class="cv-row-item">
                    <div class="cv-row-top">
                      <div class="cv-row-title">{{ ex.name }}</div>
                      <div class="cv-row-meta">{{ ex.date }}<template v-if="ex.score"> · {{ ex.score }}</template></div>
                    </div>
                  </div>
                </section>

              </div>

              <!-- RIGHT -->
              <div class="cv-col">

                <section v-if="cv.work_experiences.length" class="cv-sec">
                  <div class="cv-sec-h"><span class="cv-dot"></span><span class="cv-sec-title">Ажлын туршлага</span><span class="cv-rule"></span></div>
                  <div v-for="exp in cv.work_experiences" :key="exp.company + exp.start_date" class="cv-work-item">
                    <div class="cv-work-top">
                      <div class="cv-work-role">{{ exp.position }}</div>
                      <div class="cv-work-dates">{{ exp.start_date }} — {{ exp.current ? "Одоог хүртэл" : exp.end_date }}</div>
                    </div>
                    <div class="cv-work-sub"><span class="cv-work-co">{{ exp.company }}</span></div>
                    <ul v-if="exp.description" class="cv-work-bullets">
                      <li v-for="line in exp.description.split('\n').filter((l) => l.trim())" :key="line">{{ line }}</li>
                    </ul>
                  </div>
                </section>

                <section v-if="cv.internships.length" class="cv-sec">
                  <div class="cv-sec-h"><span class="cv-dot"></span><span class="cv-sec-title">Дадлага</span><span class="cv-rule"></span></div>
                  <div v-for="int in cv.internships" :key="int.company + int.start_date" class="cv-work-item">
                    <div class="cv-work-top">
                      <div class="cv-work-role">{{ int.position }}</div>
                      <div class="cv-work-dates">{{ int.start_date }} — {{ int.end_date }}</div>
                    </div>
                    <div class="cv-work-sub"><span class="cv-work-co">{{ int.company }}</span></div>
                    <ul v-if="int.description" class="cv-work-bullets">
                      <li v-for="line in int.description.split('\n').filter((l) => l.trim())" :key="line">{{ line }}</li>
                    </ul>
                  </div>
                </section>

                <section v-if="cv.education.length" class="cv-sec">
                  <div class="cv-sec-h"><span class="cv-dot"></span><span class="cv-sec-title">Боловсрол</span><span class="cv-rule"></span></div>
                  <div v-for="edu in cv.education" :key="edu.school" class="cv-row-item">
                    <div class="cv-row-top">
                      <div class="cv-row-title">{{ edu.degree }}<template v-if="edu.field"> {{ edu.field }}</template></div>
                      <div class="cv-row-meta">{{ edu.start_date }} — {{ edu.current ? "Одоог хүртэл" : edu.end_date }}</div>
                    </div>
                    <div class="cv-row-sub">{{ edu.school }}<template v-if="edu.gpa"> · GPA: {{ edu.gpa }}</template></div>
                  </div>
                </section>

                <section v-if="cv.awards.length" class="cv-sec">
                  <div class="cv-sec-h"><span class="cv-dot"></span><span class="cv-sec-title">Шагнал</span><span class="cv-rule"></span></div>
                  <div v-for="aw in cv.awards" :key="aw.name" class="cv-row-item">
                    <div class="cv-row-top"><div class="cv-row-title">{{ aw.name }}</div><div class="cv-row-meta">{{ aw.date }}</div></div>
                    <div class="cv-row-sub">{{ aw.organization }}</div>
                    <div v-if="aw.description" class="cv-row-desc">{{ aw.description }}</div>
                  </div>
                </section>

              </div>
            </div>

          </div>

          <!-- PAGE FOOTER -->
          <footer class="cv-page-foot">
            <div class="cv-foot-left"><span class="cv-foot-mark"></span><span>{{ fullName }} · CV · {{ printDate }}</span></div>
            <div>Page 1 / 1</div>
          </footer>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
/* ── CV print design system ── */
.cv-page {
  width: 210mm;
  height: 297mm;
  background: #fff;
  position: relative;
  overflow: hidden;
  font-family: "DM Sans", -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif;
  -webkit-font-smoothing: antialiased;
  color: #1c1a23;
  font-size: 11px;
  line-height: 1.4;
  -webkit-print-color-adjust: exact;
  print-color-adjust: exact;
}

.cv-page-inner {
  padding: 22mm 20mm 16mm;
  display: flex;
  flex-direction: column;
  gap: 14px;
  height: 100%;
  box-sizing: border-box;
}

/* Header */
.cv-head {
  display: grid;
  gap: 20px;
  align-items: end;
  padding-bottom: 14px;
  border-bottom: 1px solid #e6e3ea;
}
.cv-head--photo { grid-template-columns: auto 1fr auto; }
.cv-head--no-photo { grid-template-columns: 1fr auto; }

.cv-photo {
  width: 88px; height: 88px;
  border-radius: 999px;
  overflow: hidden;
  box-shadow: 0 0 0 1px #e6e3ea, 0 2px 6px rgba(28,26,35,.06);
  flex-shrink: 0;
}
.cv-photo img { width: 100%; height: 100%; object-fit: cover; display: block; }

.cv-name {
  font-size: 30px; font-weight: 700; letter-spacing: -1px;
  line-height: 1; color: #1c1a23;
}
.cv-role {
  font-size: 13px; font-weight: 500; color: #6b5191;
  margin-top: 6px; letter-spacing: -0.2px;
}

.cv-contact {
  font-size: 10px;
  font-family: "JetBrains Mono", monospace;
  color: #6b6976;
  line-height: 1.9;
  text-align: right;
}
.cv-contact-row { display: flex; gap: 6px; justify-content: flex-end; align-items: center; }
.cv-lbl { color: #9a98a4; font-size: 9px; text-transform: uppercase; letter-spacing: 0.6px; }

/* Two-column body */
.cv-body {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 24px;
  flex: 1;
}
.cv-col { display: flex; flex-direction: column; gap: 14px; min-width: 0; }

/* Section header */
.cv-sec { break-inside: avoid; }
.cv-sec-h {
  display: flex; align-items: center; gap: 7px;
  margin-bottom: 9px;
}
.cv-dot {
  width: 5px; height: 5px; border-radius: 999px;
  background: #6b5191; flex-shrink: 0;
}
.cv-sec-title {
  font-size: 9.5px; font-weight: 700;
  letter-spacing: 1.4px; text-transform: uppercase;
  color: #6b5191;
  font-family: "JetBrains Mono", monospace;
  white-space: nowrap;
}
.cv-rule { flex: 1; height: 1px; background: #e6e3ea; }

/* About */
.cv-about { font-size: 12px; line-height: 1.6; color: #3c3a44; }

/* Work entries */
.cv-work-item { margin-bottom: 12px; break-inside: avoid; }
.cv-work-item:last-child { margin-bottom: 0; }
.cv-work-top {
  display: flex; justify-content: space-between; gap: 10px;
  align-items: baseline; margin-bottom: 2px;
}
.cv-work-role { font-size: 12.5px; font-weight: 600; color: #1c1a23; letter-spacing: -0.2px; }
.cv-work-dates {
  font-size: 10px; font-family: "JetBrains Mono", monospace;
  color: #6b6976; white-space: nowrap; letter-spacing: 0.2px;
}
.cv-work-sub { font-size: 11px; color: #6b6976; margin-bottom: 5px; }
.cv-work-co { color: #6b5191; font-weight: 600; }
.cv-work-bullets {
  margin: 0; padding: 0; list-style: none;
  display: flex; flex-direction: column; gap: 2px;
}
.cv-work-bullets li {
  font-size: 11px; line-height: 1.5; color: #3c3a44;
  padding-left: 13px; position: relative;
}
.cv-work-bullets li::before {
  content: ""; position: absolute; left: 2px; top: 7px;
  width: 4px; height: 4px; border-radius: 999px;
  background: #6b5191; opacity: 0.55;
}

/* Row items (education, training, exam, award) */
.cv-row-item { margin-bottom: 9px; break-inside: avoid; }
.cv-row-item:last-child { margin-bottom: 0; }
.cv-row-top {
  display: flex; justify-content: space-between; gap: 10px;
  align-items: baseline;
}
.cv-row-title { font-size: 11.5px; font-weight: 600; color: #1c1a23; letter-spacing: -0.1px; }
.cv-row-meta {
  font-size: 10px; font-family: "JetBrains Mono", monospace;
  color: #6b6976; white-space: nowrap;
}
.cv-row-sub { font-size: 10.5px; color: #6b6976; margin-top: 1px; }
.cv-row-desc { font-size: 10.5px; color: #3c3a44; margin-top: 3px; padding-left: 2px; }

/* Skills */
.cv-skill-group { margin-bottom: 9px; break-inside: avoid; }
.cv-skill-group:last-child { margin-bottom: 0; }
.cv-skill-label {
  font-size: 9px; font-weight: 600; color: #6b6976;
  text-transform: uppercase; letter-spacing: 0.8px;
  margin-bottom: 4px;
  font-family: "JetBrains Mono", monospace;
}
.cv-skill-tags { display: flex; flex-wrap: wrap; gap: 3px; }
.cv-skill-tag {
  font-size: 10.5px; padding: 2px 7px;
  border-radius: 4px;
  background: rgba(107, 81, 145, 0.08);
  color: #574077;
  font-weight: 500;
  letter-spacing: -0.1px;
}

/* Page footer */
.cv-page-foot {
  position: absolute;
  bottom: 12mm; left: 20mm; right: 20mm;
  display: flex; justify-content: space-between;
  font-size: 8.5px; font-family: "JetBrains Mono", monospace;
  color: #9a98a4; letter-spacing: 0.4px; text-transform: uppercase;
}
.cv-foot-left { display: flex; gap: 6px; align-items: center; }
.cv-foot-mark {
  width: 5px; height: 5px; border-radius: 999px;
  background: #6b5191; opacity: 0.35;
}
</style>
